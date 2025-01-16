package storage_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/aspirin100/gRPC-SSO/internal/storage"
	"github.com/aspirin100/gRPC-SSO/internal/storage/sqlite"
)

const StoragePath = "../sso.db"

var Storage, _ = sqlite.New(slog.Default(), StoragePath)

func TestSaveUser(t *testing.T) {

	cases := []struct {
		testName    string
		email       string
		passHash    []byte
		expectedErr error
	}{
		{
			testName:    "ok case",
			email:       "test-mail",
			passHash:    []byte("test-pass"),
			expectedErr: nil,
		},
		{
			testName:    "user exists case",
			email:       "test-mail",
			passHash:    []byte("test-pass"),
			expectedErr: storage.ErrUserExists,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.testName, func(t *testing.T) {
			_, err := Storage.SaveUser(context.Background(),
				tcase.email, tcase.passHash)

			require.EqualValues(t, tcase.expectedErr, err)
		})
	}

}

func TestGetUser(t *testing.T) {
	cases := []struct {
		testName    string
		email       string
		expectedErr error
	}{
		{
			testName:    "ok case",
			email:       "test-mail",
			expectedErr: nil,
		},
		{
			testName:    "user not found case",
			email:       "wrong-test-mail",
			expectedErr: storage.ErrUserNotFound,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.testName, func(t *testing.T) {
			_, err := Storage.GetUser(context.Background(),
				tcase.email)

			require.EqualValues(t, tcase.expectedErr, err)
		})
	}
}

func TestIsAdmin(t *testing.T) {
	cases := []struct {
		testName    string
		userID      string
		expectedErr error
	}{
		{
			testName:    "ok case",
			userID:      "de31e37e-686a-4c7c-92b1-cb9ae5f1b953",
			expectedErr: nil,
		},
		{
			testName:    "user not found case",
			userID:      uuid.Nil.String(),
			expectedErr: storage.ErrUserNotFound,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.testName, func(t *testing.T) {
			_, err := Storage.IsAdmin(context.Background(),
				tcase.userID)

			require.EqualValues(t, tcase.expectedErr, err)
		})
	}
}
