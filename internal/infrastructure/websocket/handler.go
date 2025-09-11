package websocket

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nhutphat1203/hestia-backend/internal/domain"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Hub) ServeWS(c *gin.Context) {
	roomID := c.Query("room")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room id required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := domain.NewClient(roomID + "-" + conn.RemoteAddr().String())

	room := h.GetOrCreateRoom(roomID)
	room.AddClient(client)

	fmt.Printf("Client %s connected to room %s\n", client.ID, roomID)

	// Goroutine gửi dữ liệu cho client
	go func() {
		for msg := range client.SendCh {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				room.RemoveClient(client.ID)
				h.RemoveRoomIfEmpty(roomID)
				return
			}
		}
	}()
}
