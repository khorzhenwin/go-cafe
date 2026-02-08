package main

import (
	"log"

	"github.com/joho/godotenv"
	appconfig "github.com/khorzhenwin/go-cafe/backend/internal/config"
	"github.com/khorzhenwin/go-cafe/backend/internal/db"
	"github.com/khorzhenwin/go-cafe/backend/internal/server"
)

func (app *application) run() error {
	_ = godotenv.Load()

	cloudDbCfg, err := appconfig.LoadAWSConfig()
	if err != nil {
		log.Fatal(err)
	}

	authCfg, err := appconfig.LoadAuthConfig()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.NewAWSClient(cloudDbCfg)
	if err != nil {
		log.Fatal(err)
	}
	// Tables are created via migrations (make migrate-up). Do not AutoMigrate here.

	srvCfg := server.Config{
		BasePath:     app.config.BASE_PATH,
		Address:      app.config.ADDRESS,
		WriteTimeout: app.config.writeTimeout,
		ReadTimeout:  app.config.readTimeout,
	}
	handler := server.New(conn, authCfg, srvCfg)
	srv := server.NewServer(handler, srvCfg)

	log.Println("Starting server on", app.config.ADDRESS)
	return srv.ListenAndServe()
}
