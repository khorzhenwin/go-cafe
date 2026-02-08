// Run migrations: go run cmd/migrate/main.go [up|down]
// Requires DB_* env vars (same as API). Migrations live in migrations/.
package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	appconfig "github.com/khorzhenwin/go-cafe/backend/internal/config"
)

func main() {
	_ = godotenv.Load()

	dir := flag.String("dir", "migrations", "migrations directory (relative to cwd)")
	flag.Parse()
	action := "up"
	if flag.NArg() > 0 {
		action = flag.Arg(0)
	}

	cfg, err := appconfig.LoadAWSConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("migrations dir: %v", err)
	}
	sourceURL := "file://" + filepath.ToSlash(absDir)
	dbURL := cfg.GetMigrationDSN()

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		log.Fatalf("migrate init: %v", err)
	}
	defer m.Close()

	switch action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate up: %v", err)
		}
		log.Println("migrate: up done")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate down: %v", err)
		}
		log.Println("migrate: down done")
	default:
		log.Fatalf("usage: migrate [up|down]")
	}
}
