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

func TestNewRefreshSession(t *testing.T) {
	refreshToken, err := tokens.NewRefreshToken()
	if err != nil {
		log.Print(err)
		t.Fail()
	}

	userID := "855ef14d-1bf3-4156-b43a-36eda2493933"

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
			userID:        "855ef14d-1bf3-4156-b43a-36eda2493933",
			expectedError: storage.ErrRefreshTokenNotFound,
		},
		{
			testName:      "invalid refresh token case",
			refreshToken:  "62ff97621d7a0d0178f5204ed80b0bf9750dd0982e6da3033f4c7f4cc08be4f1",
			userID:        "855ef14d-1bf3-4156-b43a-36eda2493933",
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
