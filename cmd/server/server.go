package server

import (
	"github.com/nhutphat1203/hestia-backend/internal/config"
	mqtt_client "github.com/nhutphat1203/hestia-backend/internal/infrastructure/mqtt"
	http_server "github.com/nhutphat1203/hestia-backend/internal/interfaces/http"
	"github.com/nhutphat1203/hestia-backend/pkg/logger"
)

type Server struct {
	cfg        *config.Config
	logger     *logger.Logger
	httpServer *http_server.HTTPServer
	mqttClient mqtt_client.Client
}

func New(cfg *config.Config, logger *logger.Logger, httpServer *http_server.HTTPServer, mqttClient mqtt_client.Client) *Server {
	return &Server{
		cfg:        cfg,
		logger:     logger,
		httpServer: httpServer,
		mqttClient: mqttClient,
	}
}
