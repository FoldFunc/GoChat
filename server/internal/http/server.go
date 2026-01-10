package http

import (
	"net/http"
	"time"

)
type Server struct {
	httpServer *http.Server
}
func New(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: addr,
			Handler: handler,
			ReadTimeout: 5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout: 120 *time.Second,
		},
	}
}
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
