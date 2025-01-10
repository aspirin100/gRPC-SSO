package main

import "github.com/aspirin100/gRPC-SSO/sso/internal/config"

func main() {
	config := config.MustLoad()
	_ = config

}
