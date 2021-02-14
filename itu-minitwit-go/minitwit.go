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
	err := db.QueryRow("SELECT * FROM user WHERE username= ?", username, 1).
		Scan(&user.User_id, &user.Username, &user.Email, &user.Pw_hash)

	return user, err
}

func GetUserById(id int, db *sql.DB) (User, error) {
	user := User{}
	err := db.QueryRow("SELECT * FROM user WHERE user_id=%d", id, 1).
		Scan(&user.User_id, &user.Username, &user.Email, &user.Pw_hash)

	return user, err
}

func LoginHandler(store *sessions.CookieStore, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session_cookie")

		userId := session.Values["user_id"]
		if isLoggedIn := userId != nil; isLoggedIn {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		var errorMsg string
		if r.Method == "POST" {
			user, err := GetUserByUsername(r.FormValue("username"), db)
			if err != nil {
				errorMsg = "Invalid username"
			} else if err = bcrypt.CompareHashAndPassword([]byte(user.Pw_hash), []byte(r.FormValue("password"))); err != nil {
				errorMsg = "Invalid password"
			} else {
				session.AddFlash("You were logged in")
				session.Values["user_id"] = user.User_id
				err = session.Save(r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, "/", http.StatusFound)
			}
		}

		response := map[string]string{"error": errorMsg}
		log.Println(response)
		// TODO render login template with error
	})
}

func LogoutHandler(store *sessions.CookieStore, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session_cookie")
		session.Values["user_id"] = nil
		session.AddFlash("You were logged out")

		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func RegisterHandler(store *sessions.CookieStore, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session_cookie")
		isLoggedIn := session.Values["user_id"] != nil
		if isLoggedIn {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		var errorMsg string
		if r.Method == "POST" {
			if len(r.FormValue("username")) == 0 {
				errorMsg = "You have to enter a username"
			} else if len(r.FormValue("email")) == 0 || !strings.Contains(r.FormValue("email"), "@") {
				errorMsg = "You have to enter a valid email address"
			} else if len(r.FormValue("password")) == 0 {
				errorMsg = "You have to enter a password"
			} else if r.FormValue("password") != r.FormValue("password2") {
				errorMsg = "The passwords do not match"
			} else if user, _ := GetUserByUsername(r.FormValue("username"), db); user.Username == r.FormValue("username") {
				errorMsg = "This username is already taken"
			} else {
				statement, err := db.Prepare("INSERT INTO user (username, email, pw_hash) values (?,?,?)")
				if err != nil {
					log.Println(err)
					return
				}
				defer statement.Close()

				pass := r.FormValue("password")
				hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
				if err != nil {
					log.Println(err)
					return
				}

				statement.Exec(r.FormValue("username"), r.FormValue("email"), hash)
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
		rows, err := db.Query("SELECT COUNT(*) FROM user")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var count int
			if err := rows.Scan(&count); err != nil {
				log.Fatal(err)
			}
			fmt.Println(count)
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
			session, _ := store.Get(r, "session_cookie")
			userId := session.Values["user_id"]
			if userId != nil {
				id := userId.(int)
				user, err := GetUserById(id, db)
				if err != nil {
					log.Print(err)
				}

				session.Values["user_id"] = user.User_id
				err = session.Save(r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
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
	r.Handle("/login", LoginHandler(store, db)).Methods("GET", "POST")
	r.Handle("/register", RegisterHandler(store, db)).Methods("GET", "POST")
	r.Handle("/logout", LogoutHandler(store, db)).Methods("GET")
	r.Handle("/add_message", TestHandler(db)).Methods("POST")
	r.Handle("{username}/", TestHandler(db)).Methods("GET")
	r.Handle("{username}/follow", TestHandler(db)).Methods("GET")
	r.Handle("{username}/unfollow", TestHandler(db)).Methods("GET")
	r.Handle("/test", TestHandler(db)).Methods("GET")
	r.Handle("/user/{id}", GetUserByIdHandler(db)).Methods("GET")

	http.ListenAndServe(":8080", r)
}
