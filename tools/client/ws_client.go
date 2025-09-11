package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ws_client <server_host> <roomID>")
		return
	}

	host := os.Args[1] // ví dụ: localhost:8080
	roomID := os.Args[2]

	u := url.URL{Scheme: "ws", Host: host, Path: "/ws/v1/env", RawQuery: "room=" + roomID}
	fmt.Println("Connecting to", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer c.Close()

	fmt.Println("Connected to server. Listening for messages...")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	// Goroutine nhận message từ server
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			fmt.Printf("Received: %s\n", message)
		}
	}()

	// Wait until user interrupts
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("Interrupt received, closing connection")
			// Gửi close message
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close error:", err)
			}
			return
		}
	}
}
