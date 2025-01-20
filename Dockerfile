FROM golang:1.24rc2-alpine3.21 AS build

RUN apk --no-cache add make git

COPY . /app/go/grpc-sso

WORKDIR /app/go/grpc-sso

RUN go get github.com/mattn/go-sqlite3

RUN make build


FROM alpine:3.21.2

COPY --from=build /app/go/grpc-sso/bin/grpc-sso /usr/local/bin/grpc-sso
COPY --from=build /app/go/grpc-sso/config/local.yaml /usr/local/bin/config.yaml

ENTRYPOINT ./grpc-sso-app --confpath="."