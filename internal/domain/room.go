package domain

type Client struct {
	ID     string
	SendCh chan []byte
}

type Room struct {
	ID      string
	Clients map[string]*Client
}

func NewClient(id string) *Client {
	return &Client{
		ID:     id,
		SendCh: make(chan []byte, 256), // buffer size 256
	}
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Clients: make(map[string]*Client),
	}
}

func (r *Room) AddClient(client *Client) {
	r.Clients[client.ID] = client
}

func (r *Room) RemoveClient(clientID string) {
	delete(r.Clients, clientID)
}

func (r *Room) Broadcast(message []byte) {
	for _, client := range r.Clients {
		select {
		case client.SendCh <- message:
		default: // nếu buffer đầy thì bỏ qua hoặc remove client
			r.RemoveClient(client.ID)
		}
	}
}

func (r *Room) HasClients() bool {
	return len(r.Clients) > 0
}
