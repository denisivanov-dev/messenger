package hub

type Hub struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan RoomMessage
	JoinRoom   chan joinReq

	clients     map[*Client]bool
	userClients map[string]*Client
	rooms       map[string]map[*Client]bool
}

func NewHub() *Hub {
	h := &Hub{
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		Broadcast:    make(chan RoomMessage, 256),
		JoinRoom:     make(chan joinReq, 64),
		clients:      make(map[*Client]bool),
		userClients:  make(map[string]*Client),
		rooms:        make(map[string]map[*Client]bool),
	}
	go h.Run()
	return h
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.handleRegister(c)

		case c := <-h.Unregister:
			h.handleUnregister(c)

		// Add client to a new room dynamically
		case jr := <-h.JoinRoom:
			h.handleJoinRoom(jr)

		// Route a message to a room
		case msg := <-h.Broadcast:
			h.handleBroadcast(msg)
		}
	}
}