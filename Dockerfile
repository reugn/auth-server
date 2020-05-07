FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go get ./...
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/auth

FROM alpine:3.9
WORKDIR /go/bin
COPY --from=build /go/src/app/bin /go/bin
COPY ./secrets ./secrets
EXPOSE 8081
ENTRYPOINT ["/go/bin/auth"]