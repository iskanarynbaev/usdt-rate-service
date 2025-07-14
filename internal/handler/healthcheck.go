package handler

import (
	"context"
	pb "usdt_rate_service/internal/grpc"
)

func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: "SERVING",
	}, nil
}
