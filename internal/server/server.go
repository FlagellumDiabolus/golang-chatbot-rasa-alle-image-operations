package server

import (
	"net/http"

	"chatbot-ai/internal/handlers"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (s *Server) SetupRoutes() {
	s.mux.HandleFunc("/chat", handlers.ChatHandler)
	s.mux.HandleFunc("/save", handlers.SaveImageHandler)
	s.mux.HandleFunc("/retrieve", handlers.RetrieveImageHandler)
	s.mux.HandleFunc("/all", handlers.ListImagesHandler)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
