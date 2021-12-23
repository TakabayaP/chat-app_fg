package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type Post struct {
	Id        int
	Body      string `json:"content"`
	UserId    int    `json:"user_id"`
	CreatedAt time.Time
}

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "dbname=chat sslmode=disable")
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("/public/"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", index)
	mux.HandleFunc("/post", post)
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "bao")
}

func read(w http.ResponseWriter, r *http.Request) {

}
func Posts(limit int) (posts []Post, err error) {
	rows, err := Db.Query("select id, body, user_id from posts limit $1", limit)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Body, &post.UserId)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

func GetPost(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRow("select id, body, user_id from posts where id=$1", id).Scan(&post.Id, &post.Body, &post.UserId)
	return
}

func (post *Post) Create() (err error) {
	statement := `insert into posts (body, user_id,created_at) values ($1, $2,$3) returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Body, post.UserId, time.Now()).Scan(&post.Id)
	return
}

func (post *Post) Update() (err error) {
	_, err = Db.Exec("update posts set body = $2, user_id = $3 where id = $1", post.Id, post.Body, post.UserId)
	return
}

func (post *Post) Delete() (err error) {
	_, err = Db.Exec("delete from posts where id = $1", post.Id)
	return
}

func DeleteAll() (err error) {
	_, err = Db.Exec("delete from posts")
	return
}

func post(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = handleGetPost(w, r)
	case "POST":
		err = handlePostPost(w, r)
	case "PUT":
		err = handlePutPost(w, r)
	case "DELETE":
		err = handleDeletePost(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGetPost(w http.ResponseWriter, r *http.Request) (err error) {
	rows, err := Db.Query("SELECT id, body, user_id, created_at  FROM posts ORDER BY created_at DESC")
	var posts []Post
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		conv := Post{}
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
	w.Write(output)
	return
}
func handlePostPost(w http.ResponseWriter, r *http.Request) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var post Post
	json.Unmarshal(body, &post)
	err = post.Create()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}
func handlePutPost(w http.ResponseWriter, r *http.Request) (err error)    { return }
func handleDeletePost(w http.ResponseWriter, r *http.Request) (err error) { return }
