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
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DRIVER = "sqlite3"
const DATABASE = "../db_backup/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

func GetUserByUsername(username string, db *gorm.DB) (User, error) {
	user := User{}
	result := db.Where("username = ?", username).First(&user)
	return user, result.Error
}

func GetUserById(id uint, db *gorm.DB) (User, error) {
	user := User{}
	result := db.First(&user, id)
	return user, result.Error
}

func LoginHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session_cookie")

		userId := session.Values["user_id"]
		if isLoggedIn := userId != ""; isLoggedIn {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		var errorMsg string
		if r.Method == "POST" {
			user, err := GetUserByUsername(r.FormValue("username"), db)
			if err != nil {
				errorMsg = "Invalid username"
				log.Println(err)
			} else if err = bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(r.FormValue("password"))); err != nil {
				errorMsg = "Invalid password"
			} else {
				session.AddFlash("You were logged in")
				session.Values["user_id"] = user.UserID
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

func LogoutHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
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

func RegisterHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
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
				pass := r.FormValue("password")
				hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)

				if err != nil {
					log.Println(err)
					return
				}

				user := User{
					Username: r.FormValue("username"),
					Email:    r.FormValue("email"),
					PwHash:   string(hash),
				}

				db.Create(&user)
				http.Redirect(w, r, "/login", http.StatusCreated)
				// TODO return successful registration status
			}
		}

		response := map[string]string{"error": errorMsg}
		log.Println(response)
		// TODO render register template with error
	})
}

func GetUserByIdHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := GetUserById((uint(id)), db)
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

func AddMessageHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "session_cookie")
		userId := session.Values["user_id"]
		if userId == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		textValue := r.FormValue("text")
		if textValue != "" {
			message := Message{Author_id: userId.(uint), Text: textValue, Pub_date: strconv.Itoa(int(time.Now().Unix())), Flagged: 0}
			result := db.Create(&message)

			if result.Error != nil {
				log.Fatal(result.Error)
			}
			session.AddFlash("Your message was recorded")
			http.Redirect(w, r, "/", http.StatusFound)
		}
	})
}

func BeforeRequestMiddleware(store *sessions.CookieStore, db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mdfn := func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session_cookie")
			userId := session.Values["user_id"]
			if userId != nil {
				id := userId.(uint)
				user, err := GetUserById(id, db)
				if err != nil {
					log.Print(err)
				}

				session.Values["user_id"] = user.UserID
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

func initDb(driver string, datasource string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(datasource), &gorm.Config{})
	sql, _ := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	return db, sql.Ping()
}

func main() {
	gorm, err := initDb(DRIVER, DATABASE)
	// sql, _ := gorm.DB()
	if err != nil {
		log.Fatal(err)
	}

	store := sessions.NewCookieStore([]byte(SECRET_KEY))

	r := mux.NewRouter()
	r.Use(BeforeRequestMiddleware(store, gorm))
	r.Handle("/", HomeHandler()).Methods("GET")
	// r.Handle("/public", TestHandler(sql)).Methods("GET")
	r.Handle("/login", LoginHandler(store, gorm)).Methods("GET", "POST")
	r.Handle("/register", RegisterHandler(store, gorm)).Methods("GET", "POST")
	r.Handle("/logout", LogoutHandler(store, gorm)).Methods("GET")
	r.Handle("/add_message", AddMessageHandler(store, gorm)).Methods("POST")
	// r.Handle("/{username}", TestHandler(sql)).Methods("GET")
	// r.Handle("/{username}/follow", TestHandler(sql)).Methods("GET")
	// r.Handle("/{username}/unfollow", TestHandler(sql)).Methods("GET")
	// r.Handle("/test", TestHandler(sql)).Methods("GET")
	r.Handle("/user/{id}", GetUserByIdHandler(gorm)).Methods("GET")

	http.ListenAndServe(":8080", r)
}
