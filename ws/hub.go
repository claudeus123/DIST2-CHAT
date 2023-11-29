// package ws

// type Room struct {
// 	ID int `json:"id"`
// 	Name string `json:"name"`
// 	Clients map[string]*Client `json:"clients"`
// }

// type Hub struct {
// 	Rooms map[string]*Room `json:"rooms"`

// }

// func NewHub() *Hub {
// 	return &Hub{
// 		Rooms: make(map[string]*Room),
// 	}
// }

package ws

import(
	"strconv"
)

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Chat struct {
	ID      uint             `json:"id"`
	User1ID    uint             `json:"user1_id"`
	User2ID    uint             `json:"user2_id"`

	Clients map[string]*Client `json:"clients"`
}


type Hub struct {
	Rooms      map[string]*Room
	Chats 		map[string]*Chat
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}



func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Chats:		make(map[string]*Chat),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

// func (h *Hub) Run() {
// 	for {
// 		select {
// 		case cl := <-h.Register:
// 			if _, ok := h.Rooms[cl.RoomID]; ok {
// 				r := h.Rooms[cl.RoomID]

// 				if _, ok := r.Clients[cl.ID]; !ok {
// 					r.Clients[cl.ID] = cl
// 				}
// 			}
// 		case cl := <-h.Unregister:
// 			if _, ok := h.Rooms[cl.RoomID]; ok {
// 				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
// 					if len(h.Rooms[cl.RoomID].Clients) != 0 {
// 						h.Broadcast <- &Message{
// 							Content:  "user left the chat",
// 							RoomID:   cl.RoomID,
// 							Username: cl.Username,
// 						}
// 					}

// 					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
// 					close(cl.Message)
// 				}
// 			}

// 		case m := <-h.Broadcast:
// 			if _, ok := h.Rooms[m.RoomID]; ok {

// 				for _, cl := range h.Rooms[m.RoomID].Clients {
// 					cl.Message <- m
// 				}
// 			}
// 		}
// 	}
// }

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			chatID := strconv.Itoa(int(cl.ChatID))
			if _, ok := h.Chats[chatID]; ok {
				r := h.Chats[chatID]

				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
				}
			}
		case cl := <-h.Unregister:
			chatID := strconv.Itoa(int(cl.ChatID))
			if _, ok := h.Chats[chatID]; ok {
				if _, ok := h.Chats[chatID].Clients[cl.ID]; ok {
					if len(h.Chats[chatID].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  "user left the chat",
							ChatID:   cl.ChatID,
							Username: cl.Username,
						}
					}

					delete(h.Chats[chatID].Clients, cl.ID)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			chatID := strconv.Itoa(int(m.ChatID))
			if chat, ok := h.Chats[chatID]; ok {
				for _, cl := range chat.Clients {
					cl.Message <- m
				}
			}
		}
	}
}