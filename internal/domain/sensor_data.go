package domain

type SensorData struct {
	RoomID      string  `json:"roomID"`
	Timestamp   int64   `json:"timestamp"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}
