.PHONY: run
run:
	go run ./cmd/sso/main.go --confpath="./config/local.yaml"

.PHONY: migrations-up
migrations-up:
	go run ./cmd/migrator/main.go \
	--storage-path ./internal/storage/sqlite/sso.db \
	--migrations-path  ./internal/storage/migrations