package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"html/template"
	"database/sql"
	//"github.com/mattn/go-sqlite3"
	"log"
	//"os"
	//"fmt"
)

var DATABASE = "/tmp/minitwit.db"
var PER_PAGE = 30
var DEBUG = true
var SECRET_KEY = "development key"

var (
	db *sql.DB
	user *string
)

func connect_db() (*sql.DB){
	db_, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		log.Fatal(err)
	}
	return db_
}

func followUser(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//username := vars["username"]
	user := "jonas"
	
	if len(user) <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
	}

	whom_id := 1

	if whom_id <= 0 {
		w.WriteHeader(http.StatusNotFound)
	}
}

func before_req(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		db = connect_db()
		user = nil
		defer db.Close()
		handler(w, r)
	}
}

type Student struct {
	Name string
}

func public_timeline(w http.ResponseWriter, r *http.Request) {
	student := Student{
		Name: "Jo",
	}
	parsedTemp, _ := template.ParseFiles("test.html")
	err := parsedTemp.Execute(w, student)
	if err != nil {
		log.Println("Error executing template: ", err)
		return
	}
}


func main() {


	r := mux.NewRouter()

	r.HandleFunc("/{username}/follow", before_req(followUser)).Methods("GET")
	r.HandleFunc("/", public_timeline).Methods("GET")

	http.ListenAndServe(":8080", r) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}