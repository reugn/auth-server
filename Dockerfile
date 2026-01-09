# syntax=docker/dockerfile:1.4
FROM golang:alpine3.23 AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/auth ./cmd/auth

FROM alpine:3.23
RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -D appuser
WORKDIR /app
COPY --from=build /go/src/app/bin/auth /app
COPY --from=build /go/src/app/config /app/
COPY ./secrets ./secrets
RUN chown -R appuser:appgroup /app
ENV AUTH_SERVER_LOCAL_CONFIG_PATH=local_repository_config.yml
USER appuser

EXPOSE 8081
ENTRYPOINT ["/app/auth", "-c", "service_config.yml"]
