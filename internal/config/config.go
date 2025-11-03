package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv             string
	ServerAddress      string
	MQTTBroker         string
	MQTTTopic          string
	MQTTUser           string
	MQTTPass           string
	LogLevel           string
	TopicQoS           byte
	WorkerCount        int
	Identify_property  string
	InfluxDBURL        string
	InfluxDBAdminToken string
	InfluxDBOrg        string
	InfluxDBBucket     string
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
	ServerIdleTimeout  time.Duration
	StaticToken        string
	MQTTSSL            bool
	JWTSecret          string
	JWTExpiration      time.Duration
	DbHost             string
	DbPort             int
	DbUser             string
	DbPassword         string
	DbName             string
	AdminAcc           string
	AdminPwd           string
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key, def string) (int, error) {
	valStr := getEnv(key, def)
	return strconv.Atoi(valStr)
}

func getEnvDuration(key, def string) (time.Duration, error) {
	valStr := getEnv(key, def)
	return time.ParseDuration(valStr)
}

func getBool(key, def string) (bool, error) {
	val := getEnv(key, def)
	return strconv.ParseBool(val)
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	topicQoS, err := getEnvInt("TOPIC_QOS", "1")
	if err != nil {
		return nil, err
	}

	workerCount, err := getEnvInt("WORKER_COUNT", "4")
	if err != nil {
		return nil, err
	}

	// Timeout mặc định
	readTimeout, err := getEnvDuration("SERVER_READ_TIMEOUT", "10s")
	if err != nil {
		return nil, err
	}
	writeTimeout, err := getEnvDuration("SERVER_WRITE_TIMEOUT", "10s")
	if err != nil {
		return nil, err
	}
	idleTimeout, err := getEnvDuration("SERVER_IDLE_TIMEOUT", "60s")
	if err != nil {
		return nil, err
	}

	jwtExpiration, err := getEnvDuration("JWT_EXPIRATION", "60m")
	if err != nil {
		return nil, err
	}

	mqttSSL, err := getBool("MQTT_SSL", "false")

	if err != nil {
		return nil, err
	}

	dbPort, err := getEnvInt("DB_PORT", "5432")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		ServerAddress:      getEnv("SERVER_ADDRESS", "0.0.0.0:8080"),
		MQTTBroker:         getEnv("MQTT_BROKER", ""),
		MQTTTopic:          getEnv("MQTT_TOPIC", ""),
		MQTTUser:           getEnv("MQTT_USERNAME", ""),
		MQTTPass:           getEnv("MQTT_PASSWORD", ""),
		LogLevel:           getEnv("LOG_LEVEL", "debug"),
		TopicQoS:           byte(topicQoS),
		WorkerCount:        workerCount,
		Identify_property:  getEnv("IDENTIFY_PROPERTY", "roomID"),
		InfluxDBURL:        getEnv("INFLUXDB_URL", ""),
		InfluxDBAdminToken: getEnv("INFLUXDB_ADMIN_TOKEN", ""),
		InfluxDBOrg:        getEnv("INFLUXDB_ORG", ""),
		InfluxDBBucket:     getEnv("INFLUXDB_BUCKET", ""),
		ServerReadTimeout:  readTimeout,
		ServerWriteTimeout: writeTimeout,
		ServerIdleTimeout:  idleTimeout,
		StaticToken:        getEnv("STATIC_TOKEN", ""),
		MQTTSSL:            mqttSSL,
		JWTSecret:          getEnv("JWT_SECRET", "very-secret-key"),
		JWTExpiration:      jwtExpiration,
		DbHost:             getEnv("DB_HOST", "localhost"),
		DbPort:             dbPort,
		DbUser:             getEnv("DB_USER", "admin"),
		DbPassword:         getEnv("DB_PASSWORD", "secretpassword"),
		DbName:             getEnv("DB_NAME", "hestia"),
		AdminAcc:           getEnv("ADMIN_ACC", "admin"),
		AdminPwd:           getEnv("ADMIN_PWD", "123456"),
	}

	return cfg, nil
}
