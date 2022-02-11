package service

import (
	"encoding/base64"
	"fmt"

	"github.com/fitant/xbin-api/src/db"
	"github.com/fitant/xbin-api/src/model"
	"github.com/fitant/xbin-api/src/types"
	"github.com/fitant/xbin-api/src/utils"
	"go.uber.org/zap"
)

type Service interface {
	CreateSnippet(snippet, language string, ephemeral bool) (*model.Snippet, error)
	FetchSnippet(id string) (*model.Snippet, error)
}

type serviceImpl struct {
	uidSize int
	salt    []byte
	sc      model.SnippetController
	cipher  types.CipherSelection
	lgr     *zap.Logger
}

func NewSnippetService(sc model.SnippetController, salt []byte, cipher types.CipherSelection, lgr *zap.Logger) Service {
	if len(salt) == 0 {
		panic("salt not specified")
	}
	return &serviceImpl{
		lgr:     lgr,
		sc:      sc,
		cipher:  cipher,
		uidSize: 2,
		salt:    salt,
	}
}

func (s *serviceImpl) CreateSnippet(snippet, language string, ephemeral bool) (*model.Snippet, error) {
	id := utils.GenerateID(s.uidSize)
	hashedID := utils.HashID([]byte(id), s.salt)
	encodedID := base64.StdEncoding.EncodeToString(hashedID)

	// Deflate snippet -> Encrypt snippet -> encode snippet
	compressedSnippet := utils.DefalteBrotli([]byte(snippet))
	encryptedSnippet := utils.Encrypt(compressedSnippet, []byte(id), s.cipher)
	snip, err := s.sc.NewSnippet(encodedID, base64.StdEncoding.EncodeToString(encryptedSnippet),
		language, ephemeral)
	if err != nil {
		if err == db.ErrDuplicateKey {
			s.uidSize++
			return s.CreateSnippet(snippet, language, ephemeral)
		}
		s.lgr.Error(fmt.Sprintf("%s : %v", "[Service] [CreateSnippet] [NewSnippet]", err))
		return nil, err
	}

	// Return generated ID instead of the stored hashed ID
	snip.ID = id

	return snip, nil
}

func (s *serviceImpl) FetchSnippet(id string) (*model.Snippet, error) {
	hashedID := utils.HashID([]byte(id), s.salt)
	encodedID := base64.StdEncoding.EncodeToString(hashedID)
	snip, err := s.sc.FindSnippet(encodedID)
	if err != nil {
		if err != db.ErrNoDocuments {
			s.lgr.Error(fmt.Sprintf("%s : %v", "[Service] [FetchSnippet] [FindSnippet]", err))
		}
		return nil, err
	}

	// Decode snippet -> Decrypt snippet -> inflate snippet
	encryptedSnippet, _ := base64.StdEncoding.DecodeString(snip.Snippet)
	decryptedSnippet := utils.Decrypt(encryptedSnippet, []byte(id), s.cipher)
	snip.Snippet = string(utils.InflateBrotli(decryptedSnippet))

	// Return generated ID instead of the stored hashed ID
	snip.ID = id

	return snip, nil
}
