package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nhutphat1203/hestia-backend/cmd/server"
	"github.com/nhutphat1203/hestia-backend/internal/config"
	influxdb_client "github.com/nhutphat1203/hestia-backend/internal/infrastructure/influxdb"
	mqtt_client "github.com/nhutphat1203/hestia-backend/internal/infrastructure/mqtt"
	"github.com/nhutphat1203/hestia-backend/internal/infrastructure/websocket"
	http_server "github.com/nhutphat1203/hestia-backend/internal/interfaces/http/http_server"
	repository "github.com/nhutphat1203/hestia-backend/internal/repositories"
	service "github.com/nhutphat1203/hestia-backend/internal/services"
	app_logger "github.com/nhutphat1203/hestia-backend/pkg/logger"
	"github.com/nhutphat1203/hestia-backend/pkg/worker"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg)

	log.Println("Config loaded")

	logger := app_logger.New(cfg)
	logger.Info("Logger initialized")

	mqttClient := mqtt_client.New(cfg, logger)

	websocketHub := websocket.NewHub()

	httpServer := http_server.New(cfg, logger, websocketHub)

	jobQueue := make(chan worker.Job, 100)

	dispatcher := worker.NewDispatcher(cfg.WorkerCount, jobQueue)

	influxDBClient := influxdb_client.NewInfluxDBClient(cfg, logger)
	defer influxDBClient.Close()

	measurementRepo := repository.NewInfluxDBRepo(influxDBClient)

	measurementService := service.NewMeasurementService(measurementRepo)

	server := server.New(cfg, logger, httpServer, mqttClient, websocketHub, dispatcher, measurementService)

	go func() {
		if err := server.Start(); err != nil {
			logger.Error("Server error: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down gracefully...")
	server.Stop()
}
