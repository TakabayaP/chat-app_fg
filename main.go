package chat

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type ChatroomDatabase interface {
	Chatrooms() []Chatroom
	CreateRoom(Chatroom) []Chatroom
}

var Db *sql.DB
var Db2 ChatroomDatabase

type chatroomPostgresDb struct{}

func (*chatroomPostgresDb) Chatrooms() []Chatroom {
	rows, err := Db.Query("SELECT id, uuid, name FROM chatrooms ORDER BY id ASC")
	var rooms []Chatroom
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		room := Chatroom{}
		err = rows.Scan(&room.Id, &room.Uuid, &room.Name)
		if err != nil {
			panic(err)
		}
		rooms = append(rooms, room)
	}
	fmt.Println(rooms)
	rows.Close()
	return rooms
}
func (d *chatroomPostgresDb) CreateRoom(chatroom Chatroom) []Chatroom {
	statement := `insert into chatrooms ( uuid, name) values ($1, $2) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(createUUID(), chatroom.Name).Scan(&chatroom.Id)
	if err != nil {
		panic(err)
	}
	return d.Chatrooms()
}
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
	Db2 = &chatroomPostgresDb{}
	mux.HandleFunc("/chatrooms", handleChatroom)
	mux.HandleFunc("/chat", chat)
	mux.HandleFunc("/ws", handleWs)

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
