package chat

import (
	"encoding/json"
	"regexp"
)

type hub struct {
	rooms      map[string]map[*Connection]bool
	broadcast  chan message
	register   chan Subscription
	unregister chan Subscription
}

type message struct {
	data []byte
	room string
}

type Chat struct {
	Name        string `json:"name"`
	Message     string `json:"message"`
	Client      string `json:"client"`
	CreatedAt   string `json:"createdAt"`
	IsRead      bool   `json:"isRead"`
	MessageType string `json:"messageType"`
	ChatRoom    string `json:"chatRoom"`
}

func (h *hub) GetRooms() map[string]map[*Connection]bool {
	return h.rooms
}

func NewHub() *hub {
	return &hub{
		broadcast:  make(chan message),
		register:   make(chan Subscription),
		unregister: make(chan Subscription),
		rooms:      make(map[string]map[*Connection]bool),
	}
}

func (h *hub) Run() {
	for {
		select {
		case s := <-h.register:
			Connections := h.rooms[s.room]

			if Connections == nil {
				Connections = make(map[*Connection]bool)
				h.rooms[s.room] = Connections
			}

			h.rooms[s.room][s.conn] = true

		case s := <-h.unregister:
			Connections := h.rooms[s.room]

			if Connections != nil {
				if _, ok := Connections[s.conn]; ok {
					delete(Connections, s.conn)
					close(s.conn.send)

					if len(Connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}

		case m := <-h.broadcast:
			Connections := h.rooms[m.room]

			for c := range Connections {
				select {
				case c.send <- m.data:
					// chatbot logic
					b := Chat{}
					json.Unmarshal(m.data, &b)
					r := regexp.MustCompile("\\(\\/.*?\\)")
					command := r.FindString(b.Message)
					if command == "" {
						continue
					}

					chatbot(c.send, command)

				default:
					close(c.send)
					delete(Connections, c)

					if len(Connections) == 0 {
						delete(h.rooms, m.room)
					}
				}
			}
		}
	}
}
