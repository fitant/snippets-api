package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fitant/xbin-api/config"
	"github.com/fitant/xbin-api/src/service"
	"github.com/fitant/xbin-api/src/view"
	"github.com/fitant/xbin-api/src/view/http/router"
	"go.uber.org/zap"
)

type webView struct {
	httpServer http.Server
	lgr        *zap.Logger
}

func (w *webView) Serve() {
	w.lgr.Info(fmt.Sprintf("[Views] [Web] [Start] Going to listening on http://%s", w.httpServer.Addr))
	go func() {
		if err := w.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			w.lgr.Fatal(fmt.Sprintf("[Views] [Web] [Start] [ListenAndServe]: %s", err.Error()))
			panic(err)
		}
	}()
}

func (w *webView) Shutdown(ctx context.Context) {
	if err := w.httpServer.Shutdown(ctx); err != nil {
		w.lgr.Error(fmt.Sprintf("[Views] [Web] [Shutdown]: %v", err))
		panic(err)
	}
}

func Init(svc service.Service, cfg *config.HTTPServerConfig, lgr *zap.Logger) view.View {
	rtr := router.New(svc, lgr)
	srv := &http.Server{
		Addr:    cfg.GetListenAddr(),
		Handler: rtr,
	}

	return &webView{
		httpServer: *srv,
		lgr:        lgr,
	}
}
