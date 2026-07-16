package config

import (
	"log/slog"

	"github.com/jackc/pgx/v5"
)

var Cfg GlobalConfig

type GlobalConfig struct {
	ENV    string
	DB     *pgx.Conn
	Logger *slog.Logger
}

func init() {

}
