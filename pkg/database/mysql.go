package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"ForgeFit_Servidor/pkg/config"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLConnection(cfg *config.Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("database configuration is nil")
	}

	dbConfig := mysql.Config{
		User:      cfg.DBUser,
		Passwd:    cfg.DBPassword,
		Net:       "tcp",
		Addr:      cfg.DBHost + ":" + cfg.DBPort,
		DBName:    cfg.DBName,
		Params:    map[string]string{"charset": "utf8mb4"},
		ParseTime: true,
		Loc:       time.UTC,
	}
	dsn := dbConfig.FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open MySQL connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping MySQL: %w", err)
	}

	return db, nil
}
