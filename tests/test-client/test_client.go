package test_client

import (
	"context"
	"testing"

	"github.com/aspirin100/gRPC-SSO/internal/config"
	ssov1 "github.com/aspirin100/gRPC-SSO/protos/gen/go/sso"
)

type TestClient struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient 
}

func New(t *testing.T) (context.Context, *TestClient){
	t.Helper()
	t.Parallel()
}