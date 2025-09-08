package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	ServerPort  string
	MQTTBroker  string
	MQTTTopic   string
	MQTTUser    string
	MQTTPass    string
	DatabaseDSN string
	LogLevel    string
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func LoadConfig() (*Config, error) {
	// Nếu không phải production thì load .env (yên tâm ignore lỗi khi file không tồn tại)
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load()
	}

	cfg := &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		MQTTBroker:  getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		MQTTTopic:   getEnv("MQTT_TOPIC", "devices/+/telemetry"),
		MQTTUser:    getEnv("MQTT_USERNAME", ""),
		MQTTPass:    getEnv("MQTT_PASSWORD", ""),
		DatabaseDSN: getEnv("DATABASE_DSN", ""),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
	// TODO: validate required fields if needed
	return cfg, nil
}
