# syntax=docker/dockerfile:1.2
FROM golang:alpine3.19 AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go get ./...
WORKDIR /go/src/app/cmd/auth
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/auth

FROM alpine:3.19.1
WORKDIR /app
COPY --from=build /go/src/app/cmd/auth/bin /app
COPY --from=build /go/src/app/config /app/
COPY ./secrets ./secrets
ENV AUTH_SERVER_LOCAL_CONFIG_PATH=local_repository_config.yml

EXPOSE 8081
ENTRYPOINT ["/app/auth", "-c", "service_config.yml"]
