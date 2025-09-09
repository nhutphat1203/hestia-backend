package server

import (
	"context"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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

func (s *Server) Start() error {
	s.logger.Info("Starting server...")

	// Connect MQTT trước
	if err := s.mqttClient.Connect(); err != nil {
		return err
	}
	s.logger.Info("MQTT client connected")

	err := s.mqttClient.Subscribe("#", 1, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("Received message on topic:", msg.Topic())
		fmt.Println("Payload:", string(msg.Payload()))
	})
	if err != nil {
		fmt.Println("Subscribe failed:", err)
	}

	errCh := make(chan error, 1)
	go func() {
		if err := s.httpServer.Start(); err != nil {
			errCh <- err
		}
	}()

	err = <-errCh

	return err
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Stop(ctx); err != nil {
		s.logger.Error("Error shutting down HTTP server: " + err.Error())
	}

	s.mqttClient.Disconnect()
	s.logger.Info("Server stopped")
}
