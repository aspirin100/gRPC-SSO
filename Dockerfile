FROM golang:alpine3.21 AS build

RUN apk add --no-cache make git build-base

COPY . /app/go/grpc-sso

WORKDIR /app/go/grpc-sso

RUN make docker-build

FROM alpine:3.21.2

COPY --from=build /app/go/grpc-sso/bin/sso-app.out /usr/local/bin/sso-app.out
COPY --from=build /app/go/grpc-sso/config/local.yaml /usr/local/bin/config.yaml
COPY --from=build /app/go/grpc-sso/internal/storage/sqlite/sso.db /usr/local/bin/sso.db

EXPOSE 443

WORKDIR /usr/local/bin

ENTRYPOINT sso-app.out --confpath="." 