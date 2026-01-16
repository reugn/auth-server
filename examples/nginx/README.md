# Nginx Integration Example

This example demonstrates how to integrate `auth-server` with [Nginx](https://nginx.org/) using the [auth_request](https://nginx.org/en/docs/http/ngx_http_auth_request_module.html) module.

## Prerequisites

1. Generate RSA keys (see `../../secrets/README.md`)
2. The example includes `service_config.yml` pre-configured with `proxy: nginx`
3. Ensure `../../config/local_repository_config.yml` exists

## Architecture

Nginx uses the `auth_request` module to make a subrequest to the auth-server's `/auth` endpoint:
- If the auth server returns `200 OK`, the request is forwarded to the backend
- If the auth server returns `401 Unauthorized`, Nginx returns an error to the client
- Nginx sets `X-Forwarded-*` headers that the auth-server parses

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

* Nginx proxy: http://localhost:8082
* Auth server: http://localhost:8081
* Backend service: Internal (example nginx container)

## Nginx Configuration Details

The `nginx.conf` file includes:

- **Token endpoint** (`/token`): Directly proxies to auth-server, no authentication required
- **Auth endpoint** (`/auth`): Used internally by `auth_request`, no authentication required
- **Protected routes** (`/`): Uses `auth_request` to validate tokens before proxying to backend
- **Headers**: Nginx sets `X-Forwarded-*` headers for the auth-server to parse

## Customization

### Protect Specific Path
```nginx
location /api {
    auth_request /auth;
    proxy_pass http://backend;
    # ... other proxy settings
}
```

### Protect Multiple Paths
```nginx
location ~ ^/(api|admin|dashboard) {
    auth_request /auth;
    proxy_pass http://backend;
    # ... other proxy settings
}
```

## Troubleshooting

### Check Nginx logs
```sh
docker logs nginx
```

### Check auth-server logs
```sh
docker logs auth-server
```

### Test auth endpoint directly
```sh
curl -H "Authorization: Bearer <token>" http://localhost:8081/auth
```

### Verify Nginx configuration
```sh
docker exec nginx nginx -t
```
