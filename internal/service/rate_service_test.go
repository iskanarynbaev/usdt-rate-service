package service

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"usdt_rate_service/internal/mocks"
	model "usdt_rate_service/internal/model"
)

func TestFetchRates_Success(t *testing.T) {
	ctx := context.Background()

	mockHTTP := &mocks.HTTPClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			jsonResp := `{
				"asks": [[0.95, 1]],
				"bids": [[0.93, 1]]
			}`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(jsonResp)),
			}, nil
		},
	}

	var savedRate model.Rate
	mockRepo := &mocks.RepositoryMock{
		SaveRateFunc: func(ctx context.Context, rate model.Rate) error {
			savedRate = rate
			return nil
		},
	}

	svc := &Service{
		repo:       mockRepo,
		httpClient: mockHTTP,
		grinexURL:  "http://fakeurl",
	}

	ask, bid, ts, err := svc.FetchRates(ctx)
	if err != nil {
		t.Fatalf("FetchRates failed: %v", err)
	}

	if ask != 0.95 {
		t.Errorf("Expected ask 0.95, got %v", ask)
	}
	if bid != 0.93 {
		t.Errorf("Expected bid 0.93, got %v", bid)
	}
	if savedRate.Ask != 0.95 || savedRate.Bid != 0.93 {
		t.Errorf("SaveRate got wrong data: %+v", savedRate)
	}
	if savedRate.Timestamp.IsZero() {
		t.Errorf("SaveRate timestamp is zero")
	}
	if ts.IsZero() {
		t.Errorf("Returned timestamp is zero")
	}
}

func TestFetchRates_HTTPError(t *testing.T) {
	ctx := context.Background()

	mockHTTP := &mocks.HTTPClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, http.ErrHandlerTimeout
		},
	}

	svc := &Service{
		repo:       &mocks.RepositoryMock{},
		httpClient: mockHTTP,
		grinexURL:  "http://fakeurl",
	}

	_, _, _, err := svc.FetchRates(ctx)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestFetchRates_InvalidJSON(t *testing.T) {
	ctx := context.Background()

	mockHTTP := &mocks.HTTPClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(`invalid json`)),
			}, nil
		},
	}

	svc := &Service{
		repo:       &mocks.RepositoryMock{},
		httpClient: mockHTTP,
		grinexURL:  "http://fakeurl",
	}

	_, _, _, err := svc.FetchRates(ctx)
	if err == nil {
		t.Fatal("Expected JSON decode error, got nil")
	}
}
