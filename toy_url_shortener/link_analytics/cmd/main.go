package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/linkanalytics/config"
	pb "github.com/linkanalytics/gen/analytics"
	"github.com/linkanalytics/internal/api"
	"github.com/linkanalytics/internal/consumer"
	"github.com/linkanalytics/internal/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	// Load the config
	cfg := config.Config()

	// Create the db connection
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	clickStore := store.NewClickStore(db)

	consumer := consumer.NewClickConsumer(
		cfg.KafkaBrokers,
		"analytics",
		"analytics-service",
		clickStore,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Start(ctx)
		log.Println("kafka consumer stopped")
	}()

	// Create grpc Server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor,
			recoveryInterceptor,
		),
	)

	analyticsHandler := api.NewAnalyticsGRPCServer(*clickStore)
	pb.RegisterAnalyticsServiceServer(grpcServer, analyticsHandler)

	// ── start gRPC listener ───────────────────────────────
	lis, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatal(err)
	}

	// start gRPC in its own goroutine — it blocks inside Serve()
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("gRPC server listening on :8085")
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit
	log.Println("shutdown signal received")

	grpcServer.GracefulStop()

	log.Println("closing database")
	db.Close()

	log.Println("closing consumer")
	cancel()
	consumer.Close()
	wg.Wait()
	log.Println("server exiting")

}

// loggingInterceptor — logs every gRPC call with duration
func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("grpc %s  dur=%s  err=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}

// recoveryInterceptor — catches panics, returns clean gRPC error
func recoveryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in %s: %v", info.FullMethod, r)
			err = status.Error(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}
