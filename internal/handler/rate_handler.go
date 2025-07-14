package handler

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	model "usdt_rate_service/internal/model"
	"usdt_rate_service/internal/service"
	logger "usdt_rate_service/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "usdt_rate_service/internal/grpc"
)

type Server struct {
	pb.UnimplementedRateServiceServer
	svc *service.Service
}

func NewServer(svc *service.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) GetRates(ctx context.Context, req *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
	if err := logger.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Log.Info("Received GetRates request")

	ask, bid, ts, err := s.svc.FetchRates(ctx)
	if err != nil {
		var httpErr *model.HTTPError
		if errors.As(err, &httpErr) {
			logger.Log.Warn("FetchRates returned HTTP error",
				zap.Int("status_code", httpErr.StatusCode),
				zap.String("message", httpErr.Msg),
			)

			switch httpErr.StatusCode {
			case 422:
				return nil, status.Errorf(codes.InvalidArgument, "failed to fetch rates: %v", httpErr)
			case 500, 502, 503, 504:
				return nil, status.Errorf(codes.Unavailable, "failed to fetch rates: %v", httpErr)
			default:
				return nil, status.Errorf(codes.Internal, "failed to fetch rates: %v", httpErr)
			}
		}

		logger.Log.Error("FetchRates failed with unexpected error", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to fetch rates: %v", err)
	}

	logger.Log.Info("GetRates succeeded",
		zap.Float64("ask", ask),
		zap.Float64("bid", bid),
		zap.Int64("timestamp", ts.Unix()),
	)

	return &pb.GetRatesResponse{
		Ask:       ask,
		Bid:       bid,
		Timestamp: ts.Unix(),
	}, nil
}
