package config

import (
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/models"
)

type GlobalConfig struct {
	ENV       string
	DB        *pgx.Conn
	DBQueries *database.Queries
	Logger    *slog.Logger
}

type GlobalAssets struct {
	Player *models.Player
	Items  map[string]*models.Item
}
