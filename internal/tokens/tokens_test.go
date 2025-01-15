package tokens_test

import (
	"log"
	"testing"
	"time"

	"github.com/aspirin100/gRPC-SSO/internal/entity"
	"github.com/aspirin100/gRPC-SSO/internal/tokens"
)

func TestNewAccessToken(t *testing.T) {
	user := entity.User{
		UserID: "some test user id",
		Email:  "test@example.com",
	}

	app := entity.App{
		ID:        1,
		Name:      "test-app-name",
		SecretKey: "test secret key",
	}

	token, err := tokens.NewAccessToken(user, app, time.Minute*15)
	if err != nil {
		log.Print(err)
		t.Fail()
	}

	log.Print(*token)
}

func TestNewRefreshToken(t *testing.T) {
	token, err := tokens.NewRefreshToken()
	if err != nil {
		log.Print(err)
		t.Fail()
	}

	log.Print(len(*token))
}
