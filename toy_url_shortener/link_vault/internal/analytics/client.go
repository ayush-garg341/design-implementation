package analytics

import (
	"context"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	pb "github.com/linkvault/gen/analytics"
)

// Client wraps the generated gRPC stub
// Define your own interface so handlers can be tested with a mock

type Client interface {
	GetStats(ctx context.Context, shortCode string) (*pb.GetStatsResponse, error)
	Close() error
}

type grpcClient struct {
	conn   *grpc.ClientConn
	client pb.AnalyticsServiceClient
}

func NewClient(addr string) (Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// keep-alive: ping server every 30s so connection stays alive
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30 * time.Second,
			Timeout: 10 * time.Second,
		}),
	)

	if err != nil {
		return nil, err
	}

	return &grpcClient{
		conn:   conn,
		client: pb.NewAnalyticsServiceClient(conn),
	}, nil
}

func (c *grpcClient) Close() error { return c.conn.Close() }

// GetStats — wraps the gRPC call, translates errors to Go errors
func (c *grpcClient) GetStats(ctx context.Context, shortCode string) (*pb.GetStatsResponse, error) {
	resp, err := c.client.GetStats(ctx, &pb.GetStatsRequest{
		ShortCode: shortCode,
	})
	if err != nil {
		return nil, translateGRPCError(err)
	}
	return resp, nil
}

// translateGRPCError — convert gRPC status codes to readable Go errors
func translateGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrNotFound
	case codes.DeadlineExceeded:
		return context.DeadlineExceeded
	case codes.Unavailable:
		return ErrUnavailable // analytics service is down
	default:
		log.Printf("analytics gRPC error: %s %s", st.Code(), st.Message())
		return ErrInternal
	}
}

var (
	ErrNotFound    = errors.New("analytics: not found")
	ErrUnavailable = errors.New("analytics: service unavailable")
	ErrInternal    = errors.New("analytics: internal error")
)
