package utils

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/leonklingele/passphrase"
)

func GenerateID(n int) string {
	passphrase.Separator = "-"
	id, _ := passphrase.Generate(n)

	idFields := strings.Split(id, "-")
	for i := range idFields {
		idFields[i] = strings.Title(idFields[i])
	}
	return strings.Join(idFields, "")
}

func DefalteBrotli(data []byte) []byte {
	var b bytes.Buffer
	w := brotli.NewWriterLevel(&b, brotli.BestCompression)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

func InflateBrotli(data []byte) []byte {
	var b bytes.Reader
	b.Read(data)
	r := brotli.NewReader(bytes.NewReader(data))
	x, _ := ioutil.ReadAll(r)
	return x
}
