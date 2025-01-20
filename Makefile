GOLANGCI_LINT_VERSION := v1.62.2

.PHONY: run
run:
	go run ./cmd/sso/main.go --confpath="./config/local.yaml"

.PHONY: migrations-up
migrations-up:
	go run ./cmd/migrator/main.go \
	--storage-path ./internal/storage/sqlite/sso.db \
	--migrations-path  ./internal/storage/migrations

.PHONY: docker-build
docker-build:
	mkdir -p bin
	CGO_ENABLED=1  GOOS=linux GOARCH=amd64 go build -o ./bin/sso-app.out ./cmd/sso/main.go 

.PHONY: cover
cover:
		go test -short -race -coverprofile=coverage.out ./... 
		go tool cover -html=coverage.out
		rm coverage.out

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

.PHONY: docker-up
docker-up:
	docker build -t sso .
	docker run --rm -d \
	-e SECRET_KEY="secret_key" \
	-e STORAGE_PATH="sso.db" \
	-e CONFIG_PATH="config.yaml" \
	-p 443:443 sso