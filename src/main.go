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
	"github.com/fitant/xbin-api/src/view/http"
)

func main() {
	cfg := config.Load()
	utils.InitLogger(cfg)

	dbInstance, err := db.NewMongoStore(cfg)
	if err != nil {
		panic(err)
	}

	sc := model.NewMongoSnippetController(dbInstance)
	svc := service.NewSnippetService(sc, cfg.Svc)

	// Initialise and start serving webview
	httpView := http.Init(svc, &cfg.Http)
	httpView.Serve()

	gracefulShutdown([]view.View{httpView})
}

func gracefulShutdown(views []view.View) {
	// Listen for interrrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Once interrupted - shutdown all views
	utils.Logger.Info("[Main] [gracefulShutdown]: Attempting GracefulShutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, view := range views {
		go view.Shutdown(ctx)
	}
}
