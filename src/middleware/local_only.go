package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// LocalOnly allows only direct loopback requests. It checks RemoteAddr, Host,
// and common forwarding headers so a public reverse proxy on the same machine
// does not accidentally expose setup-only endpoints.
func LocalOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if !IsLoopbackRequest(ctx.Request()) {
				return echo.NewHTTPError(http.StatusForbidden, "local access required")
			}
			return next(ctx)
		}
	}
}

func IsLoopbackRequest(req *http.Request) bool {
	return IsLoopbackRemoteAddr(req.RemoteAddr) &&
		IsLoopbackHost(req.Host) &&
		ForwardedHeadersAreLoopback(req.Header)
}

func IsLoopbackRemoteAddr(remoteAddr string) bool {
	return IsLoopbackHost(remoteAddr)
}

func IsLoopbackHost(remoteAddr string) bool {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}
	host = strings.Trim(strings.ToLower(strings.TrimSpace(host)), "[]")
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func ForwardedHeadersAreLoopback(header http.Header) bool {
	for _, value := range header.Values("X-Forwarded-For") {
		for _, part := range strings.Split(value, ",") {
			if !IsLoopbackHost(part) {
				return false
			}
		}
	}
	for _, value := range header.Values("X-Real-IP") {
		if !IsLoopbackHost(value) {
			return false
		}
	}
	for _, value := range header.Values("Forwarded") {
		if !forwardedHeaderIsLoopback(value) {
			return false
		}
	}
	return true
}

func forwardedHeaderIsLoopback(value string) bool {
	for _, entry := range strings.Split(value, ",") {
		for _, part := range strings.Split(entry, ";") {
			key, val, ok := strings.Cut(strings.TrimSpace(part), "=")
			if !ok || !strings.EqualFold(strings.TrimSpace(key), "for") {
				continue
			}
			val = strings.Trim(strings.TrimSpace(val), `"`)
			if !IsLoopbackHost(val) {
				return false
			}
		}
	}
	return true
}
