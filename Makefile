GOLANGCI_LINT_VERSION := v1.62.2

.PHONY: run
run:
	go run ./cmd/sso/main.go --confpath="./config/local.yaml"

.PHONY: migrations-up
migrations-up:
	go run ./cmd/migrator/main.go \
	--storage-path ./internal/storage/sqlite/sso.db \
	--migrations-path  ./internal/storage/migrations

.PHONY: cover
cover:
		go test -short -race -coverprofile=coverage.out ./... 
		go tool cover -html=coverage.out
		rm coverage.out

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run