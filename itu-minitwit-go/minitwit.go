package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const DRIVER = "sqlite3"
const DATABASE = "/tmp/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

type User struct {
	User_id  int64
	Username string
	Email    string
	Pw_hash  string
}

func GetUserByIdHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM user LIMIT 5")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		r.Context()
		fmt.Println(rows.Columns())
	})
}

func TestHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM user LIMIT 5")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		fmt.Println(rows.Columns())
	})
}

func HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to itu-minitwit"))
	})
}

func initDb(driver string, datasource string) (*sql.DB, error) {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		log.Fatal(err)
	}

	return db, db.Ping()
}

func main() {
	db, err := initDb(DRIVER, DATABASE)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	// r.Use(dbConnMiddleware)
	r.Handle("/", HomeHandler()).Methods("GET")
	r.Handle("/", TestHandler(db)).Methods("GET")
	r.Handle("/user/{id}", GetUserByIdHandler(db)).Methods("GET")

	http.ListenAndServe(":8080", r)
}
