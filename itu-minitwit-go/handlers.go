package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const DRIVER = "sqlite3"
const DATABASE = "../db_backup/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

// Prometheus metrics
var (
	minitwit_ui_usertimeline_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_ui_usertimeline_requests",
		Help: "The count of HTTP requests to the /{username} endpoint of the frontend API.",
	})
	minitwit_ui_personaltimeline_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_ui_personaltimeline_requests",
		Help: "The count of HTTP requests to the /personaltimeline endpoint of the frontend API.",
	})
	minitwit_ui_unfollow_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_ui_unfollow_requests",
		Help: "The count of HTTP requests to the /{username}/unfollow endpoint of the frontend API.",
	})
	minitwit_ui_follow_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_ui_follow_requests",
		Help: "The count of HTTP requests to the /{username}/follow endpoint of the frontend API.",
	})
	minitwit_ui_addmessage_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_ui_addmessage_requests",
		Help: "The count of HTTP requests to the /add_message endpoint of the frontend API.",
	})
	minitwit_ui_homepage_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_ui_homepage_requests",
		Help: "The count of HTTP requests to the / (home) endpoint of the frontend API.",
	})
	minitwit_ui_register_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_ui_register_requests",
		Help: "The count of HTTP requests to the /register endpoint of the frontend API.",
	})
	minitwit_ui_login_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_ui_login_requests",
		Help: "The count of HTTP requests to the /login endpoint of the frontend API.",
	})
	minitwit_ui_logout_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_ui_logout_requests",
		Help: "The count of HTTP requests to the /logout endpoint of the frontend API.",
	})
	minitwit_ui_total_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_ui_total_requests",
		Help: "The count of HTTP requests to the frontend API.",
	})
)

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
		minitwit_ui_login_requests.Inc()
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
				log.WithFields(log.Fields{}).Error("Error occured when executing the login template.")
			}
		} else if r.Method == "POST" {
			user, err := GetUserByUsername(r.FormValue("username"), db)
			if err != nil {
				errorMsg = "Invalid username"
				log.WithFields(log.Fields{"user": r.FormValue("username")}).Error("User entered invalid username.")
			} else if err = bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(r.FormValue("password"))); err != nil {
				errorMsg = "Invalid password"
				log.WithFields(log.Fields{"user": r.FormValue("username")}).Error("User entered invalid password.")
			} else {
				session.AddFlash("You were logged in")
				session.Values["user_id"] = user.UserID
				log.WithFields(log.Fields{"user": r.FormValue("username"), "userId": user.UserID}).Info("User successfully logged in.")
				err = session.Save(r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					log.WithFields(log.Fields{"method": "loginHandler.go"}).Error("Error occured when saving the session")
					return
				}
				http.Redirect(w, r, "/", http.StatusFound)
			}
			//renders sign in page again with error
			viewContent := ViewContent{
				Error:        true,
				ErrorMessage: errorMsg,
			}

			if err := templates["login"].Execute(w, viewContent); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.WithFields(log.Fields{}).Error("Error occured when executing the login template.")
			}
		}
	})
}

func LogoutHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		minitwit_ui_logout_requests.Inc()
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
		minitwit_ui_register_requests.Inc()
		session, _ := store.Get(r, "session_cookie")
		userId := session.Values["user_id"]
		isLoggedIn := userId != "" && userId != nil
		if isLoggedIn {
			fmt.Println("user already signed in")
			http.Redirect(w, r, "/", http.StatusFound)
		}

		if r.Method == "GET" {
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
			user, _ := GetUserByUsername(formUsername, db)
			usernameTaken := user.Username == formUsername

			if isEmpty_formUsername {
				errorMsg = "username is empty"
			} else if isEmpty_formEmail {
				errorMsg = "email is empty"
			} else if isEmpty_formPassword {
				errorMsg = "password is empty"
			} else if isEmpty_formPasswordConfirm {
				errorMsg = "password repeat is empty"
			} else if incorrectFormat_formEmail {
				errorMsg = "You have to enter a valid email address"
			} else if formPassword != formPasswordConfirm {
				//Passwords does not match
				errorMsg = "password and repeated password does not match"
			} else if usernameTaken {
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
					Success:        true,
					SuccessMessage: "User successfully created",
				}

				if err := templates["login"].Execute(w, viewContent); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}

			if errorMsg != "" {
				//renders register page again with error
				viewContent := ViewContent{
					Error:        true,
					ErrorMessage: errorMsg,
				}

				if err := templates["register"].Execute(w, viewContent); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		}
	})
}

func HomeHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		minitwit_ui_homepage_requests.Inc()
		session, _ := store.Get(r, "session_cookie")

		userId := session.Values["user_id"]
		isLoggedIn := userId != "" && userId != nil

		posts := GetPublicPosts(10, db)

		viewContent := ViewContent{
			SignedIn: isLoggedIn,
			Posts:    posts,
		}
		if err := templates["public_timeline"].Execute(w, viewContent); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func AddMessageHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		minitwit_ui_addmessage_requests.Inc()

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

			//Increment number of http requests in Prometheus
			minitwit_ui_total_requests.Inc()

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

		minitwit_ui_follow_requests.Inc()

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

		minitwit_ui_unfollow_requests.Inc()

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
		http.Redirect(w, r, fmt.Sprintf("/%s", whomUsername), http.StatusFound)
	})
}

func PersonalTimeline(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		minitwit_ui_personaltimeline_requests.Inc()

		session, _ := store.Get(r, "session_cookie")
		userId := session.Values["user_id"]
		isLoggedIn := userId != "" && userId != nil
		currentUserName := "" //username of the current user signed in

		//redirect to public if not logged in
		if !isLoggedIn {
			http.Redirect(w, r, "/", 302)
		}

		user, err := GetUserById(userId.(uint), db)
		checkErr(err)
		currentUserName = user.Username

		if r.Method == "GET" {
			viewContent := ViewContent{
				SignedIn: isLoggedIn,
				SameUser: true, //must be
				Posts:    GetPostsByUser(currentUserName, db),
				Username: currentUserName,
			}

			if err := templates["personal_timeline"].Execute(w, viewContent); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if r.Method == "POST" {
			r.ParseForm()
			var errorMsg string
			formText := r.FormValue("text")
			isEmpty_formText := formText == ""

			if isEmpty_formText {
				errorMsg = "post is empty"
			} else {

				if userId == nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if formText != "" {
					message := Message{
						Author_id: userId.(uint),
						Text:      formText,
						Pub_date:  int(time.Now().Unix()),
						Flagged:   0,
					}
					result := db.Create(&message)

					if result.Error != nil {
						log.Fatal(result.Error)
					}
				}

				viewContent := ViewContent{
					SignedIn:       isLoggedIn,
					SameUser:       true, //must be
					Posts:          GetPostsByUser(currentUserName, db),
					Username:       currentUserName,
					Success:        true,
					SuccessMessage: "Post successfully created",
				}

				if err := templates["personal_timeline"].Execute(w, viewContent); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}

			if errorMsg != "" {
				//display error
				viewContent := ViewContent{
					Error:        true,
					ErrorMessage: errorMsg,
				}

				if err := templates["personal_timeline"].Execute(w, viewContent); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}

		}
	})
}

func UserTimeline(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		minitwit_ui_usertimeline_requests.Inc()

		session, _ := store.Get(r, "session_cookie")
		userId := session.Values["user_id"]
		isLoggedIn := userId != "" && userId != nil
		currentUserName := "" //username of the current user signed in

		if isLoggedIn {
			user, err := GetUserById(userId.(uint), db)
			checkErr(err)
			currentUserName = user.Username
		}

		vars := mux.Vars(r)
		usernameVisited := vars["username"]

		//check if username visited exists
		if !CheckUsernameExists(usernameVisited, db) {
			http.Redirect(w, r, "/", 302)
		}

		//if visited user is same as logged in user redirect to personal timeline
		if currentUserName == usernameVisited {
			PersonalTimeline(store, db)
		}

		viewContent := ViewContent{
			SignedIn:         isLoggedIn,
			SameUser:         currentUserName == usernameVisited,
			Posts:            GetPostsByUser(usernameVisited, db),
			Username:         usernameVisited,
			AlreadyFollowing: CheckIfUserIsFollowed(currentUserName, usernameVisited, db),
		}

		if err := templates["personal_timeline"].Execute(w, viewContent); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
