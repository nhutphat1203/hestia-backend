package websocket

import (
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

	// Goroutine gửi dữ liệu cho client
	go func() {
		defer func() {
			room.RemoveClient(client.ID)
			h.RemoveRoomIfEmpty(roomID)
			conn.Close()
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

func (h *Hub) ServeWS_Token_Query(c *gin.Context, authenticator domain.Authenticator) {
	roomID := c.Query("room")
	if roomID == "" {
		response.SendError(c, http.StatusBadRequest, "room id required", errorf.Validation)
		return
	}

	// LẤY TOKEN TỪ QUERY
	token := c.Query("token")
	if token == "" {
		response.SendError(c, http.StatusUnauthorized, "token required", errorf.Unauthorized)
		return
	}

	ok, err := authenticator.Authenticate(token)
	if !ok || err != nil {
		response.SendError(c, errorf.HttpStatus(errorf.InvalidToken), errorf.Message(errorf.InvalidToken), errorf.InvalidToken)
		return
	}

	// Upgrade sau khi xác thực
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "failed to upgrade to websocket", errorf.Internal)
		return
	}

	client := domain.NewClient(roomID + "-" + conn.RemoteAddr().String())

	room := h.GetOrCreateRoom(roomID)
	room.AddClient(client)

	// gửi dữ liệu
	go func() {
		defer func() {
			room.RemoveClient(client.ID)
			h.RemoveRoomIfEmpty(roomID)
			conn.Close()
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
