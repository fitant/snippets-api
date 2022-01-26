package service

import (
	"encoding/base64"
	"fmt"

	"github.com/fitant/xbin-api/src/db"
	"github.com/fitant/xbin-api/src/model"
	"go.uber.org/zap"
)

type Service interface {
	CreateSnippet(snippet, language string, ephemeral bool) (*model.Snippet, error)
	FetchSnippet(id string) (*model.Snippet, error)
}

type serviceImpl struct {
	uidSize int
	sc      model.SnippetController
	lgr     *zap.Logger
}

func NewSnippetService(sc model.SnippetController, lgr *zap.Logger) Service {
	return &serviceImpl{
		lgr:     lgr,
		sc:      sc,
		uidSize: 2,
	}
}

func (s *serviceImpl) CreateSnippet(snippet, language string, ephemeral bool) (*model.Snippet, error) {
	id := generateID(s.uidSize)
	_, err := s.FetchSnippet(id)
	if err != db.ErrNoDocuments {
		s.uidSize++
		return s.CreateSnippet(snippet, language, ephemeral)
	}

	snip, err := s.sc.NewSnippet(id, base64.StdEncoding.EncodeToString([]byte(snippet)),
		language, ephemeral)
	if err != nil {
		if err == db.ErrDuplicateKey {
			s.uidSize++
			return s.CreateSnippet(snippet, language, ephemeral)
		}
		s.lgr.Error(fmt.Sprintf("%s : %v", "[Service] [CreateSnippet] [NewSnippet]", err))
		return nil, err
	}

	return snip, nil
}

func (s *serviceImpl) FetchSnippet(id string) (*model.Snippet, error) {
	snip, err := s.sc.FindSnippet(id)
	if err != nil {
		if err != db.ErrNoDocuments {
			s.lgr.Error(fmt.Sprintf("%s : %v", "[Service] [FetchSnippet] [FindSnippet]", err))
		}
		return nil, err
	}

	str, _ := base64.StdEncoding.DecodeString(snip.Snippet)
	snip.Snippet = string(str)

	return snip, nil
}
