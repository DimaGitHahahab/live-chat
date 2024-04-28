package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"chat-service/internal/server"
	"chat-service/pkg/config"
	"chat-service/pkg/signal"

	log "github.com/sirupsen/logrus"
)

type App struct {
	config  *config.Config
	sigQuit chan os.Signal
	server  *server.Server
}

func New(cfg *config.Config) *App {
	sigQuit := signal.GetShutdownChannel()

	srv := server.New(cfg)

	return &App{
		config:  cfg,
		sigQuit: sigQuit,
		server:  srv,
	}
}

func (a *App) Run() {
	go func() {
		log.Infoln("Starting server on port ", a.config.HttpPort)
		if err := a.server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("Failed to start server: ", err)
		}
	}()

	<-a.sigQuit
	log.Infoln("Gracefully shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalln("Failed to shutdown the server gracefully: ", err)
	}

	log.Infoln("Server shutdown is successful")
}
