package jobs

import (
	"github.com/nhutphat1203/hestia-backend/internal/infrastructure/websocket"
)

type BroadcastJob struct {
	Payload      []byte
	RoomID       string
	websocketHub *websocket.Hub
}

func NewBroadcastJob(payload []byte, roomID string, websocketHub *websocket.Hub) *BroadcastJob {
	return &BroadcastJob{
		Payload:      payload,
		RoomID:       roomID,
		websocketHub: websocketHub,
	}
}

func (j *BroadcastJob) Execute() error {
	room := j.websocketHub.GetOrCreateRoom(j.RoomID)
	room.Broadcast(j.Payload)
	return nil
}
