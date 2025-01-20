package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/aspirin100/gRPC-SSO/internal/entity"
	"github.com/aspirin100/gRPC-SSO/internal/storage"
	"github.com/aspirin100/gRPC-SSO/internal/tokens"
	"github.com/aspirin100/gRPC-SSO/pkg/logger/sl"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidPassword      = errors.New("wrong password")
	ErrUserExists           = errors.New("user already exists")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
)

type Auth struct {
	logg        *slog.Logger
	authManager AuthManager
	secretKey   string
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

type AuthManager interface {
	UserSaver
	UserProvider
	AppProvider
	RefreshSessionManager
}

// storage interfaces.
type UserSaver interface {
	SaveUser(ctx context.Context,
		email string,
		passHash []byte) (userID string, err error)
}

type UserProvider interface {
	IsAdmin(ctx context.Context, userID string) (*bool, error)
	GetUser(ctx context.Context, email string) (*entity.User, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, appID int32) (*entity.App, error)
}

type RefreshSessionManager interface {
	NewRefreshSession(ctx context.Context,
		refreshToken, userID string, refreshTTL time.Duration) error
	ValidateRefreshToken(ctx context.Context, userID, refreshToken string) error
}

func New(logg *slog.Logger,
	authManager AuthManager,
	accessTTL,
	refreshTTL time.Duration,
	secretKey string) *Auth {
	return &Auth{
		logg:        logg,
		authManager: authManager,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
		secretKey:   secretKey,
	}
}

func (a *Auth) RegisterUser(ctx context.Context,
	email,
	password string) (
	*string, error) {
	const op = "service/auth.Register"

	logg := a.logg.With(slog.String("op", op))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logg.Error("hashing error", sl.Err(err))

		return nil, fmt.Errorf("password hashing error: %w", err)
	}

	userID, err := a.authManager.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			logg.Warn("user already exists", sl.Err(err))

			return nil, ErrUserExists //nolint:wrapcheck
		}

		logg.Error("SaveUser error", sl.Err(err))

		return nil, fmt.Errorf("failed to save new user: %w", err)
	}

	return &userID, nil
}

func (a *Auth) Login(ctx context.Context,
	email,
	password string,
	appID int32) (
	*entity.TokenPair, error) {
	const op = "service/auth.Login"

	logg := a.logg.With(slog.String("op", op))

	user, err := a.authManager.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			logg.Info("user not found", sl.Err(err))
		}

		logg.Error("failed to get user", sl.Err(err))

		return nil, ErrInvalidCredentials //nolint:wrapcheck
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		logg.Info("invalid credentials", sl.Err(err))

		return nil, ErrInvalidPassword //nolint:wrapcheck
	}

	app, err := a.authManager.GetApp(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			logg.Error("app not found", sl.Err(err))

			return nil, fmt.Errorf("%s: %w", op, err)
		}

		logg.Error("failed to get app", sl.Err(err))

		return nil, fmt.Errorf("failed to get app: %w", err)
	}

	logg.Info("user successfully logged")

	accessToken, err := tokens.NewAccessToken(user.UserID, app.ID, a.accessTTL, a.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := tokens.NewRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	// inserts new refresh token into database (refresh_session table)
	err = a.authManager.NewRefreshSession(ctx, user.UserID,
		*refreshToken, a.refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &entity.TokenPair{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID string) (
	*bool, error) {
	const op = "service/auth.IsAdmin"

	logg := a.logg.With(slog.String("op", op))

	isAdmin, err := a.authManager.IsAdmin(ctx, userID)
	if err != nil {
		logg.Error("checking if user is admin failed", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (a *Auth) RefreshTokenPair(
	ctx context.Context,
	userID, refreshToken string,
	appID int32) (*entity.TokenPair, error) {
	const op = "service/auth.RefreshTokenPair"

	logg := a.logg.With(slog.String("op", op))

	err := a.authManager.ValidateRefreshToken(ctx, userID, refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRefreshTokenNotFound):
			logg.Info("refresh token not found", sl.Err(err))

			return nil, ErrRefreshTokenNotFound //nolint:wrapcheck
		case errors.Is(err, tokens.ErrInvalidRefreshToken):
			logg.Info("refresh token is invalid", sl.Err(err))

			return nil, ErrInvalidRefreshToken //nolint:wrapcheck
		default:
			logg.Error("validate refresh token error", sl.Err(err))

			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	accessToken, err := tokens.NewAccessToken(userID, appID, a.accessTTL, a.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	newRefreshToken, err := tokens.NewRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	// inserts new refresh token into database (refresh_session table)
	err = a.authManager.NewRefreshSession(ctx,
		*newRefreshToken, userID, a.refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &entity.TokenPair{
		AccessToken:  *accessToken,
		RefreshToken: *newRefreshToken,
	}, nil
}
