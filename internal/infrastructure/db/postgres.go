package db

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresConnection(cfg Config, sugar *zap.SugaredLogger) (*sql.DB, error) {
	url := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Path:     cfg.DBName,
		RawQuery: "sslmode=" + url.QueryEscape(cfg.SSLMode),
	}

	dsn := url.String()

	sugar.Debugw("Connecting to PostgreSQL database",
		"host", cfg.Host,
		"port", cfg.Port,
		"dbname", cfg.DBName,
		"sslmode", cfg.SSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	sugar.Infow("Successfully connected to PostgreSQL database",
		"host", cfg.Host,
		"port", cfg.Port,
		"dbname", cfg.DBName,
	)

	return db, nil
}
