package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
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
		//TODO: Update latest value?? Or does the middleware handle that?
		numberOfMessagesHeaderResult := r.URL.Query().Get("no")
		//default set to 100
		numberOfMessages := 100
		if numberOfMessagesHeaderResult != "" {
			parsed, err := strconv.Atoi(numberOfMessagesHeaderResult)
			if err == nil {
				numberOfMessages = parsed
			}
		}
		if r.Method == "GET" {
			//type used for scanning the query results
			type result struct {
				UserID    int
				Username  string
				Email     string
				PwHash    string
				MessageID int
				AuthorID  int
				Text      string
				PubDate   int
				Flagged   int
			}

			var queryData []result
			db.Table("messages").
				Joins("JOIN users on users.user_id = messages.author_id").
				Where("messages.flagged = 0").
				Order("messages.pub_date DESC").
				Limit(numberOfMessages).
				Scan(&queryData)

			type filteredMessage struct {
				Content string `json:"content"`
				PubDate int `json:"pub_date"`
				User string `json:"user"`
			}

			filteredMessages := []filteredMessage{}

			for i := range queryData {
				filteredMsg := filteredMessage{}
				filteredMsg.Content = queryData[i].Text
				filteredMsg.PubDate = queryData[i].PubDate
				filteredMsg.User = queryData[i].Username
				filteredMessages = append(filteredMessages, filteredMsg)
			}

			json.NewEncoder(w).Encode(&filteredMessages)
			return
		}
	})
}

func MessagesPerUserHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Update latest value?? Or does the middleware handle that?
		numberOfMessagesHeaderResult := r.URL.Query().Get("no")
		//default set to 100
		numberOfMessages := 100
		if numberOfMessagesHeaderResult != "" {
			parsed, err := strconv.Atoi(numberOfMessagesHeaderResult)
			if err == nil {
				numberOfMessages = parsed
			}
		}

		if r.Method == "GET" {
			username := mux.Vars(r)["username"]
			user, err := GetUserByUsername(username, db)
			if err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			//type used for scanning the query results
			type result struct {
				UserID    int
				Username  string
				Email     string
				PwHash    string
				MessageID int
				AuthorID  int
				Text      string
				PubDate   int
				Flagged   int
			}

			var queryData []result
			db.Table("messages").
				Joins("JOIN users on users.user_id = messages.author_id").
				Where("messages.flagged = 0 AND users.user_id = ?", user.UserID).
				Order("messages.pub_date DESC").
				Limit(numberOfMessages).
				Scan(&queryData)

			type filteredMessage struct {
				Content string `json:"content"`
				PubDate int `json:"pub_date"`
				User string `json:"user"`
			}

			filteredMessages := []filteredMessage{}

			for i := range queryData {
				filteredMsg := filteredMessage{}
				filteredMsg.Content = queryData[i].Text
				filteredMsg.PubDate = queryData[i].PubDate
				filteredMsg.User = queryData[i].Username
				filteredMessages = append(filteredMessages, filteredMsg)
			}

			json.NewEncoder(w).Encode(&filteredMessages)
			return
		} else if r.Method == "POST" {
			var msg Message_
			err := json.NewDecoder(r.Body).Decode(&msg)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			username := mux.Vars(r)["username"]
			user, err := GetUserByUsername(username, db)
			if err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			message := Message{
				Author_id: user.UserID,
				Text:      msg.Content,
				Pub_date:  int(time.Now().Unix()),
				Flagged:   0,
			}

			result := db.Create(&message)
			if result.Error != nil {
				log.Fatal(result.Error)
			}

			w.WriteHeader(http.StatusNoContent)
		}
	})
}

func FollowHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := mux.Vars(r)["username"]
		user, err := GetUserByUsername(username, db)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		type body struct {
			Follow   string `json:"follow"`
			Unfollow string `json:"unfollow"`
		}
		var requestBody body
		json.NewDecoder(r.Body).Decode(&requestBody)
		if r.Method == "POST" && requestBody.Follow != "" {
			userToFollow, err := GetUserByUsername(requestBody.Follow, db)
			if err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			follower := Follower{
				WhoID:  user.UserID,
				WhomID: userToFollow.UserID,
			}
			result := db.Create(&follower)
			if result.Error != nil {
				log.Fatal("Something went wrong when following")
			}

			w.WriteHeader(http.StatusNoContent)
			return
		} else if r.Method == "POST" && requestBody.Unfollow != "" {
			userToUnfollow, err := GetUserByUsername(requestBody.Unfollow, db)
			if err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			follower := Follower{
				WhoID:  user.UserID,
				WhomID: userToUnfollow.UserID,
			}
			result := db.Where("who_id = ? and whom_id = ?", follower.WhoID, follower.WhomID).Delete(&follower)
			if result.Error != nil {
				log.Fatal("Something when wrong when unfollowing")
			}

			w.WriteHeader(http.StatusNoContent)
			return
		} else if r.Method == "GET" {
			numberOfFollowersHeaderResult := r.URL.Query().Get("no")
			//default set to 100
			numberOfFollowers := 100
			if numberOfFollowersHeaderResult != "" {
				parsed, err := strconv.Atoi(numberOfFollowersHeaderResult)
				if err == nil {
					numberOfFollowers = parsed
				}
			}

			type result struct {
				Username string
			}

			followers := []string{}

			db.Table("users").
				Select("users.username").
				Joins("JOIN followers ON followers.whom_id = users.user_id").
				Where("followers.who_id = ?", user.UserID).
				Limit(numberOfFollowers).
				Scan(&followers)

			type response struct {
				Follows []string `json:"follows"`
			}
			followersResponse := response{Follows: followers}
			json.NewEncoder(w).Encode(&followersResponse)
			return
		}
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
