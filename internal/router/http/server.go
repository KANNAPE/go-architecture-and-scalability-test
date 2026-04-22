package http

import (
	"encoding/json"
	"net/http"
	"time"

	"kannape.com/upfluence-test/internal/services/stream"
)

type Server struct {
	router *http.ServeMux

	streamRepo stream.IRepository
}

func NewServer(streamRepo stream.IRepository) *Server {
	server := &Server{
		router:     http.NewServeMux(),
		streamRepo: streamRepo,
	}
	
	server.router.HandleFunc("GET /analysis", server.analysis)

	return server
}

func (server *Server) Start() error {
	return http.ListenAndServe(":8080", server.router)
}

func (server *Server) analysis(w http.ResponseWriter, r *http.Request) {
	data, err := server.streamRepo.GetStream(time.Second * 5)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
