package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
				Status:    http.StatusBadRequest,
				Error_msg: error_msg,
			}
			json.NewEncoder(w).Encode(&e)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Pwd), bcrypt.MinCost)

		user := User{
			Username: u.Username,
			PwHash:   string(hash),
			Email:    u.Email,
		}

		db.Create(&user)

		e := Response{
			Status:    http.StatusNoContent,
			Error_msg: "",
		}
		json.NewEncoder(w).Encode(&e)
	})
}

func MessagesHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
	})
}

func MessagesPerUserHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
	})
}

func FollowHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
	})
}

func AuthenticationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mdfn := func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
				response := Response{
					Status:    http.StatusForbidden,
					Error_msg: "You are not authorized to use this resource!",
				}
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(&response)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(mdfn)
	}
}

func LatestMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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
