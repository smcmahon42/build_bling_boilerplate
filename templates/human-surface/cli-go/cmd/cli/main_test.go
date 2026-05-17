package main

import (
	"encoding/base64"
	"testing"
)

func TestParseCapabilityTokenAcceptsJSONObject(t *testing.T) {
	raw := `{"token_id":"token-1"}`

	parsed, err := parseCapabilityToken(raw)
	if err != nil {
		t.Fatalf("expected JSON token to parse, got %v", err)
	}
	if string(parsed) != raw {
		t.Fatalf("expected token to round-trip, got %s", parsed)
	}
}

func TestParseCapabilityTokenAcceptsBase64URLJSON(t *testing.T) {
	raw := `{"token_id":"token-1"}`
	encoded := base64.RawURLEncoding.EncodeToString([]byte(raw))

	parsed, err := parseCapabilityToken(encoded)
	if err != nil {
		t.Fatalf("expected base64url JSON token to parse, got %v", err)
	}
	if string(parsed) != raw {
		t.Fatalf("expected decoded token, got %s", parsed)
	}
}

func TestParseCapabilityTokenRejectsNonObject(t *testing.T) {
	if _, err := parseCapabilityToken(`"token-1"`); err == nil {
		t.Fatal("expected non-object token error")
	}
}
