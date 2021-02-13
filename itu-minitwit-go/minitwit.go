package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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

func GetUserByUsername(username string, db *sql.DB) (User, error) {
	user := User{}
	err := db.QueryRow(fmt.Sprintf("SELECT * FROM user WHERE username=%s", username), 1).
		Scan(&user.User_id, &user.Username, &user.Email, &user.Pw_hash)

	return user, err
}

func GetUserById(id int, db *sql.DB) (User, error) {
	user := User{}
	err := db.QueryRow(fmt.Sprintf("SELECT * FROM user WHERE user_id=%d", id), 1).
		Scan(&user.User_id, &user.Username, &user.Email, &user.Pw_hash)

	return user, err
}

func RegisterHandler(store *sessions.CookieStore, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session_cookie")
		if err != nil {
			log.Println(err)
		}

		isLoggedIn := session.Values["user_id"] != nil
		if isLoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		var errorMsg string
		if r.Method == "POST" {
			if len(r.Form.Get("username")) == 0 {
				errorMsg = "You have to enter a username"
			} else if len(r.Form.Get("email")) == 0 || !strings.Contains(r.Form.Get("email"), "@") {
				errorMsg = "You have to enter a valid email address"
			} else if len(r.Form.Get("password")) == 0 {
				errorMsg = "You have to enter a password"
			} else if r.Form.Get("password") != r.Form.Get("password2") {
				errorMsg = "The passwords do not match"
			} else if user, _ := GetUserByUsername(r.Form.Get("username"), db); user.Username == r.Form.Get("username") {
				errorMsg = "This username is already taken"
			} else {
				statement, err := db.Prepare("INSERT INTO user (username, email, pw_hash) values (?,?,?)")
				if err != nil {
					log.Println(err)
					return
				}
				defer statement.Close()

				pass := r.Form.Get("password")
				hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
				if err != nil {
					log.Println(err)
					return
				}

				statement.Exec(r.Form.Get("username"), r.Form.Get("email"), hash)
				// TODO return successful registration status
			}
		}

		response := map[string]string{"error": errorMsg}
		log.Println(response)
		// w.Write([]byte(json.Marshal(response)))
		// TODO render register template with error
	})
}

func GetUserByIdHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := GetUserById(id, db)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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
			var user_id int
			var username, email, pw_hash string
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

func BeforeRequestMiddleware(store *sessions.CookieStore, db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mdfn := func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "session_cookie")
			if err != nil {
				log.Println(err)
			}

			userId := session.Values["user_id"]
			if userId != nil {
				id := userId.(int)
				user, err := GetUserById(id, db)
				if err != nil {
					log.Print(err)
				}

				session.Values["user_id"] = user.User_id
				session.Values["username"] = user.Username
				err = session.Save(r, w)
				if err != nil {
					log.Println(err)
				}
			}

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
	r.Use(BeforeRequestMiddleware(store, db))
	r.Handle("/", HomeHandler()).Methods("GET")
	r.Handle("/public", TestHandler(db)).Methods("GET")
	r.Handle("/login", TestHandler(db)).Methods("GET", "POST")
	r.Handle("/register", RegisterHandler(store, db)).Methods("GET", "POST")
	r.Handle("/logout", TestHandler(db)).Methods("GET")
	r.Handle("/add_message", TestHandler(db)).Methods("POST")
	r.Handle("{username}/", TestHandler(db)).Methods("GET")
	r.Handle("{username}/follow", TestHandler(db)).Methods("GET")
	r.Handle("{username}/unfollow", TestHandler(db)).Methods("GET")
	r.Handle("/test", TestHandler(db)).Methods("GET")
	r.Handle("/user/{id}", GetUserByIdHandler(db)).Methods("GET")

	http.ListenAndServe(":8080", r)
}
