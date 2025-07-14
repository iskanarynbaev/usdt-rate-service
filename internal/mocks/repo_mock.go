// mocks/repository.go
package mocks

import (
	"context"
	model "usdt_rate_service/internal/model"
)

type RepositoryMock struct {
	SaveRateFunc    func(ctx context.Context, rate model.Rate) error
	GetLastRateFunc func(ctx context.Context) (model.Rate, error)
}

func (m *RepositoryMock) SaveRate(ctx context.Context, rate model.Rate) error {
	return m.SaveRateFunc(ctx, rate)
}

func (m *RepositoryMock) GetLastRate(ctx context.Context) (model.Rate, error) {
	return m.GetLastRateFunc(ctx)
}
