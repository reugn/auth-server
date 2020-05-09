# auth-server
[![GoDoc](https://godoc.org/github.com/reugn/auth_server?status.svg)](https://godoc.org/github.com/reugn/auth_server)

Simple authentication and authorization server.  
`auth-server` can act as a proxy middleware or be configured in a stand-alone mode as well. It doesn't require any third-party software integration. Use one of the available repositories to configure backend storage, or implement one of your own.

## Introduction
* **Authentication** is used by a server when the server needs to know exactly who is accessing their information or site.
* **Authorization** is a process by which a server determines if the client has permission to use a resource or access a file.

Building an authentication and authorization strategy is always a challenging process.
Just a number of quick questions that immediately arise:
* Should we set up separate services for authentication and authorization
* How do we handle access token creation and who is responsible for this
* Should we alter our REST service to support authorization flow

`auth-server` project tries to accumulate all those capabilities and act as a transparent authentication and authorization proxy middleware.

## Architecture
![](./images/architecture_diagram_1.png)

1. Client requests an access token (JWT), using a basic authentication header:
    ```
    GET /token HTTP/1.1
    Host: localhost:8081
    Authorization: Basic YWRtaW46MTIzNA==
    ```

2. The proxy server routes this request to `auth-server` to issue a token.  
    Response body:  
    `{"access_token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODg5MzMyNTIsImlhdCI6MTU4ODkyOTY1MiwidXNlciI6ImFkbWluIiwicm9sZSI6MX0.LUx9EYsfBZGwbEsofBTT_5Lo3Y_3lk7T8pWLv3bw-XKVOqb_GhaRkVE90QR_sI-bWTkYCFIG9cPYmMXzmPLyjbofgsqTOzH6OaXi3IqxwZRtRGFtuqMoqXkakX5n38mvI3XkIOwFkNosHrpMtIq-HdqB3tfiDJc3YMsYfPbqyRBnBYJu2K51NslGQSiqKSnS_4KeLeaqqdpC7Zdb9Fo-r7EMn3FFuyPEab1iBsrcUYG3qnsKkvDhaq_jEGHflao7dEPEWaiGvJywXWaKR6XyyGtVx0H-OPfgvh1vUCLUUci2K3xE-IxjfRrHx3dSzdqFgJq_n4bVXpO9iNVYOZLccQ","token_type":"Bearer","expires_in":3600000}`

3. Client sends an authenticated request to the proxy server:
    ```
    GET /foo HTTP/1.1
    Host: localhost:8081
    Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODg5MzMyNTIsImlhdCI6MTU4ODkyOTY1MiwidXNlciI6ImFkbWluIiwicm9sZSI6MX0.LUx9EYsfBZGwbEsofBTT_5Lo3Y_3lk7T8pWLv3bw-XKVOqb_GhaRkVE90QR_sI-bWTkYCFIG9cPYmMXzmPLyjbofgsqTOzH6OaXi3IqxwZRtRGFtuqMoqXkakX5n38mvI3XkIOwFkNosHrpMtIq-HdqB3tfiDJc3YMsYfPbqyRBnBYJu2K51NslGQSiqKSnS_4KeLeaqqdpC7Zdb9Fo-r7EMn3FFuyPEab1iBsrcUYG3qnsKkvDhaq_jEGHflao7dEPEWaiGvJywXWaKR6XyyGtVx0H-OPfgvh1vUCLUUci2K3xE-IxjfRrHx3dSzdqFgJq_n4bVXpO9iNVYOZLccQ
    ```

4. Proxy invokes `auth-server` as an authentication/authorization middleware. In case the token was successfully authenticated/authorized, the request will be routed to the target service. Otherwise, an auth error code will be returned to the client.

## Prerequisites
* `auth-server` written in Golang.
To install the latest stable version of Go, visit https://golang.org/dl/

* To run the project using Docker, visit their [page](https://www.docker.com/get-started) to get started. Docker images are available under the [GitHub Packages](https://github.com/reugn/auth-server/packages).

* Install `docker-compose` to get started with the examples.

* Read the following [instructions](./secrets/README.md) to generate keys.

## Examples
Examples are available under the examples folder.

To run `auth-server` as a [Traefik](https://docs.traefik.io/) middleware:
* `cd examples/traefik`
* `docker-compose up -d`

## Contributing
Contributions are very welcome!

## License
Licensed under the Apache 2.0 License.