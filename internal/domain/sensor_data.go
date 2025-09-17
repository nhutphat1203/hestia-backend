package domain

type Measure struct {
	T   float64 `json:"t"`
	H   float64 `json:"h"`
	P   float64 `json:"p"`
	Lux float64 `json:"lux"`
}

type Units struct {
	T   string `json:"t"`
	H   string `json:"h"`
	P   string `json:"p"`
	Lux string `json:"lux"`
}

type Meta struct {
	Seq    int    `json:"seq"`
	Source string `json:"source"`
	Units  Units  `json:"units"`
}

type SensorData struct {
	SchemaVersion int     `json:"schemaVersion"`
	RoomID        string  `json:"roomId"`
	Type          string  `json:"type"`
	Ts            int64   `json:"ts"`
	Measure       Measure `json:"measure"`
	Score         int     `json:"score"`
	State         string  `json:"state"`
	Meta          Meta    `json:"meta"`
}
