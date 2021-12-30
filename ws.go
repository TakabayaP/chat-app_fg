package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Chat)
var upgrader = websocket.Upgrader{}

func handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients[conn] = true

	for {
		var chat Chat
		err := conn.ReadJSON(&chat)
		fmt.Println(chat)
		if err != nil {
			delete(clients, conn)
			panic(err)
		}
		broadcast <- chat
	}
}
func handleMessages() {
	for {
		message := <-broadcast
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
