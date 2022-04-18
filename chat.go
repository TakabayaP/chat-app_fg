package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Chat struct {
	Id        int
	Body      string `json:"content"`
	UserId    int    `json:"user_id"`
	CreatedAt time.Time
}

func (chat *Chat) Create() (err error) {
	statement := `insert into posts (body, user_id,created_at) values ($1, $2,$3) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(chat.Body, chat.UserId, time.Now()).Scan(&chat.Id)
	if err != nil {
		panic(err)
	}
	return
}

func chat(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = handleGetChat(w, r)
	case "POST":
		err = handlePostChat(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGetChat(w http.ResponseWriter, r *http.Request) (err error) {
	rows, err := Db.Query("SELECT id, body, user_id, created_at  FROM posts ORDER BY created_at ASC")
	var posts []Chat
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		conv := Chat{}
		if err = rows.Scan(&conv.Id, &conv.Body, &conv.UserId, &conv.CreatedAt); err != nil {
			panic(err)
		}
		posts = append(posts, conv)
	}
	rows.Close()
	output, err := json.MarshalIndent(&posts, "", "\t\t")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(output)
	return
}

func handlePostChat(w http.ResponseWriter, r *http.Request) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var chat Chat
	json.Unmarshal(body, &chat)
	err = chat.Create()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}
