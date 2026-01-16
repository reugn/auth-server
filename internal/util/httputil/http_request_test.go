package httputil

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testRemoteAddr = "192.168.1.1:12345"

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		expected   string
	}{
		{
			name:       "X-Forwarded-For single IP",
			headers:    map[string]string{"X-Forwarded-For": "192.168.1.1"},
			remoteAddr: "10.0.0.1:12345",
			expected:   "192.168.1.1",
		},
		{
			name:       "X-Forwarded-For multiple IPs",
			headers:    map[string]string{"X-Forwarded-For": "192.168.1.1, 10.0.0.2, 172.16.0.1"},
			remoteAddr: "10.0.0.1:12345",
			expected:   "192.168.1.1",
		},
		{
			name:       "X-Forwarded-For with spaces",
			headers:    map[string]string{"X-Forwarded-For": "  192.168.1.100  "},
			remoteAddr: "10.0.0.1:12345",
			expected:   "192.168.1.100",
		},
		{
			name:       "X-Real-IP header",
			headers:    map[string]string{"X-Real-IP": "172.16.0.50"},
			remoteAddr: "10.0.0.1:12345",
			expected:   "172.16.0.50",
		},
		{
			name:       "X-Forwarded-For takes precedence over X-Real-IP",
			headers:    map[string]string{"X-Forwarded-For": "192.168.1.1", "X-Real-IP": "172.16.0.50"},
			remoteAddr: "10.0.0.1:12345",
			expected:   "192.168.1.1",
		},
		{
			name:       "RemoteAddr fallback with port",
			headers:    map[string]string{},
			remoteAddr: "203.0.113.50:54321",
			expected:   "203.0.113.50",
		},
		{
			name:       "RemoteAddr fallback without port",
			headers:    map[string]string{},
			remoteAddr: "203.0.113.50",
			expected:   "203.0.113.50",
		},
		{
			name:       "IPv6 address in X-Forwarded-For",
			headers:    map[string]string{"X-Forwarded-For": "2001:db8::1"},
			remoteAddr: "10.0.0.1:12345",
			expected:   "2001:db8::1",
		},
		{
			name:       "Empty X-Forwarded-For falls back to X-Real-IP",
			headers:    map[string]string{"X-Forwarded-For": "", "X-Real-IP": "172.16.0.50"},
			remoteAddr: "10.0.0.1:12345",
			expected:   "172.16.0.50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.RemoteAddr = tt.remoteAddr
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			result := GetClientIP(req)
			if result != tt.expected {
				t.Errorf("GetClientIP() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestLogAttrs(t *testing.T) {
	t.Run("basic attributes", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = testRemoteAddr
		req.Header.Set("User-Agent", "TestAgent/1.0")

		attrs := LogAttrs(req)

		if len(attrs) != 2 {
			t.Errorf("LogAttrs() returned %d attributes, want 2", len(attrs))
		}

		// Check IP attribute
		ipAttr, ok := attrs[0].(slog.Attr)
		if !ok {
			t.Fatal("first attribute is not slog.Attr")
		}
		if ipAttr.Key != "ip" || ipAttr.Value.String() != "192.168.1.1" {
			t.Errorf("IP attribute = %v, want ip=192.168.1.1", ipAttr)
		}

		// Check user_agent attribute
		uaAttr, ok := attrs[1].(slog.Attr)
		if !ok {
			t.Fatal("second attribute is not slog.Attr")
		}
		if uaAttr.Key != "user_agent" || uaAttr.Value.String() != "TestAgent/1.0" {
			t.Errorf("user_agent attribute = %v, want user_agent=TestAgent/1.0", uaAttr)
		}
	})

	t.Run("with additional attributes", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = testRemoteAddr
		req.Header.Set("User-Agent", "TestAgent/1.0")

		extraAttrs := []slog.Attr{
			slog.String("username", "testuser"),
			slog.Int("status", 200),
		}

		attrs := LogAttrs(req, extraAttrs...)

		if len(attrs) != 4 {
			t.Errorf("LogAttrs() returned %d attributes, want 4", len(attrs))
		}

		// Check additional attributes
		usernameAttr, ok := attrs[2].(slog.Attr)
		if !ok {
			t.Fatal("third attribute is not slog.Attr")
		}
		if usernameAttr.Key != "username" || usernameAttr.Value.String() != "testuser" {
			t.Errorf("username attribute = %v, want username=testuser", usernameAttr)
		}

		statusAttr, ok := attrs[3].(slog.Attr)
		if !ok {
			t.Fatal("fourth attribute is not slog.Attr")
		}
		if statusAttr.Key != "status" || statusAttr.Value.Int64() != 200 {
			t.Errorf("status attribute = %v, want status=200", statusAttr)
		}
	})

	t.Run("empty user agent", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = testRemoteAddr

		attrs := LogAttrs(req)

		if len(attrs) != 2 {
			t.Errorf("LogAttrs() returned %d attributes, want 2", len(attrs))
		}

		uaAttr, ok := attrs[1].(slog.Attr)
		if !ok {
			t.Fatal("second attribute is not slog.Attr")
		}
		if uaAttr.Key != "user_agent" || uaAttr.Value.String() != "" {
			t.Errorf("user_agent attribute = %v, want user_agent=\"\"", uaAttr)
		}
	})
}
