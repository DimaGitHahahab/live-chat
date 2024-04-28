package server

import (
	"context"
	"net/http"
	"strings"

	"chat-service/internal/handler"
	"chat-service/internal/service"
	"chat-service/pkg/config"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config) *Server {
	chatService := service.NewChat(strings.Split(cfg.KafkaBrokers, ","), cfg.KafkaTopic, cfg.StorageServiceUrl)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.HttpPort,
			Handler: handler.NewChatHandler(chatService),
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
