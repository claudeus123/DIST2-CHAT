package router

import (
	"time"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/claudeus123/DIST2-CHAT/ws"
)

var r *gin.Engine

func InitRouter(wsHandler *ws.Handler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// r.Use(middlewares.Validate)

	r.POST("/ws/createRoom", wsHandler.CreateRoom)
	r.POST("/ws/createChat", wsHandler.CreateChat)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/joinChat/:chatID", wsHandler.JoinChat)
	r.GET("/ws/getRooms", wsHandler.GetRooms)
	r.GET("/ws/getChats", wsHandler.GetChats)
	r.GET("/ws/getAvailableChats/:userId", wsHandler.GetAvailableChats)
	r.GET("/ws/getClients/:roomId", wsHandler.GetClients)
	r.GET("/ws/getChatClients/:chatID", wsHandler.GetChatClients)
}

func Initialize(wsHandler *ws.Handler){
	wsHandler.InitializeChats()
}

func Start(addr string) error {
	return r.Run(addr)
}