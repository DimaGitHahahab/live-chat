package server

import (
	"context"
	"net"
	"net/http"

	"storage-service/internal/handler"
	"storage-service/internal/service"
)

type Server struct {
	httpServer *http.Server
}

func New(storage *service.Storage, port string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    net.JoinHostPort("", port),
			Handler: handler.NewChatHandler(storage),
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
