package router

import (
	"github.com/fitant/xbin-api/src/service"
	"github.com/fitant/xbin-api/src/view/http/handler/snippet"
	"github.com/fitant/xbin-api/src/view/http/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func New(svc service.Service, lgr *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/snippets", func(sr chi.Router) {
		sr.With(middleware.WithIngestion()).Put(
			"/{filename}", snippet.Create(svc, lgr))
		sr.With(middleware.WithIngestion()).Post(
			"/", snippet.Create(svc, lgr))
		sr.Get("/r/{snippetID}", snippet.Get(svc, lgr, "raw"))
		sr.Get("/{snippetID}", snippet.Get(svc, lgr, "json"))
	})

	return r
}
