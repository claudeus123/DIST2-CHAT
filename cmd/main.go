package main

import (
	"log"
	"github.com/claudeus123/DIST2-CHAT/db"
)


func main() {
	_, err := db.ConnectDb()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	log.Println("Conexión exitosa a la base de datos")
}