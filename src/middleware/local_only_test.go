package middleware

import (
	"net/http"
	"testing"
)

func TestIsLoopbackRequest(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "127.0.0.1:12345",
		Host:       "localhost:8000",
		Header:     http.Header{},
	}
	if !IsLoopbackRequest(req) {
		t.Fatal("loopback request rejected")
	}
}

func TestIsLoopbackRequestRejectsPublicHost(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "127.0.0.1:12345",
		Host:       "pay.example.com",
		Header:     http.Header{},
	}
	if IsLoopbackRequest(req) {
		t.Fatal("request with public host accepted")
	}
}

func TestIsLoopbackRequestRejectsForwardedPublicClient(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "127.0.0.1:12345",
		Host:       "localhost:8000",
		Header: http.Header{
			"X-Forwarded-For": []string{"203.0.113.10"},
		},
	}
	if IsLoopbackRequest(req) {
		t.Fatal("request with public forwarded client accepted")
	}
}
