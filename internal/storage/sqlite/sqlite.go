package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aspirin100/gRPC-SSO/internal/storage"
	"github.com/aspirin100/gRPC-SSO/pkg/logger/sl"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(logg *slog.Logger, storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	log := logg.With(slog.String("op", op))

	db, err := sql.Open("sqlite3", storagePath)
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

	queryStatement, err := s.db.Prepare(SaveUserQuery)
	if err != nil {
		return "", fmt.Errorf("failed to prepare query: %w", err)
	}

	_, err = queryStatement.ExecContext(ctx, userID, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) &&
			sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return "", fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return "", fmt.Errorf("failed to save user: %w", err)
	}

	return userID, nil
}



const (
	SaveUserQuery = `insert into users(userID, email, passHash) values($1, $2, $3)`
	GetUserQuery = `select (userID, email, passHash) from users where email = $1`
)
