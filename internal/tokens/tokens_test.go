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
	}

	app := entity.App{
		ID: 1,
	}

	token, err := tokens.NewAccessToken(user.UserID, app.ID, time.Minute*15, []byte("secret_test_key"))
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
