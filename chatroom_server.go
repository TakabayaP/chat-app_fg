package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Chatroom struct {
	Id   int
	Uuid string
	Name string
}

func (chatroom *Chatroom) Create() (err error) {
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
	return
}

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

func handleChatroomSelect(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
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
		rows.Close()
		output, err := json.MarshalIndent(&rooms, "", "\t\t")
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write(output)
	}
}
