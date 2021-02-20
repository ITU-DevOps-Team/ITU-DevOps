package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDb(driver string, datasource string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(datasource), &gorm.Config{})
	sql, _ := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	return db, sql.Ping()
}

func main() {
	gorm, err := initDb(DRIVER, DATABASE)
	// sql, _ := gorm.DB()
	if err != nil {
		log.Fatal(err)
	}

	store := sessions.NewCookieStore([]byte(SECRET_KEY))

	r := mux.NewRouter()
	r.Use(BeforeRequestMiddleware(store, gorm))
	r.Handle("/", HomeHandler()).Methods("GET")
<<<<<<< HEAD
	// r.Handle("/public", TestHandler(sql)).Methods("GET")
	r.Handle("/login", LoginHandler(store, gorm)).Methods("GET", "POST")
	r.Handle("/register", RegisterHandler(store, gorm)).Methods("GET", "POST")
	r.Handle("/logout", LogoutHandler(store, gorm)).Methods("GET")
	r.Handle("/add_message", AddMessageHandler(store, gorm)).Methods("POST")
	// r.Handle("/{username}", TestHandler(sql)).Methods("GET")
	// r.Handle("/{username}/follow", TestHandler(sql)).Methods("GET")
	// r.Handle("/{username}/unfollow", TestHandler(sql)).Methods("GET")
	// r.Handle("/user/{id}", GetUserByIdHandler(gorm)).Methods("GET")
=======
	r.Handle("/public", TestHandler(db)).Methods("GET")
	r.Handle("/login", LoginHandler(store, db)).Methods("GET", "POST")
	r.Handle("/register", RegisterHandler(store, db)).Methods("GET", "POST")
	r.Handle("/logout", LogoutHandler(store, db)).Methods("GET")
	r.Handle("/add_message", TestHandler(db)).Methods("POST")
	r.Handle("/{username}", TestHandler(db)).Methods("GET")
	r.Handle("/{username}/follow", FollowUserHandler(store, db)).Methods("GET")
	r.Handle("/{username}/unfollow", UnfollowUserHandler(store, db)).Methods("GET")
	r.Handle("/test", TestHandler(db)).Methods("GET")
	r.Handle("/user/{id}", GetUserByIdHandler(db)).Methods("GET")
	r.Handle("/get_message", GetMessageByString(store, db)).Methods("GET")
>>>>>>> e92eb79e232baa0a94276e56bc8dd81cb0095aaa

	http.ListenAndServe(":8080", r)
}
