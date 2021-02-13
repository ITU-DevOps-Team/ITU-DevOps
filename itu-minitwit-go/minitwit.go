package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

const DRIVER = "sqlite3"
const DATABASE = "/tmp/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

type User struct {
	User_id  int
	Username string
	Email    string
	Pw_hash  string
}

func GetUserByIdHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := User{}
		db.QueryRow(fmt.Sprintf("SELECT * FROM user WHERE user_id=%d", id), 1).
			Scan(&user.User_id, &user.Username, &user.Email, &user.Pw_hash)

		file := path.Join(".", "templates", "greet.html")
		tmpl, err := template.ParseFiles(file)
		if err != nil {
			log.Fatal(err)
		}

		tmpl.Execute(w, user)
	})
}

func TestHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM user LIMIT 5")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var user_id, username, email, pw_hash string
			if err := rows.Scan(&user_id, &username, &email, &pw_hash); err != nil {
				log.Fatal(err)
			}
			fmt.Println(user_id, username, email, pw_hash)
		}
	})
}

func HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to itu-minitwit"))
	})
}

func BeforeRequestMiddleware(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mdfn := func(w http.ResponseWriter, r *http.Request) {
			// TODO: middleware logic here
			// check for user session in store

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(mdfn)
	}
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

	store := sessions.NewCookieStore([]byte(SECRET_KEY))
	r := mux.NewRouter()

	r.Use(BeforeRequestMiddleware(store))
	r.Handle("/", HomeHandler()).Methods("GET")
	r.Handle("/test", TestHandler(db)).Methods("GET")
	r.Handle("/user/{id}", GetUserByIdHandler(db)).Methods("GET")

	http.ListenAndServe(":8080", r)
}
