package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback for local testing if not set
		dsn = "host=localhost port=5432 user=postgres password=postgres dbname=animedb sslmode=disable"
		log.Println("DATABASE_URL not set, using default:", dsn)
	}

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to PostgreSQL successfully")
}
