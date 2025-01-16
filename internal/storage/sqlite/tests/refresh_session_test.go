package storage_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/aspirin100/gRPC-SSO/internal/storage"
	"github.com/aspirin100/gRPC-SSO/internal/tokens"
	"github.com/stretchr/testify/require"
)

var userID = "48fd672b-81cb-4a37-a33f-8e2b1e0031d5"
var refreshTokenPrepared = "e8072fc7dfc55a232794a3b8d14e16c4e997bc5f490cd9cd57c8d4e8ff194554"

func TestNewRefreshSession(t *testing.T) {
	refreshToken, err := tokens.NewRefreshToken()
	if err != nil {
		log.Print(err)
		t.Fail()
	}

	err = Storage.NewRefreshSession(context.Background(),
		*refreshToken,
		userID, time.Minute*60)
	if err != nil {
		log.Print(err)
		t.Fail()
	}
}

func TestValidateRefreshToken(t *testing.T) {
	cases := []struct {
		testName      string
		refreshToken  string
		userID        string
		expectedError error
	}{
		{
			testName:      "not found case",
			refreshToken:  "",
			userID:        userID,
			expectedError: storage.ErrRefreshTokenNotFound,
		},
		{
			testName:      "invalid refresh token case",
			refreshToken:  refreshTokenPrepared,
			userID:        userID,
			expectedError: tokens.ErrInvalidRefreshToken,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.testName, func(t *testing.T) {
			err := Storage.ValidateRefreshToken(context.Background(),
				tcase.refreshToken, tcase.userID)

			require.EqualValues(t, tcase.expectedError, err)
		})
	}
}
