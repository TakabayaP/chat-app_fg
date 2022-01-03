package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "dbname=chat sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/chatrooms",handleChatroomSelect)
	mux.HandleFunc("/chat", chat)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWs(w, r)
	})

	go handleMessages()
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "chat server running!")
}
