package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/onahvictor/rss-scrapper/internal/database"
)

type application struct {
	config *Config
	DB     *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .ENV")
	}

	config := loadConfig()

	db, err := sql.Open("postgres", config.db.dsn)
	if err != nil {
		log.Fatal("database not connected")
	}

	log.Println("database connection esterblished")
	queries := database.New(db)
	app := &application{
		config: config,
		DB:     queries,
	}

	go startScrapping(queries, 10, time.Minute)
	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
