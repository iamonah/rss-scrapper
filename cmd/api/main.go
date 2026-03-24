package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
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

	ctx, cancel := context.WithCancel(context.Background())
	
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		Scrapping(queries, 10, time.Minute, ctx)
		wg.Done()
		}()

	err = app.Serve()
	if err != nil {
		log.Printf("Server shutdown abruptly %v", err)
		return
	}
	
	cancel()
	wg.Wait()
	log.Println("application shutdown complete")
}
