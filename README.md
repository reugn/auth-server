# auth-server
[![Build](https://github.com/reugn/auth-server/actions/workflows/build.yml/badge.svg)](https://github.com/reugn/auth-server/actions/workflows/build.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/reugn/auth-server)](https://pkg.go.dev/github.com/reugn/auth-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/reugn/auth-server)](https://goreportcard.com/report/github.com/reugn/auth-server)

This project offers a toolkit for building and configuring a tailored authentication and authorization service.

`auth-server` can act as a proxy middleware or be configured in a stand-alone mode. It doesn't require any third-party software integration.
Leverage existing backend [storage repositories](internal/repository) for storing security policies or develop a custom one to suit your specific requirements.
For information on how to configure repositories using environment variables, refer to the [repository configuration](docs/repository_configuration.md) page.

> [!NOTE] 
> This project's security has not been thoroughly evaluated. Proceed with caution when setting up your own auth provider.

## Introduction
* **Authentication** is used by a server when the server needs to know exactly who is accessing their information or site.
* **Authorization** is a process by which a server determines if the client has permission to use a resource or access a file.

The inherent complexity of crafting an authentication and authorization strategy raises a barrage of immediate questions:

* Would it be beneficial to utilize separate services for authentication and authorization purposes?
* What is the process for creating access tokens, and who is tasked with this responsibility?
* Is it necessary to adapt our REST service to support an authorization flow?

The `auth-server` project aims to address these concerns by serving as a transparent authentication and authorization proxy middleware.

## Architecture
![architecture_diagram](docs/images/architecture_diagram_1.png)

1. The user requests an access token (JWT), using a basic authentication header:
    ```
    GET /token HTTP/1.1
    Host: localhost:8081
    Authorization: Basic YWRtaW46MTIzNA==
    ```

2. The proxy server routes this request to `auth-server` to issue a token.  
    Response body:  
    `{"access_token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODg5MzMyNTIsImlhdCI6MTU4ODkyOTY1MiwidXNlciI6ImFkbWluIiwicm9sZSI6MX0.LUx9EYsfBZGwbEsofBTT_5Lo3Y_3lk7T8pWLv3bw-XKVOqb_GhaRkVE90QR_sI-bWTkYCFIG9cPYmMXzmPLyjbofgsqTOzH6OaXi3IqxwZRtRGFtuqMoqXkakX5n38mvI3XkIOwFkNosHrpMtIq-HdqB3tfiDJc3YMsYfPbqyRBnBYJu2K51NslGQSiqKSnS_4KeLeaqqdpC7Zdb9Fo-r7EMn3FFuyPEab1iBsrcUYG3qnsKkvDhaq_jEGHflao7dEPEWaiGvJywXWaKR6XyyGtVx0H-OPfgvh1vUCLUUci2K3xE-IxjfRrHx3dSzdqFgJq_n4bVXpO9iNVYOZLccQ","token_type":"Bearer","expires_in":3600000}`

3. The user sends an authenticated request to the proxy server:
    ```
    GET /foo HTTP/1.1
    Host: localhost:8081
    Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODg5MzMyNTIsImlhdCI6MTU4ODkyOTY1MiwidXNlciI6ImFkbWluIiwicm9sZSI6MX0.LUx9EYsfBZGwbEsofBTT_5Lo3Y_3lk7T8pWLv3bw-XKVOqb_GhaRkVE90QR_sI-bWTkYCFIG9cPYmMXzmPLyjbofgsqTOzH6OaXi3IqxwZRtRGFtuqMoqXkakX5n38mvI3XkIOwFkNosHrpMtIq-HdqB3tfiDJc3YMsYfPbqyRBnBYJu2K51NslGQSiqKSnS_4KeLeaqqdpC7Zdb9Fo-r7EMn3FFuyPEab1iBsrcUYG3qnsKkvDhaq_jEGHflao7dEPEWaiGvJywXWaKR6XyyGtVx0H-OPfgvh1vUCLUUci2K3xE-IxjfRrHx3dSzdqFgJq_n4bVXpO9iNVYOZLccQ
    ```

4. Proxy invokes `auth-server` as an authentication/authorization middleware. In case the token was successfully authenticated/authorized, the request will be routed to the target service. Otherwise, an auth error code will be returned to the client.

## Installation and Prerequisites
* `auth-server` is written in Golang.
To install the latest stable version of Go, visit the [releases page](https://golang.org/dl/).

* Read the following [instructions](./secrets/README.md) to generate keys required to sign the token. Specify the location of the generated certificates in the service configuration file. An example of the configuration file can be found [here](config/service_config.yml).

* The following example shows how to run the service using a configuration file:
    ```
    ./auth -c service_config.yml
    ```

* To run the project using Docker, visit their [page](https://www.docker.com/get-started) to get started. Docker images are available under the [GitHub Packages](https://github.com/reugn/auth-server/packages).

* Install `docker-compose` to get started with the examples.

## Examples
Examples are available under the [examples](examples) folder.

To run `auth-server` as a [Traefik](https://docs.traefik.io/) middleware:
```
cd examples/traefik
docker-compose up -d
```

## License
Licensed under the Apache 2.0 License.
