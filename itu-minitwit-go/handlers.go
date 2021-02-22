package main

import (
	"fmt"
	"html/template"
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

var templates map[string]*template.Template

func LoadTemplates() {
	var layoutTemplate = "templates/layout.gohtml"
	templates = make(map[string]*template.Template)

	templates["login"] = template.Must(template.ParseFiles(layoutTemplate, "templates/login.gohtml"))
	templates["register"] = template.Must(template.ParseFiles(layoutTemplate, "templates/register.gohtml"))
	templates["personal_timeline"] = template.Must(template.ParseFiles(layoutTemplate, "templates/personal_timeline.gohtml"))
	templates["public_timeline"] = template.Must(template.ParseFiles(layoutTemplate, "templates/public_timeline.gohtml"))
}

//LoginHandler ...
func LoginHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session_cookie")

		userId := session.Values["user_id"]
		isLoggedIn := userId != "" && userId != nil
		if isLoggedIn {
			fmt.Println("user already signed in -> redirecting to /")
			http.Redirect(w, r, "/", http.StatusFound)
		}

		var errorMsg string

		if r.Method == "GET" {
			if err := templates["login"].Execute(w, nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if r.Method == "POST" {
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
			//renders sign in page again with error
			viewContent := ViewContent{
				Error: true,
				ErrorMessage: errorMsg,
			}

			if err := templates["login"].Execute(w, viewContent); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			}
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
		userId := session.Values["user_id"]
		isLoggedIn := userId != "" && userId != nil
		if isLoggedIn {
			fmt.Println("user already signed in")
			http.Redirect(w, r, "/", http.StatusFound)
		}

		if r.Method == "GET"{
			if err := templates["register"].Execute(w, nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if r.Method == "POST" {

			//parsing form posted by user
			r.ParseForm()
			var errorMsg string

			formUsername := r.FormValue("username")
			formEmail := r.FormValue("email")
			formPassword := r.FormValue("password")
			formPasswordConfirm := r.FormValue("password_confirm")

			isEmpty_formUsername := formUsername == ""
			isEmpty_formEmail := formEmail == ""
			isEmpty_formPassword := formPassword == ""
			isEmpty_formPasswordConfirm := formPasswordConfirm == ""
			incorrectFormat_formEmail := !strings.Contains(formEmail, "@") || !strings.Contains(formEmail, ".")
			user, _ := GetUserByUsername(formUsername, db);
			usernameTaken := user.Username == formUsername

			if (isEmpty_formUsername){
				errorMsg = "username is empty"
			} else if (isEmpty_formEmail){
				errorMsg = "email is empty"
			} else if (isEmpty_formPassword){
				errorMsg = "password is empty"
			} else if (isEmpty_formPasswordConfirm){
				errorMsg = "password repeat is empty"
			} else if incorrectFormat_formEmail {
				errorMsg = "You have to enter a valid email address"
			} else if (formPassword != formPasswordConfirm){
				//Passwords does not match
				errorMsg = "password and repeated password does not match"
			} else if usernameTaken{
				//Username is already taken
				errorMsg = "username already exist"

			} else {
				//Sign up user
				hash, err := bcrypt.GenerateFromPassword([]byte(formPassword), bcrypt.MinCost)

				if err != nil {
					log.Println(err)
					return
				}

				user := User{
					Username: formUsername,
					Email:    formEmail,
					PwHash:   string(hash),
				}

				db.Create(&user)
				//renders sign in page again with error
				viewContent := ViewContent{
					Success: true,
					SuccessMessage: "User successfully created",
				}

				if err := templates["login"].Execute(w, viewContent); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}

			if errorMsg != "" {
				//renders register page again with error
				viewContent := ViewContent{
					Error: true,
					ErrorMessage: errorMsg,
				}

				if err := templates["register"].Execute(w, viewContent); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
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
