package db

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"wb-tech/internal/config"
)

func OpenDB(ctx context.Context, cfg config.DBConfig) (*sqlx.DB, error) {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PgHost, cfg.PgPort, cfg.PgUser, cfg.PgPassword, cfg.PgDatabase,
	)

	db, err := sqlx.Open("pgx", dbInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database sqlx: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}
