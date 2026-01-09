# Envoy Integration Example

This example demonstrates how to integrate `auth-server` with [Envoy](https://www.envoyproxy.io/) using the [External Authorization](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter) (`ext_authz`) filter.

## Prerequisites

1. Generate RSA keys (see `../../secrets/README.md`)
2. The example includes `service_config.yml` pre-configured with `proxy: envoy`
3. Ensure `../../config/local_repository_config.yml` exists

## Architecture

Envoy uses the `ext_authz` HTTP filter to make an external authorization call to the auth-server's `/auth` endpoint:
- If the auth server returns `200 OK`, the request is forwarded to the backend
- If the auth server returns `401 Unauthorized`, Envoy returns an error to the client
- Envoy sets `X-Envoy-*` and `X-Forwarded-*` headers that the auth-server parses

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

5. Check Envoy admin interface:
```sh
curl http://localhost:9901/stats
```

## Configuration

* Envoy proxy: http://localhost:8082
* Envoy admin: http://localhost:9901
* Auth server: http://localhost:8081

## Envoy Configuration Details

The `envoy.yaml` file includes:

- **External Authorization Filter** (`ext_authz`): Calls auth-server for all requests except bypassed routes
- **Token endpoint** (`/token`): Directly routes to auth-server, bypasses auth
- **Auth endpoint** (`/auth`): Directly routes to auth-server, bypasses auth
- **Health endpoints** (`/health`, `/ready`, `/version`): Bypass auth
- **Protected routes** (`/`): Require authentication via `ext_authz` filter
- **Header forwarding**: Uses Lua filter to set `X-Original-Method` and `X-Original-Uri` headers for auth-server

## Customization

### Protect Specific Routes
To protect only specific paths, modify the route configuration:
```yaml
routes:
  # Public route - no auth
  - match:
      path: "/public"
    route:
      cluster: backend
  # Protected route - requires auth
  - match:
      prefix: "/api"
    route:
      cluster: backend
    # ext_authz filter will handle auth for this route
```

### Adjust Timeouts
Modify timeout values in the configuration:
```yaml
http_service:
  server_uri:
    timeout: 10s  # Increase auth timeout
```

### Add More Headers
To forward additional headers to auth-server:
```yaml
authorization_request:
  allowed_headers:
    patterns:
      - exact: authorization
      - exact: x-custom-header
      - prefix: x-forwarded-
```

## Troubleshooting

### Check Envoy logs
```sh
docker logs envoy
```

### Check auth-server logs
```sh
docker logs auth-server
```

### Test auth endpoint directly
```sh
curl -H "Authorization: Bearer <token>" http://localhost:8081/auth
```

### Check Envoy stats
```sh
curl http://localhost:9901/stats | grep ext_authz
```

### View Envoy configuration
```sh
curl http://localhost:9901/config_dump
```

### Verify Envoy configuration
```sh
docker exec envoy envoy --mode validate -c /etc/envoy/envoy.yaml
```

## Differences from Traefik/Nginx

- **Configuration**: Envoy uses YAML configuration files vs Traefik's dynamic config or Nginx's nginx.conf
- **Filter-based**: Uses HTTP filters (`ext_authz`) vs Traefik's middleware or Nginx's `auth_request`
- **Admin interface**: Envoy provides a built-in admin interface on port 9901
- **Header handling**: Envoy sets both `X-Forwarded-*` and `X-Envoy-*` headers
- **More granular control**: Can configure per-route auth policies

## Advanced Configuration

### Per-Route Auth Policies

You can configure different auth policies for different routes:
```yaml
routes:
  - match:
      prefix: "/admin"
    route:
      cluster: backend
    typed_per_filter_config:
      envoy.filters.http.ext_authz:
        "@type": type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthzPerRoute
        # Custom auth config for this route
```

### Response Headers

To forward auth response headers to backend:
```yaml
allowed_upstream_headers:
  patterns:
    - exact: x-auth-user
    - exact: x-auth-role
```

This allows the auth-server to set headers that Envoy will forward to the backend service.
