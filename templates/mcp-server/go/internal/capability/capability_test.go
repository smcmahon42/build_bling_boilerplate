package capability

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/example/project/mcp-server/internal/primitives"
)

func TestCheckGrantsExactScope(t *testing.T) {
	token := validToken()
	token.Capabilities[0].Scope = json.RawMessage(`{"thread_id":"thread-1"}`)

	denied := Check(StubVerifier{}, token, "agent-a", "invoke.summarize.thread", json.RawMessage(`{"thread_id":"thread-1","audience":"reviewer"}`), time.Now())
	if denied != nil {
		t.Fatalf("expected scoped capability to grant action, got %v", denied)
	}
}

func TestCheckDeniesMismatchedScope(t *testing.T) {
	token := validToken()
	token.Capabilities[0].Scope = json.RawMessage(`{"thread_id":"thread-1"}`)

	denied := Check(StubVerifier{}, token, "agent-a", "invoke.summarize.thread", json.RawMessage(`{"thread_id":"thread-2"}`), time.Now())
	if denied == nil || denied.Code != "capability.denied" {
		t.Fatalf("expected capability.denied, got %v", denied)
	}
}

func TestCheckAllowsUnscopedCapability(t *testing.T) {
	token := validToken()

	denied := Check(StubVerifier{}, token, "agent-a", "invoke.summarize.thread", json.RawMessage(`{"thread_id":"thread-1"}`), time.Now())
	if denied != nil {
		t.Fatalf("expected unscoped capability to grant action, got %v", denied)
	}
}

func validToken() *primitives.CapabilityToken {
	return &primitives.CapabilityToken{
		TokenID:   "token-1",
		Issuer:    "agent-root",
		Subject:   "agent-a",
		ExpiresAt: time.Now().Add(time.Hour),
		Signature: "placeholder-signature",
		Capabilities: []primitives.Capability{
			{Action: "invoke.summarize.thread"},
		},
	}
}
