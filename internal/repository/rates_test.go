package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	pgconn "github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	model "usdt_rate_service/internal/model"
)

type MockDB struct {
	ExecFunc func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

func (m *MockDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, sql, args...)
	}
	return pgconn.NewCommandTag(""), nil
}

func TestSaveRate_Success(t *testing.T) {
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag(""), nil
		},
	}

	repo := NewPostgresRepo(mockDB)
	rate := model.Rate{
		Ask:       1.23,
		Bid:       1.22,
		Timestamp: time.Now(),
	}

	err := repo.SaveRate(context.Background(), rate)
	assert.NoError(t, err)
}

func TestSaveRate_Error(t *testing.T) {
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag(""), errors.New("some db error")
		},
	}

	repo := NewPostgresRepo(mockDB)
	rate := model.Rate{
		Ask:       1.23,
		Bid:       1.22,
		Timestamp: time.Now(),
	}

	err := repo.SaveRate(context.Background(), rate)
	assert.Error(t, err)
}
