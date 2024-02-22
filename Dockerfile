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
ENV AUTH_SERVER_PRIVATE_KEY_PATH=secrets/privkey.pem
ENV AUTH_SERVER_PUBLIC_KEY_PATH=secrets/cert.pem

EXPOSE 8080
ENTRYPOINT ["/app/auth", "-c", "service_config.yml"]
