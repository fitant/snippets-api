package model

import (
	"bytes"

	"github.com/fitant/xbin-api/src/storageprovider"
)

var ErrNotFound = storageprovider.ErrNotFound

type SnippetController interface {
	NewSnippet(id, language string, ephemeral bool, snip []byte) error
	FindSnippet(name string) (*Snippet, error)
}

type Snippet struct {
	ID        string
	Ephemeral bool
	Language  string
	Snippet   []byte
}

type mongoSnippetController struct {
	sp *storageprovider.S3Provider
}

func NewMongoSnippetController(sp *storageprovider.S3Provider) SnippetController {
	return &mongoSnippetController{
		sp: sp,
	}
}

func (msc *mongoSnippetController) NewSnippet(id, language string, ephemeral bool, snip []byte) error {
	err := msc.sp.UploadSnippet(bytes.NewReader(snip), id, language, ephemeral)
	if err != nil {
		return err
	}
	return nil
}

func (msc *mongoSnippetController) FindSnippet(id string) (*Snippet, error) {
	data, language, ephemeral, err := msc.sp.DownloadSnippet(id)
	if err != nil {
		return nil, err
	}

	// Create new Snippet and return
	return &Snippet{
		ID:        id,
		Ephemeral: ephemeral,
		Language:  language,
		Snippet:   data,
	}, nil
}
