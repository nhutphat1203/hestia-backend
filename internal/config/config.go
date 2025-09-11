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
	TopicQoS    string
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		MQTTBroker:  getEnv("MQTT_BROKER", ""),
		MQTTTopic:   getEnv("MQTT_TOPIC", ""),
		MQTTUser:    getEnv("MQTT_USERNAME", ""),
		MQTTPass:    getEnv("MQTT_PASSWORD", ""),
		DatabaseDSN: getEnv("DATABASE_DSN", ""),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
		TopicQoS:    getEnv("TOPIC_QOS", "1"),
	}

	return cfg, nil
}
