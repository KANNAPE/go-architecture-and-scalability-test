package http

import (
	"net/http"

	"kannape.com/upfluence-test/internal/services/stream"
)

type Server struct {
	router *http.ServeMux

	streamRepo stream.IRepository
}

func NewServer(streamRepo stream.IRepository) *Server {
	return &Server{
		router:     http.NewServeMux(),
		streamRepo: streamRepo,
	}
}

func (server *Server) Start() error {
	services := services{
		streamService: stream.NewService(server.streamRepo),
	}

	registerRoutes(server, services)

	return http.ListenAndServe(":8080", server.router)
}
