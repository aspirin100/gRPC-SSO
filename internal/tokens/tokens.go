package tokens

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidRefreshToken = errors.New("refresh token is expired")
)

const (
	RefreshTokenBytesLen = 32
)

func NewAccessToken(userID string,
	appID int32,
	ttl time.Duration,
	secretKey string) (
	*string, error) {
	claims := jwt.MapClaims{
		"appID":     appID,
		"userID":    userID,
		"expiresAt": time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("jwt token signing failure: %w", err)
	}

	return &signed, nil
}

func NewRefreshToken() (*string, error) {
	randBytes := make([]byte, RefreshTokenBytesLen)

	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src) //nolint:gosec

	_, err := r.Read(randBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	token := hex.EncodeToString(randBytes)

	return &token, nil
}