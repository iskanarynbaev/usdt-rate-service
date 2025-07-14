package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
	models "usdt_rate_service/internal/model"
	"usdt_rate_service/internal/repository"
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
	req, err := http.NewRequestWithContext(ctx, "GET", s.grinexURL, nil)
	if err != nil {
		return 0, 0, time.Time{}, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, 0, time.Time{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, 0, time.Time{}, &models.HTTPError{
			StatusCode: resp.StatusCode,
			Msg:        string(bodyBytes),
		}
	}

	var data grinexResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, time.Time{}, err
	}

	if len(data.Asks) == 0 || len(data.Bids) == 0 {
		return 0, 0, time.Time{}, errors.New("empty asks or bids")
	}

	askVal, ok := data.Asks[0][0].(float64)
	if !ok {
		return 0, 0, time.Time{}, errors.New("invalid ask format")
	}
	bidVal, ok := data.Bids[0][0].(float64)
	if !ok {
		return 0, 0, time.Time{}, errors.New("invalid bid format")
	}

	now := time.Now().UTC()

	err = s.repo.SaveRate(ctx, models.Rate{
		Ask:       askVal,
		Bid:       bidVal,
		Timestamp: now,
	})

	if err != nil {
		return 0, 0, time.Time{}, err
	}

	return askVal, bidVal, now, nil
}
