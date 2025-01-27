package storage

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/aspirin100/gRPC-SSO/internal/storage"
)

const StoragePath = "../sso.db"

var Storage, _ = storage.New(slog.Default(), StoragePath)

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
			userID:      "4f1eb1b5-0520-4917-b865-3f2d460f9603",
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

func TestGetApp(t *testing.T) {
	cases := []struct {
		testName    string
		appID       int32
		expectedErr error
	}{
		{
			testName:    "ok case",
			appID:       1,
			expectedErr: nil,
		},
		{
			testName:    "app not found case",
			appID:       0,
			expectedErr: storage.ErrAppNotFound,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.testName, func(t *testing.T) {
			_, err := Storage.GetApp(context.Background(),
				tcase.appID)

			require.EqualValues(t, tcase.expectedErr, err)
		})
	}
}
