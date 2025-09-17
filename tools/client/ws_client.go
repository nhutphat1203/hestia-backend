package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: ws_client <server_host:port> <roomID> <token>")
		fmt.Println("Example: ws_client localhost:80 1001 supersecrettoken")
		return
	}

	host := os.Args[1]   // ví dụ: localhost:80
	roomID := os.Args[2] // ví dụ: 1001
	token := os.Args[3]  // token auth

	// Tạo URL WebSocket
	u := url.URL{
		Scheme:   "ws",
		Host:     host,
		Path:     "/ws/v1/env",
		RawQuery: "room=" + roomID,
	}
	fmt.Println("Connecting to", u.String())

	// Thêm header Authorization
	header := http.Header{}
	header.Add("Authorization", "Bearer "+token)

	// Kết nối WebSocket
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer c.Close()
	fmt.Println("Connected to server. Listening for messages...")

	// Bắt tín hiệu Ctrl+C để đóng kết nối
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

	// Vòng lặp chính: đợi tín hiệu Ctrl+C hoặc connection lỗi
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("Interrupt received, closing connection...")
			// Gửi close message
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close error:", err)
			}
			return
		}
	}
}
