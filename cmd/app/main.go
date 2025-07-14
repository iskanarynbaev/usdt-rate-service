package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"usdt_rate_service/internal/config"
	pb "usdt_rate_service/internal/grpc"
	"usdt_rate_service/internal/handler"
	"usdt_rate_service/internal/repository"
	"usdt_rate_service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pool.Close()

	repo := repository.NewPostgresRepo(pool)
	svc := service.NewService(repo, cfg.GrinexURL)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRateServiceServer(grpcServer, handler.NewServer(svc))

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down GRPC server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("Starting GRPC server on port %s...", cfg.GRPCPort)
	log.Printf("Work URL is %s...", cfg.GrinexURL)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
