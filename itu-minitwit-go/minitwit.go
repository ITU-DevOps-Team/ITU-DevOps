package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

func initDb(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	sql, _ := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	return db, sql.Ping()
}

func ReadDVariables() (string, error) {
	var err error

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		err = errors.New("env var missing (DB_NAME)")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		err = errors.New("env var missing (DB_USER)")
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		err = errors.New("env var missing (DB_PASS)")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		err = errors.New("env var missing (DB_HOST)")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		err = errors.New("env var missing (DB_PORT)")
	}

return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort), err
}


//Initialize prometheus
func init() {
	prometheus.MustRegister(minitwit_ui_total_requests)
	prometheus.MustRegister(minitwit_ui_login_requests)
	prometheus.MustRegister(minitwit_ui_logout_requests)
	prometheus.MustRegister(minitwit_ui_register_requests)
	prometheus.MustRegister(minitwit_ui_homepage_requests)
	prometheus.MustRegister(minitwit_ui_addmessage_requests)
	prometheus.MustRegister(minitwit_ui_follow_requests)
	prometheus.MustRegister(minitwit_ui_unfollow_requests)
	prometheus.MustRegister(minitwit_ui_personaltimeline_requests)
	prometheus.MustRegister(minitwit_ui_usertimeline_requests)
}



func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Failed to read .env file")
		log.Println("Getting environment variables from env instead...")
	}

	dsn, err := ReadDVariables()
	if err != nil {
		log.Fatal(err)
	}

	gorm, err := initDb(dsn)
	if err != nil {
		log.Fatal(err)
	}

	LoadTemplates()
	store := sessions.NewCookieStore([]byte(SECRET_KEY))

	r := mux.NewRouter()
	r.Use(BeforeRequestMiddleware(store, gorm))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))

	r.Handle("/", HomeHandler(store, gorm)).Methods("GET")
	r.Handle("/login", LoginHandler(store, gorm)).Methods("GET", "POST")
	r.Handle("/register", RegisterHandler(store, gorm)).Methods("GET", "POST")
	r.Handle("/logout", LogoutHandler(store, gorm)).Methods("GET")
	r.Handle("/add_message", AddMessageHandler(store, gorm)).Methods("POST")
	r.Handle("/personaltimeline", PersonalTimeline(store, gorm)).Methods("GET", "POST")
	r.Handle("/metrics",promhttp.Handler())
	r.Handle("/{username}", UserTimeline(store, gorm)).Methods("GET")
	r.Handle("/{username}/follow", FollowUserHandler(store, gorm)).Methods("GET")
	r.Handle("/{username}/unfollow", UnfollowUserHandler(store, gorm)).Methods("GET")

	http.ListenAndServe(":8080", r)
}
