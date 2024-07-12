package services

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer http.Server
}

type HttpServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

func NewServer(cfg HttpServerConfig, handler http.Handler) *Server {
	return &Server{
		httpServer: http.Server{
			Addr:         cfg.Host + ":" + cfg.Port,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
