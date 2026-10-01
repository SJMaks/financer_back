package httpserver

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	http *http.Server
}

func New(engine *gin.Engine, addr string) *Server {
	return &Server{
		http: &http.Server{
			Addr:    addr,
			Handler: engine,
		},
	}
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
