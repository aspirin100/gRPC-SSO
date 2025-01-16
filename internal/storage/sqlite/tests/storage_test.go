package storage_test

import (
	"context"
	"log/slog"
	"testing"

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
			testName:    "user exists case",
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
