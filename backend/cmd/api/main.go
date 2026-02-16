package main

import (
	"log"
	"time"
)

// @title go-cafe backend API
// @version 1.0
// @description API for users, cafe listings, and ratings.
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

type application struct {
	config config
}

type config struct {
	BASE_PATH    string
	ADDRESS      string
	writeTimeout time.Duration
	readTimeout  time.Duration
}

func main() {
	cfg := config{
		BASE_PATH:    "/api/v1",
		ADDRESS:      ":8080",
		writeTimeout: time.Second * 10,
		readTimeout:  time.Second * 5,
	}

	app := &application{
		config: cfg,
	}

	log.Fatal(app.run())
}
