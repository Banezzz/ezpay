package http_client

import "testing"

func TestValidateOutboundURLBlocksPrivateTargets(t *testing.T) {
	for _, rawURL := range []string{
		"http://127.0.0.1:8080/callback",
		"http://10.0.0.1/callback",
		"http://169.254.169.254/latest/meta-data/",
		"http://[::1]/callback",
	} {
		if err := ValidateOutboundURL(rawURL); err == nil {
			t.Fatalf("ValidateOutboundURL(%q) = nil, want error", rawURL)
		}
	}
}

func TestValidateOutboundURLRejectsUnsafeShapes(t *testing.T) {
	for _, rawURL := range []string{
		"ftp://example.com/callback",
		"//example.com/callback",
		"http://user:pass@example.com/callback",
	} {
		if err := ValidateOutboundURL(rawURL); err == nil {
			t.Fatalf("ValidateOutboundURL(%q) = nil, want error", rawURL)
		}
	}
}

func TestValidateOutboundURLAllowsPrivateTargetsWhenExplicitlyEnabled(t *testing.T) {
	t.Setenv("EPUSDT_ALLOW_PRIVATE_CALLBACKS", "true")
	if err := ValidateOutboundURL("http://127.0.0.1:8080/callback"); err != nil {
		t.Fatalf("ValidateOutboundURL with private callbacks enabled: %v", err)
	}
}
