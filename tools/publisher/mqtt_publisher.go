package main

import (
	"flag"
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

	broker := flag.String("broker", "example.com", "Địa chỉ MQTT broker")
	topic := flag.String("topic", "test/topic", "Topic để publish dữ liệu")
	interval := flag.Int("interval", 1, "Khoảng thời gian giữa các message (giây)")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)
	opts := mqtt.NewClientOptions().AddBroker(*broker)
	opts.SetClientID(fmt.Sprintf("publisher-%d", time.Now().UnixNano()))
	// Nếu broker yêu cầu username/password
	if cfg.MQTTUser != "" {
		opts.SetUsername(cfg.MQTTUser)
	}
	if cfg.MQTTPass != "" {
		opts.SetPassword(cfg.MQTTPass)
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	fmt.Printf("🚀 Connected to broker %s\n", *broker)
	fmt.Printf("📤 Publishing to topic '%s' every %d second(s)\n", *topic, *interval)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				payload := fmt.Sprintf(`
				{
					"roomID": "0001",
					"temperature": %d,
					"humidity": %d
				}`, rand.Intn(1000), rand.Intn(100))
				token := client.Publish(*topic, 0, false, payload)
				token.Wait()
				fmt.Printf("📤 Published: %s\n", payload)
				time.Sleep(time.Duration(*interval) * time.Second)
			}
		}
	}()

	<-sigc
	fmt.Println("\n🛑 Interrupted! Stopping publisher...")
	close(done)
	time.Sleep(1 * time.Second)
}
