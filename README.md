# auth-server
[![GoDoc](https://godoc.org/github.com/reugn/auth_server?status.svg)](https://godoc.org/github.com/reugn/auth_server)

Simple authentication and authorization server
* **Authentication** is used by a server when the server needs to know exactly who is accessing their information or site.
* **Authorization** is a process by which a server determines if the client has permission to use a resource or access a file.

Building an authentication and authorization strategy is always a challenging process.
Just a number of quick questions that immediately arise:
* Should we grab separate services for authentication and authorization
* How do we handle token creation and who is responsible for this
* Should we alter our REST service to support auth flows

`auth-server` project tries to accumulate all those capabilities and act as a transparent proxy server auth middleware.

## Architecture
![](./images/architecture_diagram_1.png)

1. Client requests an authentication token (JWT) with a basic authentication
2. Proxy server routs this request to `auth-server` to issue a token
3. Client sends an authenticated request (with bearer token) to the proxy server
4. Proxy invokes `auth-server` as an authentication/authorization middleware. In case the token was successfully authenticated/authorized, the request will be routed to the target REST service. Otherwise an auth error code will be returned to the client.

## Prerequisites
`auth-server` written in Golang.
To install the latest stable version of Go, visit https://golang.org/dl/

To run the project using Docker visit their [page](https://www.docker.com/get-started) to get started.

Install `docker-compose` to get started with the examples.

Read the following [instructions](./secrets/README.md) to generate keys.

## Examples
Examples are available under the examples folder.

To run [Traefik](https://docs.traefik.io/) configuration:
* `cd examples/traefik`
* `docker-compose up -d`

## License
Licensed under the Apache 2.0 License.