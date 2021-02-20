package main

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

func LatestHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
	})
}

func ResgisterApiHandler(store *sessions.CookieStore, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO
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
