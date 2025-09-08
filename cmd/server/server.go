package server

import (
	"github.com/nhutphat1203/hestia-backend/internal/config"
	http_server "github.com/nhutphat1203/hestia-backend/internal/interfaces/http"
	"github.com/nhutphat1203/hestia-backend/pkg/logger"
)

type Server struct {
	cfg        *config.Config
	logger     *logger.Logger
	httpServer *http_server.HTTPServer
}

func New(cfg *config.Config, logger *logger.Logger, httpServer *http_server.HTTPServer) *Server {
	return &Server{
		cfg:        cfg,
		logger:     logger,
		httpServer: httpServer,
	}
}
