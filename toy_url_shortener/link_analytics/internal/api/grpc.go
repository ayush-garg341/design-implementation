package api

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/linkanalytics/gen/analytics"
	"github.com/linkanalytics/internal/store"
)

// AnalyticsGRPCServer implements the generated AnalyticsServiceServer interface

type AnalyticsGRPCServer struct {
	pb.UnimplementedAnalyticsServiceServer
	store store.ClickStore
}

func NewAnalyticsGRPCServer(s store.ClickStore) *AnalyticsGRPCServer {
	return &AnalyticsGRPCServer{store: s}
}

// GetStats — called by link service when user requests their link stats
func (s *AnalyticsGRPCServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	if req.ShortCode == "" {
		return nil, status.Error(codes.InvalidArgument, "short_code is required")
	}

	// ctx already carries the deadline from the calling service
	// pass it straight to every I/O call — no new context needed
	stats, err := s.store.GetClickStats(ctx, req.ShortCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "no stats for this link")
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "db timeout")
		}

		log.Printf("GetStats db error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.GetStatsResponse{
		ShortCode:  stats.ShortCode,
		ClickCount: int64(stats.ClickCount),
		LongUrl:    stats.LongUrl,
		UserId:     stats.UserId,
	}, nil
}
