package mqtt_client

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nhutphat1203/hestia-backend/internal/config"
	"github.com/nhutphat1203/hestia-backend/pkg/logger"
)

type Client interface {
	Connect() error
	Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error
	Publish(topic string, qos byte, retained bool, payload interface{}) error
	Unsubscribe(topics ...string) error
}

type MQTTClient struct {
	client mqtt.Client
	cfg    *config.Config
	logger *logger.Logger
}

func New(cfg *config.Config, logger *logger.Logger) *MQTTClient {
	opts := mqtt.NewClientOptions()
	//opts.AddBroker(cfg.MQTT.Broker)
	//opts.SetClientID(cfg.MQTT.ClientID)

	client := mqtt.NewClient(opts)

	return &MQTTClient{
		client: client,
		cfg:    cfg,
		logger: logger,
	}
}

func (m *MQTTClient) Connect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	//m.logger.Info("MQTT connected to " + m.cfg.MQTT.Broker)
	return nil
}

func (m *MQTTClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	token := m.client.Subscribe(topic, qos, callback)
	token.Wait()
	return token.Error()
}

func (m *MQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	token := m.client.Publish(topic, qos, retained, payload)
	token.Wait()
	return token.Error()
}

func (m *MQTTClient) Unsubscribe(topics ...string) error {
	token := m.client.Unsubscribe(topics...)
	token.Wait()
	return token.Error()
}
