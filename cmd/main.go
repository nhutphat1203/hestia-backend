package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nhutphat1203/hestia-backend/cmd/server"
	"github.com/nhutphat1203/hestia-backend/internal/config"
	mqtt_client "github.com/nhutphat1203/hestia-backend/internal/infrastructure/mqtt"
	http_server "github.com/nhutphat1203/hestia-backend/internal/interfaces/http"
	app_logger "github.com/nhutphat1203/hestia-backend/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config loaded")

	logger := app_logger.New(cfg)
	logger.Info("Logger initialized")

	mqttClient := mqtt_client.New(cfg, logger)

	httpServer := http_server.New(cfg, logger)

	server := server.New(cfg, logger, httpServer, mqttClient)

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
