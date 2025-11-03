package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nhutphat1203/hestia-backend/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Cannot load config: %v", err)
	}

	// Tạo MQTT options từ config
	opts := mqtt.NewClientOptions().
		AddBroker("127.0.0.1:1883").
		SetClientID(fmt.Sprintf("publisher-%d", time.Now().UnixNano()))

	// Nếu có username/password
	if cfg.MQTTUser != "" {
		opts.SetUsername(cfg.MQTTUser)
	}
	if cfg.MQTTPass != "" {
		opts.SetPassword(cfg.MQTTPass)
	}

	// Nếu MQTT_SSL = true, bật TLS
	if cfg.MQTTSSL {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // có thể đổi sang false nếu bạn có chứng chỉ hợp lệ
			ClientAuth:         tls.NoClientCert,
		}
		opts.SetTLSConfig(tlsConfig)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ Failed to connect MQTT: %v", token.Error())
	}
	defer client.Disconnect(250)

	fmt.Printf("🚀 Connected to broker %s\n", cfg.MQTTBroker)
	fmt.Printf("📤 Publishing to topic '%s'\n", cfg.MQTTTopic)

	// Bắt tín hiệu Ctrl+C
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				payload := fmt.Sprintf(`{
					"schemaVersion": 1,
					"roomId": "kit-01",
					"type": "env",
					"ts": %d,
					"measure": {
						"t": %.2f,
						"h": %.2f,
						"p": %.2f,
						"lux": %d
					},
					"score": %d,
					"state": "OK",
					"meta": {
						"seq": %d,
						"source": "gw-01@fw1.3.2",
						"units": { "t": "C", "h": "%%", "p": "hPa", "lux": "lux" }
					}
				}`,
					time.Now().UnixMilli(),
					20+10*rand.Float64(),  // t
					40+60*rand.Float64(),  // h
					980+40*rand.Float64(), // p
					rand.Intn(200),        // lux
					rand.Intn(100),        // score
					rand.Intn(10000),      // seq
				)

				token := client.Publish(cfg.MQTTTopic, 0, false, payload)
				token.Wait()
				fmt.Printf("📤 Published: %s\n", payload)
				time.Sleep(2 * time.Second) // có thể thêm ENV để điều khiển
			}
		}
	}()

	<-sigc
	fmt.Println("\n🛑 Interrupted! Stopping publisher...")
	close(done)
	time.Sleep(1 * time.Second)
}
