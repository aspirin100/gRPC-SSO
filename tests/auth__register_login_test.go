package auth_test

import (
	"log"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	ssov1 "github.com/aspirin100/gRPC-SSO/protos/gen/go/sso"
	test_client "github.com/aspirin100/gRPC-SSO/tests/test-client"
)

const (
	emptyAppID          = 0
	appID               = 1
	passwordWrongMaxLen = 17
	passwordDefaultLen  = 10

	testSecretKey = "local_secret_key"
)

func TestRegisterLogin(t *testing.T) {
	ctx, testClient := test_client.New(t)

	email := gofakeit.Email()
	password := generatePassword(passwordDefaultLen)

	log.Println(testClient.Cfg.SecretKey)

	registerReponse, err := testClient.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	require.NotEmpty(t, registerReponse)

	loginResponse, err := testClient.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppID:    appID,
	})
	require.NoError(t, err)

	var claims jwt.MapClaims

	_, err = jwt.ParseWithClaims(
		loginResponse.GetAccessToken(),
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(testSecretKey), nil
		})
	require.NoError(t, err)

	require.EqualValues(t, appID, int(claims["appID"].(float64)))                  //nolint:forcetypeassert
	require.EqualValues(t, registerReponse.GetUserID(), claims["userID"].(string)) //nolint:forcetypeassert
}

func generatePassword(length int) string {
	return gofakeit.Password(true, true, true, true, false, length)
}
