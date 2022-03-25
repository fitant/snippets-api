package service

import (
	"encoding/base64"
	"fmt"

	"github.com/fitant/xbin-api/config"
	"github.com/fitant/xbin-api/src/db"
	"github.com/fitant/xbin-api/src/model"
	"github.com/fitant/xbin-api/src/types"
	"github.com/fitant/xbin-api/src/utils"
)

type Service interface {
	CreateSnippet(snippet, language string, ephemeral bool) (*model.Snippet, error)
	FetchSnippet(id string) (*model.Snippet, error)
}

type serviceImpl struct {
	sc        model.SnippetController
	overrides map[string]string
}

func NewSnippetService(sc model.SnippetController, cfg config.Service) Service {
	encryptionKeys = make(chan types.EncryptionStack, 10)
	go populateEncryptionStack(2)
	return &serviceImpl{
		sc:        sc,
		overrides: cfg.Overrides,
	}
}

func (s *serviceImpl) CreateSnippet(snippet, language string, ephemeral bool) (*model.Snippet, error) {
	keys := <-encryptionKeys

	// Deflate snippet -> Encrypt snippet -> encode snippet
	compressedSnippet := utils.DefalteBrotli([]byte(snippet))
	encryptedSnippet := utils.Encrypt(compressedSnippet, keys.Key, keys.Salt)
	snip, err := s.sc.NewSnippet(keys.Hash, base64.StdEncoding.EncodeToString(encryptedSnippet),
		language, ephemeral)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [CreateSnippet] [NewSnippet]", err))
		return nil, err
	}

	// Return generated ID instead of the stored hashed ID
	snip.ID = keys.ID

	return snip, nil
}

func (s *serviceImpl) FetchSnippet(id string) (*model.Snippet, error) {
	if s.overrides[id] != "" {
		id = s.overrides[id]
	}
	hashedID := utils.HashID([]byte(id))
	encodedID := base64.StdEncoding.EncodeToString(hashedID)
	snip, err := s.sc.FindSnippet(encodedID)
	if err != nil {
		if err != db.ErrNoDocuments {
			utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [FetchSnippet] [FindSnippet]", err))
		}
		return nil, err
	}

	// Decode snippet -> Decrypt snippet -> inflate snippet
	encryptedSnippet, _ := base64.StdEncoding.DecodeString(snip.Snippet)
	decryptedSnippet := utils.Decrypt(encryptedSnippet, []byte(id))
	snip.Snippet = string(utils.InflateBrotli(decryptedSnippet))

	// Return generated ID instead of the stored hashed ID
	snip.ID = id

	return snip, nil
}
