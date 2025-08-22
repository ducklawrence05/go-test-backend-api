package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ducklawrence05/go-test-backend-api/internal/initialize"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate [up|down|step n]")
	}

	cfg, err := initialize.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	pgCfg := cfg.Postgres
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		pgCfg.Username, pgCfg.Password, pgCfg.Host, pgCfg.Port, pgCfg.Dbname,
	)

	migrationsPath := "file://./migrate/migrations"

	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	cmd := os.Args[1]
	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "step":
		if len(os.Args) < 3 {
			log.Fatal("Usage: migrate step <n>")
		}
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid step count: %v", err)
		}
		if err := m.Steps(n); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown command: %s", cmd)
	}
}
