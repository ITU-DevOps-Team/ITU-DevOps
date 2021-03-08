package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DRIVER = "sqlite3"
const DATABASE = "../db_backup/minitwit.db"
const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

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

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=enable", dbHost, dbUser, dbPass, dbName, dbPort), err
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

	r := mux.NewRouter()
	r.Use(LatestMiddleware(gorm))
	r.Use(AuthenticationMiddleware())
	r.Handle("/latest", LatestHandler(gorm)).Methods("GET")
	r.Handle("/register", RegisterApiHandler(gorm)).Methods("POST")
	r.Handle("/msgs", MessagesHandler(gorm)).Methods("GET")
	r.Handle("/msgs/{username}", MessagesPerUserHandler(gorm)).Methods("GET", "POST")
	r.Handle("/fllws/{username}", FollowHandler(gorm)).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
