package server

import (
	"context"
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nhutphat1203/hestia-backend/internal/config"
	mqtt_client "github.com/nhutphat1203/hestia-backend/internal/infrastructure/mqtt"
	"github.com/nhutphat1203/hestia-backend/internal/infrastructure/websocket"
	http_server "github.com/nhutphat1203/hestia-backend/internal/interfaces/http"
	"github.com/nhutphat1203/hestia-backend/internal/jobs"
	"github.com/nhutphat1203/hestia-backend/pkg/logger"
	"github.com/nhutphat1203/hestia-backend/pkg/worker"
)

type Server struct {
	cfg          *config.Config
	logger       *logger.Logger
	httpServer   *http_server.HTTPServer
	mqttClient   mqtt_client.Client
	websocketHub *websocket.Hub
	dispatcher   *worker.Dispatcher
}

func New(cfg *config.Config,
	logger *logger.Logger,
	httpServer *http_server.HTTPServer,
	mqttClient mqtt_client.Client,
	websocketHub *websocket.Hub,
	dispatcher *worker.Dispatcher) *Server {
	return &Server{
		cfg:          cfg,
		logger:       logger,
		httpServer:   httpServer,
		mqttClient:   mqttClient,
		websocketHub: websocketHub,
		dispatcher:   dispatcher,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting server...")

	go s.dispatcher.Run()

	// Connect MQTT trước
	if err := s.mqttClient.Connect(); err != nil {
		return err
	}
	s.logger.Info("MQTT client connected")

	err := s.mqttClient.Subscribe(s.cfg.MQTTTopic, s.cfg.TopicQoS, func(client mqtt.Client, msg mqtt.Message) {
		s.logger.Debug("Received message on topic: " + msg.Topic())
		s.logger.Debug("Payload: " + string(msg.Payload()))
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			return
		}

		roomID, ok := payload[s.cfg.IDENTIFY_PROPERTY].(string)
		if !ok {
			return
		}

		s.dispatcher.JobQueue <- jobs.NewBroadcastJob(msg.Payload(), roomID, s.websocketHub)
	})

	if err != nil {
		s.logger.Error("Subscribe failed: " + err.Error())
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
