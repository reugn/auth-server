# Traefik Integration Example

This example demonstrates how to integrate `auth-server` with [Traefik](https://traefik.io/) using the [ForwardAuth](https://doc.traefik.io/traefik/middlewares/http/forwardauth/) middleware.

## Prerequisites

1. Generate RSA keys (see `../../secrets/README.md`)
2. The example includes `service_config.yml` pre-configured with `proxy: traefik`
3. Ensure `../../config/local_repository_config.yml` exists

## Architecture

Traefik uses the ForwardAuth middleware to make a subrequest to the auth-server's `/auth` endpoint:
- If the auth server returns `200 OK`, the request is forwarded to the backend
- If the auth server returns `401 Unauthorized`, Traefik returns an error to the client
- Traefik sets `X-Forwarded-*` headers that the auth-server parses

## Usage

1. Start services:
```sh
docker compose up -d
```

2. Get a token:
```sh
curl -u admin:1234 http://localhost:8082/token
```

3. Use token to access protected routes:
```sh
curl -H "Authorization: Bearer <token>" http://localhost:8082/health
```

4. Access without token (will fail):
```sh
curl http://localhost:8082/health
```

## Configuration

* Traefik dashboard: http://localhost:8080
* Auth server: http://localhost:8081
* Public endpoint: http://localhost:8082

## Traefik Configuration Details

The ForwardAuth middleware:
- Sends requests to `http://auth-server:8081/auth`
- Forwards `X-Forwarded-Method`, `X-Forwarded-Uri`, and `X-Forwarded-Prefix` headers
- Can optionally forward auth response headers to the backend

## Customization

### Protect Specific Routes
Add the `auth-server` middleware only to routes that need protection in your dynamic config.

### Forward User Headers
To forward authenticated user info to the backend:
```yaml
forwardAuth:
  address: "http://auth-server:8081/auth"
  authResponseHeaders:
    - "X-Auth-User"
    - "X-Auth-Role"
```

## Troubleshooting

### Check Traefik logs
```sh
docker logs traefik
```

### Check auth-server logs
```sh
docker logs auth-server
```

### Test auth endpoint directly
```sh
curl -H "Authorization: Bearer <token>" http://localhost:8081/auth
```

### Access Traefik dashboard
Open http://localhost:8080 in your browser.
