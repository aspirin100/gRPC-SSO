package grpclient

import (
	"context"
	"fmt"
	"time"

	ssov1 "github.com/aspirin100/gRPC-SSO/protos/gen/go/sso"
	"github.com/google/uuid"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api ssov1.AuthClient
}

func New(
	ctx context.Context,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "grpclient.New"

	retryOpts := []retry.CallOption{
		retry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		retry.WithMax(uint(retriesCount)),
		retry.WithPerRetryTimeout(timeout),
	}

	// new client with retry interceptor
	cc, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			retry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: ssov1.NewAuthClient(cc),
	}, nil
}

func (cl *Client) IsAdmin(ctx context.Context, userID uuid.UUID) (*bool, error) {
	const op = "grpclient.IsAdmin"

	isAdmin, err := cl.api.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserID: userID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &isAdmin.IsAdmin, nil
}
