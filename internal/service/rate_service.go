package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"time"
	models "usdt_rate_service/internal/model"
	"usdt_rate_service/internal/repository"
	logger "usdt_rate_service/internal/utils"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Service struct {
	repo       repository.Repository
	httpClient HTTPClient
	grinexURL  string
}

func NewService(repo repository.Repository, grinexURL string) *Service {
	return &Service{
		repo:       repo,
		httpClient: &http.Client{Timeout: 5 * time.Second},
		grinexURL:  grinexURL,
	}
}

type grinexResponse struct {
	Asks [][]interface{} `json:"asks"`
	Bids [][]interface{} `json:"bids"`
}

func (s *Service) FetchRates(ctx context.Context) (ask, bid float64, timestamp time.Time, err error) {
	if err := logger.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Log.Info("Starting FetchRates", zap.String("url", s.grinexURL))

	req, err := http.NewRequestWithContext(ctx, "GET", s.grinexURL, nil)
	if err != nil {
		logger.Log.Error("Failed to create HTTP request", zap.Error(err))
		return 0, 0, time.Time{}, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		logger.Log.Error("HTTP request to Grinex failed", zap.Error(err))
		return 0, 0, time.Time{}, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.Log.Warn("Failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Log.Warn("Received non-200 response from Grinex",
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", string(bodyBytes)),
		)
		return 0, 0, time.Time{}, &models.HTTPError{
			StatusCode: resp.StatusCode,
			Msg:        string(bodyBytes),
		}
	}

	var data grinexResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Log.Error("Failed to decode Grinex response", zap.Error(err))
		return 0, 0, time.Time{}, err
	}

	if len(data.Asks) == 0 || len(data.Bids) == 0 {
		logger.Log.Warn("Received empty asks or bids")
		return 0, 0, time.Time{}, errors.New("empty asks or bids")
	}

	askVal, ok := data.Asks[0][0].(float64)
	if !ok {
		logger.Log.Error("Invalid format for ask value")
		return 0, 0, time.Time{}, errors.New("invalid ask format")
	}
	bidVal, ok := data.Bids[0][0].(float64)
	if !ok {
		logger.Log.Error("Invalid format for bid value")
		return 0, 0, time.Time{}, errors.New("invalid bid format")
	}

	now := time.Now().UTC()

	err = s.repo.SaveRate(ctx, models.Rate{
		Ask:       askVal,
		Bid:       bidVal,
		Timestamp: now,
	})
	if err != nil {
		logger.Log.Error("Failed to save rate to DB", zap.Error(err))
		return 0, 0, time.Time{}, err
	}

	logger.Log.Info("Successfully fetched and saved rate",
		zap.Float64("ask", askVal),
		zap.Float64("bid", bidVal),
		zap.Time("timestamp", now),
	)

	return askVal, bidVal, now, nil
}
