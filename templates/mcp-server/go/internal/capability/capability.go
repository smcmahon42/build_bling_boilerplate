// Package capability verifies CapabilityTokens before an action is permitted.
// The Verifier interface lets projects plug in real signature crypto
// (Ed25519 with a KMS-backed key, etc.) without changing callers.
package capability

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/example/project/mcp-server/internal/primitives"
)

// Verifier authenticates a CapabilityToken. Implementations check the
// signature against a trusted issuer key, nothing else — the remaining
// checks (expiry, subject, action, scope) live in Check.
type Verifier interface {
	VerifySignature(token *primitives.CapabilityToken) error
}

// Check enforces the token's declared constraints for a specific action.
// Returns a typed error whose Code is suitable for Result.error.code.
func Check(
	v Verifier,
	token *primitives.CapabilityToken,
	selfAgentID string,
	action string,
	requestedScope json.RawMessage,
	now time.Time,
) *Denied {
	if token == nil {
		return &Denied{Code: "capability.missing", Message: "no capability token presented"}
	}
	if len(token.Capabilities) == 0 {
		return &Denied{Code: "capability.empty", Message: "token grants no capabilities"}
	}
	if token.Subject != selfAgentID {
		return &Denied{Code: "capability.wrong_subject", Message: "token not issued to this agent"}
	}
	if !token.ExpiresAt.IsZero() && !now.Before(token.ExpiresAt) {
		return &Denied{Code: "capability.expired", Message: "token is expired"}
	}
	if err := v.VerifySignature(token); err != nil {
		return &Denied{Code: "capability.bad_signature", Message: err.Error()}
	}
	if !grants(token, action, requestedScope) {
		return &Denied{Code: "capability.denied", Message: "token does not grant the requested action"}
	}
	return nil
}

// grants reports whether any capability in the token matches the requested
// action and scope. Empty capability scopes grant the action broadly; non-empty
// scopes must be a subset of the caller-provided requested scope after JSON
// canonicalization. Projects with hierarchical resources should replace this
// with a kind-specific matcher.
func grants(token *primitives.CapabilityToken, action string, requestedScope json.RawMessage) bool {
	for _, capability := range token.Capabilities {
		if capability.Action != action {
			continue
		}
		if len(capability.Scope) == 0 {
			return true
		}
		if scopeMatches(capability.Scope, requestedScope) {
			return true
		}
	}
	return false
}

func scopeMatches(grantedScope, requestedScope json.RawMessage) bool {
	if len(requestedScope) == 0 {
		return false
	}
	var granted map[string]json.RawMessage
	if err := json.Unmarshal(grantedScope, &granted); err != nil {
		return false
	}
	var requested map[string]json.RawMessage
	if err := json.Unmarshal(requestedScope, &requested); err != nil {
		return false
	}
	for key, grantedValue := range granted {
		requestedValue, ok := requested[key]
		if !ok {
			return false
		}
		if !jsonEqual(grantedValue, requestedValue) {
			return false
		}
	}
	return true
}

func jsonEqual(left, right json.RawMessage) bool {
	var leftValue any
	if err := json.Unmarshal(left, &leftValue); err != nil {
		return false
	}
	var rightValue any
	if err := json.Unmarshal(right, &rightValue); err != nil {
		return false
	}
	leftCanonical, err := json.Marshal(leftValue)
	if err != nil {
		return false
	}
	rightCanonical, err := json.Marshal(rightValue)
	if err != nil {
		return false
	}
	return bytes.Equal(leftCanonical, rightCanonical)
}

// Denied is the rejection reason shaped to map cleanly into Result.error.
type Denied struct {
	Code    string
	Message string
}

func (d *Denied) Error() string { return d.Code + ": " + d.Message }

// StubVerifier accepts any non-empty signature. Replace with a real verifier
// (e.g., Ed25519 with a trust root) before using this skeleton in production.
type StubVerifier struct{}

func (StubVerifier) VerifySignature(token *primitives.CapabilityToken) error {
	if token.Signature == "" {
		return errors.New("signature missing")
	}
	// TODO: verify token.Signature against the issuer's public key using
	// a canonical serialization of the token (all fields except Signature).
	return nil
}
