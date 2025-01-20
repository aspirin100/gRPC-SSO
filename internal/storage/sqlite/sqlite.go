package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/aspirin100/gRPC-SSO/internal/entity"
	"github.com/aspirin100/gRPC-SSO/internal/storage"
	"github.com/aspirin100/gRPC-SSO/internal/tokens"
	"github.com/aspirin100/gRPC-SSO/pkg/logger/sl"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sqlx.DB
}

func New(logg *slog.Logger, storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	log := logg.With(slog.String("op", op))

	db, err := sqlx.Open("sqlite3", storagePath)
	if err != nil {
		log.Error("db open error", sl.Err(err))
		
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveUser(ctx context.Context,
	email string,
	passHash []byte) (userID string, err error) {
	userID = uuid.NewString()

	const op = "storage.sqlite3.SaveUser"

	_, err = s.db.ExecContext(
		ctx,
		SaveUserQuery,
		userID,
		email,
		passHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) &&
			sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return "", storage.ErrUserExists
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (s *Storage) GetUser(ctx context.Context, email string) (*entity.User, error) {
	const op = "storage.sqlite.GetUser"

	user := &entity.User{}

	err := s.db.GetContext(ctx, user, GetUserQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID string) (*bool, error) {
	const op = "storage.sqlite.GetUser"

	var isAdmin bool

	err := s.db.GetContext(ctx, &isAdmin, IsAdminQuery, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &isAdmin, nil
}

func (s *Storage) GetApp(ctx context.Context, appID int32) (*entity.App, error) {
	const op = "storage.sqlite.GetUser"

	app := entity.App{}

	err := s.db.GetContext(ctx, &app, GetAppQuery, appID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrAppNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &app, nil
}

func (s *Storage) NewRefreshSession(
	ctx context.Context,
	refreshToken, userID string,
	refreshTTL time.Duration) error {
	const op = "storage.sqlite.NewRefreshSession"

	_, err := s.db.ExecContext(
		ctx,
		NewRefreshSessionQuery,
		refreshToken,
		userID,
		time.Now().Add(refreshTTL).Unix())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ValidateRefreshToken(ctx context.Context, refreshToken, userID string) error {
	const op = "storage.sqlite.ValidateRefreshToken"

	result := struct {
		ExpiresAt int64 `db:"expiresAt"`
		IsUsed    bool  `db:"isUsed"`
	}{}

	err := s.db.GetContext(ctx,
		&result,
		ValidateRefreshTokenQuery,
		refreshToken, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrRefreshTokenNotFound
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	if result.IsUsed {
		return tokens.ErrInvalidRefreshToken
	}

	if time.Now().Unix() >= result.ExpiresAt {
		return tokens.ErrInvalidRefreshToken
	}

	_, err = s.db.ExecContext(ctx, SetRefreshTokenUsedQuery, refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

const (
	SaveUserQuery             = `insert into users(id, email, passHash) values(?, ?, ?)`
	GetUserQuery              = `select id, email, passHash from users where email = ?`
	IsAdminQuery              = `select isAdmin from users where id = ?`
	GetAppQuery               = `select id, name from apps where id = ?`
	ValidateRefreshTokenQuery = `select expiresAt, isUsed from
	refresh_session where refreshToken = ? AND userID = ?`
	SetRefreshTokenUsedQuery = `update refresh_session set isUsed = true where refreshToken = ?`
	NewRefreshSessionQuery   = `insert into
	refresh_session(refreshToken, userID, expiresAt)
	values(?, ?, ?)`
)
