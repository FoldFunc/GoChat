package http

import (
	"log"
	"net/http"
	"time"
)
type Server struct {
	httpServer *http.Server
	certFile string
	keyFile string
}
func New(addr string, handler http.Handler, certFile, keyFile string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: addr,
			Handler: handler,
			ReadTimeout: 5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout: 120 *time.Second,
		},
		certFile: certFile,
		keyFile: keyFile,
	}
}
func (s *Server) Start() error {
	if s.certFile != "" && s.keyFile != "" {
		log.Println("Running an https server")
		return s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile)
	}
	log.Println("Running an http server")
	return s.httpServer.ListenAndServe()
}
