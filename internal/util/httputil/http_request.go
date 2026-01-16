package httputil

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
)

// GetClientIP extracts the client IP address from the request.
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxy scenarios)
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if ip := strings.Split(forwarded, ",")[0]; ip != "" {
			return strings.TrimSpace(ip)
		}
	}

	// Check X-Real-IP header (alternative proxy header)
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// LogAttrs returns log attributes for request logging.
// It includes IP address and user agent. Additional attributes can be provided.
func LogAttrs(r *http.Request, attrs ...slog.Attr) []any {
	attrsList := make([]any, 0, len(attrs)+2)
	attrsList = append(attrsList,
		slog.String("ip", GetClientIP(r)),
		slog.String("user_agent", r.UserAgent()))

	for _, attr := range attrs {
		attrsList = append(attrsList, attr)
	}

	return attrsList
}
