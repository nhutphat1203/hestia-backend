package main

import (
	"log"

	"github.com/nhutphat1203/hestia-backend/cmd/server"
	"github.com/nhutphat1203/hestia-backend/internal/config"
	mqtt_client "github.com/nhutphat1203/hestia-backend/internal/infrastructure/mqtt"
	http_server "github.com/nhutphat1203/hestia-backend/internal/interfaces/http"
	"github.com/nhutphat1203/hestia-backend/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Config loaded: %+v", cfg)

	logger := logger.New(cfg)
	logger.Info("Logger initialized")

	mqttClient := mqtt_client.New(cfg, logger)

	httpServer := http_server.New(cfg, logger)

	server := server.New(cfg, logger, httpServer, mqttClient)

	if err := server.Start(); err != nil {
		logger.Error("Failed to start HTTP server: " + err.Error())
	}

}
