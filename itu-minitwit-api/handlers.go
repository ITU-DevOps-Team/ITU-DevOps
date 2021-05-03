package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//Prometheus metrics
var (
	minitwit_api_register_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_api_register_requests",
		Help: "The count of requests to the /register endpoint of the backend API",
	})
	minitwit_api_messages_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_api_messages_requests",
		Help: "The count of requests to the /msgs endpoint of the backend API",
	})
	minitwit_api_messages_per_user_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_api_messages_per_user_requests",
		Help: "The count of requests to the /msgs/{username} endpoint of the backend API",
	})
	minitwit_api_follow_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_api_follow_requests",
		Help: "The count of requests to the /fllws/{username} endpoint of the backend API",
	})
	minitwit_api_total_requests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_api_total_requests",
		Help: "The total count of requests to the backend API",
	})
	minitwit_api_latest_execution_time_in_ns = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "minitwit_api_latest_execution_time_in_ns",
		Help:    "Histogram of the execution time of the latest middleware of the backend API ",
		Buckets: []float64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000},
	})
	minitwit_api_register_execution_time_in_ns = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "minitwit_api_register_execution_time_in_ns",
		Help:    "Histogram of the execution time of the RegisterApiHandler of the backend API ",
		Buckets: []float64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000},
	})
	minitwit_api_messages_execution_time_in_ns = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "minitwit_api_messages_execution_time_in_ns",
		Help:    "Histogram of the execution time of the MessagesHandler of the backend API ",
		Buckets: []float64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000},
	})
	minitwit_api_messages_per_user_execution_time_in_ns = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "minitwit_api_messages_per_user_execution_time_in_ns",
		Help:    "Histogram of the execution time of the MessagesPerUserHandler of the backend API ",
		Buckets: []float64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000},
	})
	minitwit_api_follow_execution_time_in_ns = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "minitwit_api_follow_execution_time_in_ns",
		Help:    "Histogram of the execution time of the FollowHandler of the backend API ",
		Buckets: []float64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000},
	})
	minitwit_api_authentication_middleware_execution_time_in_ns = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "minitwit_api_authentication_middleware_execution_time_in_ns",
		Help:    "Histogram of the execution time of the AuthenticationMiddleware of the backend API ",
		Buckets: []float64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000},
	})
	minitwit_api_latest_middleware_execution_time_in_ns = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "minitwit_api_latest_middleware_execution_time_in_ns",
		Help:    "Histogram of the execution time of the LatestMiddleware of the backend API ",
		Buckets: []float64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000},
	})
)

func LatestHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var t0 = time.Now()

		latestObj, _ := GetLatest(db)
		json.NewEncoder(w).Encode(&latestObj)

		var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
		minitwit_api_latest_execution_time_in_ns.Observe(elapsed)

	})
}

func RegisterApiHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var t0 = time.Now()
		minitwit_api_register_requests.Inc()

		var u User_
		errormsg := ""
		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
			minitwit_api_register_execution_time_in_ns.Observe(elapsed)
			return
		}

		if u.Username == "" {
			errormsg = "You have to enter a username"
		} else if u.Email == "" || !strings.ContainsAny(u.Email, "@") {
			errormsg = "You have to enter a email"
		} else if u.Pwd == "" {
			errormsg = "You have to enter a password"
		} else if user_check, _ := GetUserByUsername(u.Username, db); user_check.Username == u.Username {
			errormsg = "The username is already taken"
		}

		if errormsg != "" {
			e := Response{
				Status:    http.StatusBadRequest,
				Error_msg: errormsg,
			}

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&e)
			var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
			minitwit_api_register_execution_time_in_ns.Observe(elapsed)
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

		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(&e)
		var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
		minitwit_api_register_execution_time_in_ns.Observe(elapsed)
	})
}

func MessagesHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var t0 = time.Now()
		minitwit_api_messages_requests.Inc()

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
				PubDate int    `json:"pub_date"`
				User    string `json:"user"`
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
			var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
			minitwit_api_messages_execution_time_in_ns.Observe(elapsed)
			return
		}
	})
}

func UserHandlerGet(db *gorm.DB, w *http.ResponseWriter, r *http.Request, t0 time.Time, numberOfMessages int, ) {
	username := mux.Vars(r)["username"]
	user, err := GetUserByUsername(username, db)
	if err != nil {
		http.Error(*w, "User not found", http.StatusNotFound)
		var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
		minitwit_api_messages_per_user_execution_time_in_ns.Observe(elapsed)
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
		PubDate int    `json:"pub_date"`
		User    string `json:"user"`
	}

	var filteredMessages []filteredMessage

	for i := range queryData {
		filteredMsg := filteredMessage{}
		filteredMsg.Content = queryData[i].Text
		filteredMsg.PubDate = queryData[i].PubDate
		filteredMsg.User = queryData[i].Username
		filteredMessages = append(filteredMessages, filteredMsg)
	}

	json.NewEncoder(*w).Encode(&filteredMessages)
	var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
	minitwit_api_messages_per_user_execution_time_in_ns.Observe(elapsed)
	return
}

func UserHandlerPost(db *gorm.DB, w *http.ResponseWriter, r *http.Request, t0 time.Time) {
	var msg Message_
	err := json.NewDecoder(r.Body).Decode(&msg)

	if err != nil {
		http.Error(*w, err.Error(), http.StatusBadRequest)
		var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
		minitwit_api_messages_per_user_execution_time_in_ns.Observe(elapsed)
		return
	}

	username := mux.Vars(r)["username"]
	user, err := GetUserByUsername(username, db)
	if err != nil {
		http.Error(*w, "User not found", http.StatusNotFound)
		var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
		minitwit_api_messages_per_user_execution_time_in_ns.Observe(elapsed)
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
		log.Println(result.Error)
	}

	(*w).WriteHeader(http.StatusNoContent)
}

func MessagesPerUserHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var t0 = time.Now()

		minitwit_api_messages_per_user_requests.Inc()

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
			UserHandlerGet(db, &w, r, t0, numberOfMessages)
		} else if r.Method == "POST" {
			UserHandlerPost(db, &w, r, t0)
		}
	})
}

func FollowHandler(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var t0 = time.Now()

		minitwit_api_follow_requests.Inc()

		username := mux.Vars(r)["username"]
		user, err := GetUserByUsername(username, db)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
			minitwit_api_follow_execution_time_in_ns.Observe(elapsed)
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
				var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
				minitwit_api_follow_execution_time_in_ns.Observe(elapsed)
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
			var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
			minitwit_api_follow_execution_time_in_ns.Observe(elapsed)
			return
		} else if r.Method == "POST" && requestBody.Unfollow != "" {
			userToUnfollow, err := GetUserByUsername(requestBody.Unfollow, db)
			if err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
				minitwit_api_follow_execution_time_in_ns.Observe(elapsed)
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
			var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
			minitwit_api_follow_execution_time_in_ns.Observe(elapsed)
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
			var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
			minitwit_api_follow_execution_time_in_ns.Observe(elapsed)
			return
		}
	})
}

func AuthenticationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		var t0 = time.Now()
		minitwit_api_total_requests.Inc()
		mdfn := func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
				response := Response{
					Status:    http.StatusForbidden,
					Error_msg: "You are not authorized to use this resource!",
				}
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(&response)
				var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
				minitwit_api_authentication_middleware_execution_time_in_ns.Observe(elapsed)
				return
			}

			next.ServeHTTP(w, r)
		}
		var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
		minitwit_api_authentication_middleware_execution_time_in_ns.Observe(elapsed)
		return http.HandlerFunc(mdfn)
	}
}

func LatestMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		var t0 = time.Now()

		mdfn := func(w http.ResponseWriter, r *http.Request) {

			//If the request comes from Prometheus, skip the latestmiddleware
			agent := r.Header.Get("User-Agent")
			if strings.Split(agent, "/")[0] != "Prometheus" {

				keys, ok := r.URL.Query()["latest"]

				if !ok || len(keys[0]) < 1 {
					log.Println("Request does not contain a new latest value")
				} else {
					latest, _ := strconv.Atoi(keys[0])
					latestObj := Latest{latest}
					AddLatest(latestObj, db)
					log.Println(latest)
				}

			}

			next.ServeHTTP(w, r)
		}
		var elapsed = float64((time.Now().Sub(t0)).Nanoseconds())
		minitwit_api_latest_middleware_execution_time_in_ns.Observe(elapsed)
		return http.HandlerFunc(mdfn)
	}
}
