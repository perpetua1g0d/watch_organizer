package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DBConfig struct {
	Host string
	Port string
	Username string
	Password string
	DBName string
	SSLMode string
}

func NewPostgresDB(cfg DBConfig) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
}