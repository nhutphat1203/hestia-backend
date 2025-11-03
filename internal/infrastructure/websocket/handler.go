package websocket

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nhutphat1203/hestia-backend/internal/domain"
	"github.com/nhutphat1203/hestia-backend/pkg/errorf"
	"github.com/nhutphat1203/hestia-backend/pkg/response"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Hub) ServeWS(c *gin.Context) {
	roomID := c.Query("room")
	if roomID == "" {
		response.SendError(c, http.StatusBadRequest, "room id required", errorf.Validation)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "failed to upgrade to websocket", errorf.Internal)
		return
	}

	client := domain.NewClient(roomID + "-" + conn.RemoteAddr().String())

	room := h.GetOrCreateRoom(roomID)
	room.AddClient(client)

	fmt.Printf("Client %s connected to room %s\n", client.ID, roomID)

	// Goroutine gửi dữ liệu cho client
	go func() {
		defer func() {
			room.RemoveClient(client.ID)
			h.RemoveRoomIfEmpty(roomID)
			conn.Close()
			fmt.Printf("Client %s disconnected\n", client.ID)
		}()
		for msg := range client.SendCh {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				room.RemoveClient(client.ID)
				h.RemoveRoomIfEmpty(roomID)
				return
			}
		}
	}()
}
