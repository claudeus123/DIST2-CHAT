package main

import (
	"log"
	"github.com/claudeus123/DIST2-CHAT/database"
	"github.com/claudeus123/DIST2-CHAT/ws"
	"github.com/claudeus123/DIST2-CHAT/router"
)


func main() {
	database.ConnectDb()
	log.Println("Se conecto a la base de datos")

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub, database.DB)
	go hub.Run()


	router.InitRouter(wsHandler)
	router.Initialize(wsHandler)
	router.Start("0.0.0.0:8080")
}