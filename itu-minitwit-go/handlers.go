package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const DRIVER = "sqlite3"
const DATABASE = "../db_backup/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

//LoginHandler ...
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
			message := Message{Author_id: userId.(uint), Text: textValue, Pub_date: int(time.Now().Unix()), Flagged: 0}
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

func FollowUserHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "session_cookie")
		userId := session.Values["user_id"]
		if userId == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		params := mux.Vars(r)
		whomUsername := params["username"]

		whom, err := GetUserByUsername(whomUsername, db)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		follower := Follower{
			WhoID:  userId.(uint),
			WhomID: whom.UserID,
		}

		result := db.Create(&follower)
		log.Println(result)
		if result.Error != nil {
			log.Fatal(result.Error)
		}
		session.AddFlash(fmt.Sprintf("You are now following %s.", whomUsername))
		http.Redirect(w, r, fmt.Sprintf("/%s", whomUsername), http.StatusFound)
	})
}

func UnfollowUserHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "session_cookie")
		userId := session.Values["user_id"]
		if userId == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		params := mux.Vars(r)
		whomUsername := params["username"]

		whom, err := GetUserByUsername(whomUsername, db)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		follower := Follower{
			WhoID:  userId.(uint),
			WhomID: whom.UserID,
		}

		result := db.Where("who_id = ? and whom_id = ?", follower.WhoID, follower.WhomID).Delete(&follower)
		log.Println(result)
		if result.Error != nil {
			log.Fatal(result.Error)
		}
		session.AddFlash(fmt.Sprintf("You are no longer following %s.", whomUsername))
		http.Redirect(w, r, "/", http.StatusFound)
	})
}
