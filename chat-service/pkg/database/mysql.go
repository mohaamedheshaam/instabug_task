package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func NewMySQLConnection(cfg MySQLConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	return db, nil
}
