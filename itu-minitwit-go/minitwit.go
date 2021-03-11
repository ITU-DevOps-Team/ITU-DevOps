package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

func initDb(driver string, datasource string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(datasource), &gorm.Config{})
	sql, _ := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	return db, sql.Ping()
}




func init() {
	prometheus.MustRegister(minitwit_http_responses_total)
}

func main() {


	gorm, err := initDb(DRIVER, DATABASE)
	// sql, _ := gorm.DB()
	if err != nil {
		log.Fatal(err)
	}


	LoadTemplates()

	store := sessions.NewCookieStore([]byte(SECRET_KEY))

	r := mux.NewRouter()
	r.Use(BeforeRequestMiddleware(store, gorm))

	//CSS
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


	// r.Handle("/user/{id}", GetUserByIdHandler(gorm)).Methods("GET")

	http.ListenAndServe(":8080", r)
}
