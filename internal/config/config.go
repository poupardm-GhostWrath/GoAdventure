package config

import (
	"log/slog"

	"github.com/google/uuid"
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
	ID        uuid.UUID
	Player    *models.Player
	Items     map[int32]*models.Item     // map[itemID]*models.Item
	Locations map[int32]*models.Location // map[locationID]*models.Location
	Stores    map[int32]*models.Store    // map[locationID]*models.Store
}
