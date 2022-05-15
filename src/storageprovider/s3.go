package storageprovider

import (
	"fmt"
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/fitant/xbin-api/config"
)

type S3Provider struct {
	bucketname string
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	service    *s3.S3
}

func InitS3StorageProvider() *S3Provider {
	sess := session.Must(session.NewSession())

	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)
	service := s3.New(sess)

	return &S3Provider{
		uploader:   uploader,
		downloader: downloader,
		service:    service,
		bucketname: config.Cfg.S3.Bucket,
	}
}

func (sp *S3Provider) UploadSnippet(data io.Reader, id string, language string, ephemeral bool) error {
	e := strconv.FormatBool(ephemeral)
	_, err := sp.uploader.Upload(&s3manager.UploadInput{
		Body:   data,
		Bucket: &sp.bucketname,
		Key:    &id,
		Metadata: map[string]*string{
			"Language":  &language,
			"Ephemeral": &e,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (sp *S3Provider) DownloadSnippet(id string) (data []byte, language string, ephemeral bool, err error) {
	objhead, err := sp.service.HeadObject(&s3.HeadObjectInput{
		Bucket: &sp.bucketname,
		Key:    &id,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				return nil, "", false, ErrNotFound
			default:
				return nil, "", false, err
			}
		}
		return nil, "", false, err
	}

	data = make([]byte, *objhead.ContentLength)

	_, err = sp.downloader.Download(aws.NewWriteAtBuffer(data), &s3.GetObjectInput{
		Bucket: &sp.bucketname,
		Key:    &id,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("+%v\n", objhead.Metadata)

	eph := objhead.Metadata["Ephemeral"]
	if eph != nil {
		ephemeral, _ = strconv.ParseBool(*eph)
		fmt.Println(ephemeral)
	}

	lang := objhead.Metadata["Language"]
	if lang != nil {
		language = *lang
		fmt.Println(lang)
	}

	return
}
