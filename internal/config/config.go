package config

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
)

var Cfg GlobalConfig

type GlobalConfig struct {
	ENV       string
	DB        *pgx.Conn
	DBQueries *database.Queries
	Logger    *slog.Logger
}

func init() {
	// Get Environmental Variables
	env := os.Getenv("ENV")
	if env == "" {
		log.Fatal("ENV not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME not set")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER not set")
	}

	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		log.Fatal("DB_PASS not set")
	}

	dbAddr := os.Getenv("DB_ADDR")
	if dbAddr == "" {
		log.Fatal("DB_ADDR not set")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Fatal("DB_PORT not set")
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbAddr,
		dbPort,
		dbName)

	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v\n", err)
	}

	dbQueries := database.New(db)

	Cfg = GlobalConfig{
		ENV:       env,
		DB:        db,
		DBQueries: dbQueries,
	}
}
