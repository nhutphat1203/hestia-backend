package websocket

import (
	"sync"

	"github.com/nhutphat1203/hestia-backend/internal/domain"
)

type Hub struct {
	Rooms map[string]*domain.Room
	lock  sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*domain.Room),
	}
}

func (h *Hub) GetOrCreateRoom(roomID string) *domain.Room {
	h.lock.Lock()
	defer h.lock.Unlock()
	room, exists := h.Rooms[roomID]
	if !exists {
		room = domain.NewRoom(roomID)
		h.Rooms[roomID] = room
	}
	return room
}

func (h *Hub) RemoveRoomIfEmpty(roomID string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if room, ok := h.Rooms[roomID]; ok && !room.HasClients() {
		delete(h.Rooms, roomID)
	}
}
