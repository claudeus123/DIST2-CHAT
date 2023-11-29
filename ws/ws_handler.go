// package ws

// type Handler struct {
// 	hub *Hub
// }

package ws

import (
	"net/http"

	"fmt"

	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/claudeus123/DIST2-CHAT/models"
	"strconv"
)

type Handler struct {
	hub *Hub
	db *gorm.DB
}

func NewHandler(h *Hub, db *gorm.DB) *Handler {
	return &Handler{
		hub: h,
		db: db,
	}
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateChatReq struct {
	ID   uint `json:"id"`
	User1ID uint `json:"user1_id"`
	User2ID uint `json:"user2_id"`
}


func (h *Handler) InitializeChats() {
	chats := make([]ChatRes, 0)

	// Consulta las salas desde la base de datos usando GORM
	var chatModels []models.Chat
	if err := h.db.Find(&chatModels).Error; err != nil {
		return
	}
	fmt.Println(chatModels)

	// Convierte los modelos de sala a RoomRes y agrega a la lista
	for _, chat := range chatModels {
		chats = append(chats, ChatRes{
			ID:   chat.ID,
			User1ID: chat.User1ID,
			User2ID: chat.User2ID,
		})

		chatIDInt := int(chat.ID)
		h.hub.Chats[strconv.Itoa(chatIDInt)] = &Chat{
			ID:      chat.ID,
			User1ID:    chat.User1ID,
			User2ID:    chat.User2ID,
			Clients: make(map[string]*Client),
		}
	}

	// c.JSON(http.StatusOK, chats)
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

func (h *Handler) JoinChat(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatID := c.Param("chatID")
	intChatID, err := strconv.Atoi(chatID)
	clientID := c.Query("userId")
	username := c.Query("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		ChatID:   uint(intChatID),
		Username: username,
	}

	m := &Message{
		Content:  "A new user has joined the room",
		ChatID:   uint(intChatID),
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}


type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ChatRes struct {
	ID       uint `json:"id"`
	User1ID  uint `json:"user1_id"`
	User2ID  uint `json:"user2_id"`
}

func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *Handler) GetChats(c *gin.Context) {
	chats := make([]ChatRes, 0)

	for _, r := range h.hub.Chats {
		chats = append(chats, ChatRes{
			ID:   r.ID,
			User1ID: r.User1ID,
			User2ID: r.User2ID,
		})
	}
	fmt.Println("chats")
	
	c.JSON(http.StatusOK, chats)
}

func (h *Handler) GetAvailableChats(c *gin.Context) {
	// chats := make([]ChatRes, 0)
	
	userId := c.Param("userId")
	// var chatModels []models.Chat
	// if err := h.db.Find(&chatModels).Error; err != nil {
	// 	return
	// }

	// for _, r := range h.hub.Chats {
	// 	chats = append(chats, ChatRes{
	// 		ID:   r.ID,
	// 		User1ID: r.User1ID,
	// 		User2ID: r.User2ID,
	// 	})
	// }
	// fmt.Println("chats")
	
	// c.JSON(http.StatusOK, chats)
	chats := make([]ChatRes, 0)

    // Obtener todos los chats desde la base de datos
    var chatModels []models.Chat
    if err := h.db.Find(&chatModels).Error; err != nil {
        // Manejar el error según tus necesidades
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Filtrar los chats para obtener solo aquellos en los que el usuario está involucrado
	for _, chat := range chatModels {
		userID, err := strconv.Atoi(userId)
		if err != nil {
			// Handle the error according to your needs
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if chat.User1ID == uint(userID) || chat.User2ID == uint(userID) {
			chats = append(chats, ChatRes{
				ID:      chat.ID,
				User1ID: chat.User1ID,
				User2ID: chat.User2ID,
			})
		}

	}

	fmt.Println("chats")

    // Responder con los chats disponibles para el usuario
    c.JSON(http.StatusOK, chats)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *gin.Context) {
	var clients []ClientRes
	roomId := c.Param("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}

// FALTA QUE AL HACER MATCH SE HAGA UN CREATE CHAT
// func (h *Handler) CreateRoom(c *gin.Context) {
// 	var req CreateRoomReq
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	h.hub.Rooms[req.ID] = &Room{
// 		ID:      req.ID,
// 		Name:    req.Name,
// 		Clients: make(map[string]*Client),
// 	}

// 	c.JSON(http.StatusOK, req)
// }