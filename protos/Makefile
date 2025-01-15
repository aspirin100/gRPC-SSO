.PHONY: generate
generate:
		protoc -I proto proto/sso/*.proto --go_out=./gen/go/sso \
		--go_opt=paths=source_relative --go-grpc_out=./gen/go/sso --go-grpc_opt=paths=source_relative
