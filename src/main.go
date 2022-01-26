package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/fitant/xbin-api/config"
	"github.com/fitant/xbin-api/src/db"
	"github.com/fitant/xbin-api/src/model"
	"github.com/fitant/xbin-api/src/service"
	"github.com/fitant/xbin-api/src/utils"
	"github.com/fitant/xbin-api/src/view"
	"github.com/fitant/xbin-api/src/view/web"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	lgr := utils.InitLogger(cfg)

	dbInstance, err := db.NewMongoStore(cfg, lgr)
	if err != nil {
		panic(err)
	}

	sc := model.NewMongoSnippetController(dbInstance, lgr)
	svc := service.NewSnippetService(sc, lgr)

	// Initialise and start serving webview
	webView := web.Init(svc, &cfg.Http, lgr)
	webView.Serve()

	gracefulShutdown([]view.View{webView}, lgr)
}

func gracefulShutdown(views []view.View, logger *zap.Logger) {
	// Listen for interrrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Once interrupted - shutdown all views
	logger.Info("[Main] [gracefulShutdown]: Attempting GracefulShutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, view := range views {
		go view.Shutdown(ctx)
	}
}
