package server

import (
	"context"
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nhutphat1203/hestia-backend/internal/config"
	"github.com/nhutphat1203/hestia-backend/internal/domain"
	mqtt_client "github.com/nhutphat1203/hestia-backend/internal/infrastructure/mqtt"
	"github.com/nhutphat1203/hestia-backend/internal/infrastructure/websocket"
	http_server "github.com/nhutphat1203/hestia-backend/internal/interfaces/http"
	"github.com/nhutphat1203/hestia-backend/internal/jobs"
	service "github.com/nhutphat1203/hestia-backend/internal/services"
	"github.com/nhutphat1203/hestia-backend/pkg/logger"
	"github.com/nhutphat1203/hestia-backend/pkg/worker"
)

type Server struct {
	cfg                *config.Config
	logger             *logger.Logger
	httpServer         *http_server.HTTPServer
	mqttClient         mqtt_client.Client
	websocketHub       *websocket.Hub
	dispatcher         *worker.Dispatcher
	measurementService *service.MeasurementService
}

func New(cfg *config.Config,
	logger *logger.Logger,
	httpServer *http_server.HTTPServer,
	mqttClient mqtt_client.Client,
	websocketHub *websocket.Hub,
	dispatcher *worker.Dispatcher,
	measurementService *service.MeasurementService,
) *Server {
	return &Server{
		cfg:                cfg,
		logger:             logger,
		httpServer:         httpServer,
		mqttClient:         mqttClient,
		websocketHub:       websocketHub,
		dispatcher:         dispatcher,
		measurementService: measurementService,
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

		var sensorData domain.SensorData
		if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
			s.logger.Error("Failed to unmarshal sensor data: " + err.Error())
			return
		}

		s.dispatcher.JobQueue <- jobs.NewSaveMeasurementJob(sensorData, s.measurementService)

		s.dispatcher.JobQueue <- jobs.NewBroadcastJob(msg.Payload(), sensorData.RoomID, s.websocketHub)
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
