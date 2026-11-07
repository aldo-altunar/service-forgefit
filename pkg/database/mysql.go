package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ForgeFit_Servidor/pkg/config"
	mysqldriver "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection(cfg *config.Config) (*sql.DB, error) {
	driverConfig := mysqldriver.NewConfig()
	driverConfig.User = cfg.DBUser
	driverConfig.Passwd = cfg.DBPassword
	driverConfig.Net = "tcp"
	driverConfig.Addr = fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort)
	driverConfig.DBName = cfg.DBName
	driverConfig.ParseTime = true
	driverConfig.Loc = time.Local
	driverConfig.Timeout = 10 * time.Second
	driverConfig.ReadTimeout = 30 * time.Second
	driverConfig.WriteTimeout = 30 * time.Second

	db, err := sql.Open("mysql", driverConfig.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("open MySQL connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to MySQL at %s:%s: %w", cfg.DBHost, cfg.DBPort, err)
	}
	return db, nil
}
