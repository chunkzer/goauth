package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/signup", signUpHandler).Methods("POST")
	r.HandleFunc("/welcome", welcomeHandler).Methods("GET")
	return r
}

var db *sql.DB
var jwtKey = []byte("secret_key")

func main() {
	initDB()
	r := newRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}

func initDB() {
	var err error

	db, err = sql.Open("postgres", "dbname=mydb sslmode=disable")
	if err != nil {
		panic(err)
	}
	log.Println("Database running successfully")
}
