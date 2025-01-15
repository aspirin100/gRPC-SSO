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
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Auth struct {
	logg        *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context,
		email string,
		passHash []byte) (userID string, err error)
}

type UserProvider interface {
	IsAdmin(ctx context.Context, userID string) (bool, error)
	GetUser(ctx context.Context, email string) (entity.User, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, appID int32) (entity.App, error)
}

func New(logg *slog.Logger,
	usrSaver UserSaver,
	usrProvider UserProvider,
	appProvider AppProvider,
	accessTTL,
	refreshTTL time.Duration) *Auth {

	return &Auth{
		logg:        logg,
		usrSaver:    usrSaver,
		usrProvider: usrProvider,
		appProvider: appProvider,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
	}
}

func (a *Auth) RegisterUser(ctx context.Context,
	email,
	password string) (
	string, error) {
	const op = "service/auth.Register"

	logg := a.logg.With(slog.String("op", op))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logg.Error("hashing error", sl.Err(err))
		return "", fmt.Errorf("password hashing error: %w", err)
	}

	userID, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		logg.Error("SaveUser error", sl.Err(err))
		return "", fmt.Errorf("failed to save new user: %w", err)
	}

	return userID, nil
}

func (a *Auth) Login(ctx context.Context,
	email,
	password string,
	appID int32) (
	*entity.TokenPair, error) {
	const op = "service/auth.Login"

	logg := a.logg.With(slog.String("op", op))

	user, err := a.usrProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			logg.Info("user not found", sl.Err(err))

			logg.Error("failed to get user", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		logg.Info("invalid credentials", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.GetApp(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
	}

	logg.Info("user successfully logged")

	accessToken, err := tokens.NewAccessToken(user, app, a.accessTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := tokens.NewRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &entity.TokenPair{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID string) (
	bool, error) {
	return false, nil
}
