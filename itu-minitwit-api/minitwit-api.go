package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DRIVER = "sqlite3"
const DATABASE = "../db_backup/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

var Latest_ = Latest{
	Latest: 0,
}

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

	r := mux.NewRouter()
	r.Use(LatestMiddleware())
	r.Use(AuthenticationMiddleware(gorm))
	r.Handle("/latest", LatestHandler()).Methods("GET")
	r.Handle("/register", RegisterApiHandler(gorm)).Methods("POST")
	r.Handle("/msgs", MessagesHandler(gorm)).Methods("GET")
	r.Handle("/msgs/{username}", MessagesPerUserHandler(gorm)).Methods("GET", "POST")
	r.Handle("/fllws/{username}", FollowHandler(gorm)).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
