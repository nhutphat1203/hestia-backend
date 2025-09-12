package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv            string
	ServerPort        string
	MQTTBroker        string
	MQTTTopic         string
	MQTTUser          string
	MQTTPass          string
	DatabaseDSN       string
	LogLevel          string
	TopicQoS          byte
	WorkerCount       int
	IDENTIFY_PROPERTY string
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	topicQoS, err := strconv.Atoi(getEnv("TOPIC_QOS", "1"))

	if err != nil {
		return nil, err
	}

	workerCount, err := strconv.Atoi(getEnv("WORKER_COUNT", "4"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppEnv:            getEnv("APP_ENV", "development"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		MQTTBroker:        getEnv("MQTT_BROKER", ""),
		MQTTTopic:         getEnv("MQTT_TOPIC", ""),
		MQTTUser:          getEnv("MQTT_USERNAME", ""),
		MQTTPass:          getEnv("MQTT_PASSWORD", ""),
		DatabaseDSN:       getEnv("DATABASE_DSN", ""),
		LogLevel:          getEnv("LOG_LEVEL", "debug"),
		TopicQoS:          byte(topicQoS),
		WorkerCount:       workerCount,
		IDENTIFY_PROPERTY: getEnv("IDENTIFY_PROPERTY", "roomID"),
	}

	return cfg, nil
}
