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

const DRIVER = "sqlite3"
const DATABASE = "../db_backup/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"


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

	latest := 0

	store := sessions.NewCookieStore([]byte(SECRET_KEY))

	r := mux.NewRouter()
	
	r.Use(LatestMiddleware(&latest))
	r.Use(BeforeRequestMiddleware(store, gorm))


	//API ROUTES
	r.Handle("/latest", LatestHandler(&latest)).Methods("GET")
	r.Handle("/register", RegisterApiHandler(gorm)).Methods("POST")
	r.Handle("/msgs", MessagesHandler(store, gorm)).Methods("GET")
	r.Handle("/msgs/{username}", MessagesPerUserHandler(store, gorm)).Methods("GET", "POST")
	r.Handle("/fllws/{username}", FollowHandler(store, gorm)).Methods("GET", "POST")

	//http.ListenAndServe(":8080", r)
	log.Fatal(http.ListenAndServe(":8080", r))
}
