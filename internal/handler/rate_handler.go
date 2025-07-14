package handler

import (
	"context"
	"errors"
	model "usdt_rate_service/internal/model"
	"usdt_rate_service/internal/service"

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
	ask, bid, ts, err := s.svc.FetchRates(ctx)
	if err != nil {
		var httpErr *model.HTTPError
		if errors.As(err, &httpErr) {
			switch httpErr.StatusCode {
			case 422:
				return nil, status.Errorf(codes.InvalidArgument, "failed to fetch rates: %v", httpErr)
			case 500, 502, 503, 504:
				return nil, status.Errorf(codes.Unavailable, "failed to fetch rates: %v", httpErr)
			default:
				return nil, status.Errorf(codes.Internal, "failed to fetch rates: %v", httpErr)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch rates: %v", err)
	}

	return &pb.GetRatesResponse{
		Ask:       ask,
		Bid:       bid,
		Timestamp: ts.Unix(),
	}, nil
}
