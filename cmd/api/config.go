package main

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	port string
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

func loadConfig() *Config {
	maxOpenConns, err := strconv.Atoi(os.Getenv("MAXOPENCONNS"))
	if err != nil {
		log.Fatalf("Invalid MAXOPENCONNS: %v", err)
	}

	maxIdleConns, err := strconv.Atoi(os.Getenv("MAXIDLECONNS"))
	if err != nil {
		log.Fatalf("Invalid MAXIDLECONNS: %v", err)
	}

	return &Config{
		port: os.Getenv("PORT"),
		env:  os.Getenv("ENV"),
		db: struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  string
		}{
			dsn:          os.Getenv("DSN"),
			maxOpenConns: maxOpenConns,
			maxIdleConns: maxIdleConns,
			maxIdleTime:  os.Getenv("MAXIDLETIME"),
		},
	}
}
