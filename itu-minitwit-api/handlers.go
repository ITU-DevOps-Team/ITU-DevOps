package main

import (
	"log"
	"strings"
	"net/http"
	"encoding/json"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
)



func LatestHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		latestObj := Latest{
			Latest: Latest_.Latest,
		}

		json.NewEncoder(w).Encode(&latestObj)

	})
}

func RegisterApiHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var u User_
		error_msg := ""
		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if u.Username == "" {
			error_msg = "You have to enter a username"
		} else if u.Email == "" || !strings.ContainsAny(u.Email, "@") {
			error_msg = "You have to enter a email"
		} else if u.Pwd == "" {
			error_msg = "You have to enter a password"
		} else if user_check, _ := GetUserByUsername(u.Username, db); user_check.Username == u.Username {
			error_msg = "The username is already taken"
		}


		if error_msg != "" {
			e := Response{
				Status: http.StatusBadRequest,
				Error_msg: error_msg,
			}
			json.NewEncoder(w).Encode(&e)
			return
		}
		
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Pwd), bcrypt.MinCost)
		
		user := User{
			Username: u.Username,
			PwHash: string(hash),
			Email: u.Email,
		}

		db.Create(&user)


		e := Response{
			Status: http.StatusNoContent,
			Error_msg: "",
		}
		json.NewEncoder(w).Encode(&e)
	})
}

func MessagesHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
	})
}

func MessagesPerUserHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
	})
}

func FollowHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
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
func LatestMiddleware() func(http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		mdfn := func(w http.ResponseWriter, r *http.Request) {
			keys, ok := r.URL.Query()["latest"]
		
			if !ok || len(keys[0]) < 1 {
			log.Println("Request does not contain a new latest value")
			} else {
				Latest_.Latest, _ = strconv.Atoi(keys[0])
				log.Println(Latest_.Latest)
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(mdfn)
	}
}