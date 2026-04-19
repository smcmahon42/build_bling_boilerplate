// Package handler is the transport-agnostic core: it takes a parsed tool
// call, runs the Task/Result pipeline, and returns a typed response.
// The MCP SDK binding in cmd/server/main.go calls Dispatch.
package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/example/project/mcp-server/internal/capability"
	"github.com/example/project/mcp-server/internal/idempotency"
	"github.com/example/project/mcp-server/internal/primitives"
	"github.com/example/project/mcp-server/internal/producer"
)

type Handler struct {
	AgentID   string
	Producers *producer.Registry
	Capability capability.Verifier
	Idempotency idempotency.Cache
	Now        func() time.Time
}

// ToolCall is the SDK-agnostic input: a tool name (= task kind) and its
// arguments as a JSON blob. The blob is expected to deserialize into a
// full primitives.Task (callers construct the Task from their MCP inputs).
type ToolCall struct {
	Tool      string
	Arguments json.RawMessage
}

// Dispatch is the single entry point. It validates inputs, checks
// capability, honors idempotency, invokes the producer, and returns a
// fully-formed Result.
func (h *Handler) Dispatch(ctx context.Context, call ToolCall) *primitives.Result {
	now := h.Now()

	var task primitives.Task
	if err := json.Unmarshal(call.Arguments, &task); err != nil {
		return h.fail(now, "", "input.invalid", "arguments are not a valid Task: "+err.Error(), false)
	}
	if task.Kind != call.Tool {
		return h.fail(now, task.TaskID, "input.invalid", "tool name does not match Task.kind", false)
	}

	if denied := capability.Check(h.Capability, task.CapabilityToken, h.AgentID, "invoke."+task.Kind, now); denied != nil {
		return h.fail(now, task.TaskID, denied.Code, denied.Message, false)
	}

	if cached, ok := h.Idempotency.Get(task.IdempotencyKey); ok {
		return cached
	}

	prod, ok := h.Producers.Lookup(task.Kind)
	if !ok {
		return h.fail(now, task.TaskID, "kind.unknown", "no producer registered for kind "+task.Kind, false)
	}

	output, perr := prod.Produce(ctx, &task)
	if perr != nil {
		return h.fail(now, task.TaskID, perr.Code, perr.Message, perr.Retryable)
	}

	stepID := newStepID(task.TaskID, now)
	result := &primitives.Result{
		ResultID: resultID(task.TaskID, output, stepID),
		TaskID:   task.TaskID,
		Status:   primitives.StatusOK,
		Output:   output,
		Provenance: primitives.Provenance{
			ProducedBy: primitives.Producer{AgentID: h.AgentID},
			ProducedAt: now,
			StepID:     stepID,
			ParentStepID: task.Provenance.StepID,
			TraceID:    task.Provenance.TraceID,
			Ring:       task.Provenance.Ring + 1,
		},
	}
	h.Idempotency.Put(task.IdempotencyKey, result)
	return result
}

func (h *Handler) fail(now time.Time, taskID, code, message string, retryable bool) *primitives.Result {
	stepID := newStepID(taskID, now)
	return &primitives.Result{
		ResultID: resultID(taskID, nil, stepID),
		TaskID:   taskID,
		Status:   primitives.StatusError,
		Error: &primitives.Error{
			Code:      code,
			Message:   message,
			Retryable: retryable,
		},
		Provenance: primitives.Provenance{
			ProducedBy: primitives.Producer{AgentID: h.AgentID},
			ProducedAt: now,
			StepID:     stepID,
		},
	}
}

func newStepID(taskID string, now time.Time) string {
	h := sha256.New()
	h.Write([]byte(taskID))
	h.Write([]byte(now.Format(time.RFC3339Nano)))
	return "step:" + hex.EncodeToString(h.Sum(nil))[:16]
}

func resultID(taskID string, output json.RawMessage, stepID string) string {
	h := sha256.New()
	h.Write([]byte(taskID))
	h.Write(output)
	h.Write([]byte(stepID))
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}
