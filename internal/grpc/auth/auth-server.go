package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/aspirin100/gRPC-SSO/internal/entity"
	authService "github.com/aspirin100/gRPC-SSO/internal/service/auth"
	ssov1 "github.com/aspirin100/gRPC-SSO/protos/gen/go/sso"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue     = 0
	passwordMaxLen = 16
)

// service layer interface.
type Auth interface {
	Login(ctx context.Context, email, password string, appID int32) (
		*entity.TokenPair, error)
	RegisterUser(ctx context.Context, email, password string) (*string, error)
	IsAdmin(ctx context.Context, userID string) (*bool, error)
	RefreshTokenPair(ctx context.Context,
		userID, refreshToken string, appID int32) (*entity.TokenPair, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func RegisterAuthServer(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) ( //nolint:dupl
	*ssov1.NewTokenPairResponse, error) {
	err := validateLogin(req)
	if err != nil {
		return nil, fmt.Errorf("login validation error: %w", err)
	}

	tokens, err := s.auth.Login(ctx, req.GetEmail(),
		req.GetPassword(), req.GetAppID())
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "wrong email or password") 
		case errors.Is(err, authService.ErrInvalidPassword):
			return nil, status.Error(codes.InvalidArgument, "wrong password") 
		default:
			return nil, status.Error(codes.Internal, "internal error") 
		}
	}

	return &ssov1.NewTokenPairResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (
	*ssov1.RegisterRespons, error) {
	err := validateRegister(req)
	if err != nil {
		return nil, fmt.Errorf("register validation error: %w", err)
	}

	userID, err := s.auth.RegisterUser(ctx, req.GetEmail(),
		req.GetPassword())
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrUserExists):
			return nil, status.Error(codes.AlreadyExists, "user already exists") 
		default:
			return nil, status.Error(codes.Internal, "internal error") 
		}
	}

	return &ssov1.RegisterRespons{
		UserID: *userID,
	}, nil
}

func (s *serverAPI) RefreshTokenPair(ctx context.Context, req *ssov1.RefreshRequest) ( //nolint:dupl
	*ssov1.NewTokenPairResponse, error) {
	err := validateRefreshRequest(req)
	if err != nil {
		return nil, fmt.Errorf("login validation error: %w", err)
	}

	tokens, err := s.auth.RefreshTokenPair(ctx, req.GetUserID(),
		req.GetRefreshToken(), req.GetAppID())
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrRefreshTokenNotFound):
			return nil, status.Error(codes.NotFound, "refresh token not found") 
		case errors.Is(err, authService.ErrInvalidRefreshToken):
			return nil, status.Error(codes.PermissionDenied, "invalid refresh token") 
		default:
			return nil, status.Error(codes.Internal, "internal error") 
		}
	}

	return &ssov1.NewTokenPairResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (
	*ssov1.IsAdminResponse, error) {
	if req.GetUserID() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required") 
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserID())
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "wrong email or password") 
		default:
			return nil, status.Error(codes.Internal, "internal error") 
		}
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: *isAdmin,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	} else if len(req.GetPassword()) > passwordMaxLen {
		return status.Error(codes.InvalidArgument, "password is too big")
	}

	if req.GetAppID() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	} else if len(req.GetPassword()) > passwordMaxLen {
		return status.Error(codes.InvalidArgument, "password is too big")
	}

	return nil
}

func validateRefreshRequest(req *ssov1.RefreshRequest) error {
	if req.GetRefreshToken() == "" {
		return status.Error(codes.InvalidArgument, "refresh token is required")
	}

	if req.GetUserID() == "" {
		return status.Error(codes.InvalidArgument, "user id is required")
	}

	return nil
}