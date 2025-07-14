package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	models "usdt_rate_service/internal/model"
)

type DBExec interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type Repository interface {
	SaveRate(ctx context.Context, rate models.Rate) error
}

type PostgresRepo struct {
	db DBExec
}

func NewPostgresRepo(db DBExec) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) SaveRate(ctx context.Context, rate models.Rate) error {
	sql := `INSERT INTO rates (ask, bid, timestamp) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, sql, rate.Ask, rate.Bid, rate.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert rate: %w", err)
	}
	return nil
}
