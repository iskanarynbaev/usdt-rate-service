package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"usdt_rate_service/internal/config"
	pb "usdt_rate_service/internal/grpc"
	"usdt_rate_service/internal/handler"
	"usdt_rate_service/internal/repository"
	"usdt_rate_service/internal/service"
	"usdt_rate_service/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// Init zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize zap logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("logger initialized")

	// Init tracer
	tp, err := utils.InitTracer()
	if err != nil {
		logger.Fatal("failed to initialize tracer", zap.Error(err))
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Fatal("failed to shutdown tracer provider", zap.Error(err))
		}
	}()

	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Fatal("failed to connect to DB", zap.Error(err))
	}
	defer pool.Close()

	repo := repository.NewPostgresRepo(pool)
	svc := service.NewService(repo, cfg.GrinexURL)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRateServiceServer(grpcServer, handler.NewServer(svc))

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		logger.Info("Shutting down GRPC server...")
		grpcServer.GracefulStop()
	}()

	logger.Info("Starting GRPC server",
		zap.String("port", cfg.GRPCPort),
		zap.String("grinex_url", cfg.GrinexURL),
	)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve gRPC", zap.Error(err))
	}
}
