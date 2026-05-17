package primitives

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTaskValidateAcceptsMinimalValidTask(t *testing.T) {
	task := validTask()

	if err := task.Validate(); err != nil {
		t.Fatalf("expected valid task, got %v", err)
	}
}

func TestTaskValidateRejectsInvalidKind(t *testing.T) {
	task := validTask()
	task.Kind = "summarize"

	if err := task.Validate(); err == nil {
		t.Fatal("expected invalid kind error")
	}
}

func TestTaskValidateRejectsNonObjectInputs(t *testing.T) {
	task := validTask()
	task.Inputs = json.RawMessage(`"thread-1"`)

	if err := task.Validate(); err == nil {
		t.Fatal("expected non-object inputs error")
	}
}

func validTask() *Task {
	return &Task{
		TaskID:         "sha256:abc123",
		Kind:           "summarize.thread",
		Version:        "1.0.0",
		Inputs:         json.RawMessage(`{"thread_id":"thread-1"}`),
		IdempotencyKey: "sha256:abc123",
		Provenance: Provenance{
			ProducedBy: Producer{AgentID: "cli"},
			ProducedAt: time.Now(),
			StepID:     "step:abc123",
		},
	}
}
