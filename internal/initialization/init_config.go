package initialization

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/config"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
)

func InitializeConfig() (*config.GlobalConfig, error) {
	// Environment
	env := os.Getenv("ENV")
	if env == "" {
		return nil, errors.New("ENV not set")
	}

	// DB URL
	dbURL, err := getDBURL()
	if err != nil {
		return nil, err
	}

	// DB Connection
	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %v", err)
	}

	// DB Queries
	dbQueries := database.New(db)

	// Logger

	cfg := config.GlobalConfig{
		ENV:       env,
		DB:        db,
		DBQueries: dbQueries,
	}

	return &cfg, nil
}

func getDBURL() (string, error) {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return "", errors.New("DB_NAME not set")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return "", errors.New("DB_USER not set")
	}

	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		return "", errors.New("DB_PASSWORD not set")
	}

	dbAddr := os.Getenv("DB_ADDR")
	if dbAddr == "" {
		return "", errors.New("DB_ADDR not set")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return "", errors.New("DB_PORT not set")
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbAddr,
		dbPort,
		dbName), nil
}
