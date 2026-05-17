// Package primitives mirrors the JSON Schemas in templates/agent-primitives/schemas/.
// Hand-written here for template simplicity; projects are encouraged to codegen
// these types from the JSON Schemas so contract and code cannot drift.
package primitives

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
)

var dottedNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+$`)
var semverPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

type Task struct {
	TaskID          string           `json:"task_id"`
	Kind            string           `json:"kind"`
	Version         string           `json:"version,omitempty"`
	Inputs          json.RawMessage  `json:"inputs"`
	Params          json.RawMessage  `json:"params,omitempty"`
	IdempotencyKey  string           `json:"idempotency_key"`
	Budget          *Budget          `json:"budget,omitempty"`
	CapabilityToken *CapabilityToken `json:"capability_token,omitempty"`
	Provenance      Provenance       `json:"provenance"`
}

// Validate enforces the transport-neutral Task invariants mirrored from
// templates/agent-primitives/schemas/task.schema.json. Kind-specific input
// validation still belongs at the Producer boundary.
func (task *Task) Validate() error {
	if strings.TrimSpace(task.TaskID) == "" {
		return errors.New("task_id is required")
	}
	if !dottedNamePattern.MatchString(task.Kind) {
		return errors.New("kind must be a dotted lowercase name")
	}
	if task.Version != "" && !semverPattern.MatchString(task.Version) {
		return errors.New("version must be semver")
	}
	if !isJSONObject(task.Inputs) {
		return errors.New("inputs must be a JSON object")
	}
	if len(task.Params) > 0 && !isJSONObject(task.Params) {
		return errors.New("params must be a JSON object")
	}
	if strings.TrimSpace(task.IdempotencyKey) == "" {
		return errors.New("idempotency_key is required")
	}
	if strings.TrimSpace(task.Provenance.ProducedBy.AgentID) == "" {
		return errors.New("provenance.produced_by.agent_id is required")
	}
	if task.Provenance.ProducedAt.IsZero() {
		return errors.New("provenance.produced_at is required")
	}
	if strings.TrimSpace(task.Provenance.StepID) == "" {
		return errors.New("provenance.step_id is required")
	}
	return nil
}

type Budget struct {
	MaxTokens      int     `json:"max_tokens,omitempty"`
	MaxWallSeconds int     `json:"max_wall_seconds,omitempty"`
	MaxCostUSD     float64 `json:"max_cost_usd,omitempty"`
}

type Status string

const (
	StatusOK      Status = "ok"
	StatusPartial Status = "partial"
	StatusError   Status = "error"
)

type Result struct {
	ResultID   string          `json:"result_id"`
	TaskID     string          `json:"task_id"`
	Status     Status          `json:"status"`
	Output     json.RawMessage `json:"output,omitempty"`
	Error      *Error          `json:"error,omitempty"`
	Confidence *float64        `json:"confidence,omitempty"`
	Evidence   []Evidence      `json:"evidence,omitempty"`
	Metrics    *Metrics        `json:"metrics,omitempty"`
	Provenance Provenance      `json:"provenance"`
}

type Error struct {
	Code      string          `json:"code"`
	Message   string          `json:"message"`
	Retryable bool            `json:"retryable,omitempty"`
	Details   json.RawMessage `json:"details,omitempty"`
}

type Metrics struct {
	TokensIn    int     `json:"tokens_in,omitempty"`
	TokensOut   int     `json:"tokens_out,omitempty"`
	WallSeconds float64 `json:"wall_seconds,omitempty"`
	CostUSD     float64 `json:"cost_usd,omitempty"`
}

type Evidence struct {
	Claim      string         `json:"claim"`
	Source     EvidenceSource `json:"source"`
	Confidence *float64       `json:"confidence,omitempty"`
}

type EvidenceSource struct {
	Kind  string `json:"kind"`
	ID    string `json:"id"`
	Range string `json:"range,omitempty"`
}

type Provenance struct {
	ProducedBy   Producer  `json:"produced_by"`
	ProducedAt   time.Time `json:"produced_at"`
	StepID       string    `json:"step_id"`
	ParentStepID string    `json:"parent_step_id,omitempty"`
	TraceID      string    `json:"trace_id,omitempty"`
	Inputs       []string  `json:"inputs,omitempty"`
	Ring         int       `json:"ring"`
}

type Producer struct {
	AgentID string `json:"agent_id"`
	Model   string `json:"model,omitempty"`
	Tool    string `json:"tool,omitempty"`
}

type CapabilityToken struct {
	TokenID        string       `json:"token_id"`
	Issuer         string       `json:"issuer"`
	Subject        string       `json:"subject"`
	Capabilities   []Capability `json:"capabilities"`
	ExpiresAt      time.Time    `json:"expires_at"`
	AttenuatedFrom string       `json:"attenuated_from,omitempty"`
	Nonce          string       `json:"nonce,omitempty"`
	Signature      string       `json:"signature,omitempty"`
}

type Capability struct {
	Action string          `json:"action"`
	Scope  json.RawMessage `json:"scope,omitempty"`
}

func isJSONObject(raw json.RawMessage) bool {
	var value map[string]json.RawMessage
	return json.Unmarshal(raw, &value) == nil
}
