package router

import (
	"github.com/fitant/xbin-api/config"
	"github.com/fitant/xbin-api/src/service"
	"github.com/fitant/xbin-api/src/view/http/handler/snippet"
	"github.com/fitant/xbin-api/src/view/http/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func New(svc service.Service, cfg *config.HTTPServerConfig, lgr *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.WithCors(cfg))

	r.Route(cfg.Enpoint, func(sr chi.Router) {
		sr.With(middleware.WithIngestion()).Put(
			"/", snippet.Create(svc, cfg, lgr))
		sr.With(middleware.WithIngestion()).Post(
			"/", snippet.Create(svc, cfg, lgr))
		sr.Get("/r/{snippetID}", snippet.Get(svc, lgr, "raw"))
		sr.Get("/{snippetID}", snippet.Get(svc, lgr, "json"))
	})

	return r
}
