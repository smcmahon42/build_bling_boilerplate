// Package producer defines the interface every task-kind implementation
// satisfies. One kind = one Producer. The Registry maps Task.kind to its
// Producer; the handler looks up and invokes.
package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/example/project/mcp-server/internal/primitives"
)

type Producer interface {
	// Kind returns the dotted kind this producer serves (e.g., "summarize.thread").
	Kind() string

	// Produce runs the task. Return a filled Result.output on success, or a
	// non-nil *primitives.Error on failure. The handler wraps either into the
	// final Result (adding result_id, provenance, metrics).
	Produce(ctx context.Context, task *primitives.Task) (output json.RawMessage, perr *primitives.Error)
}

// Registry dispatches Tasks to Producers by kind.
type Registry struct {
	byKind map[string]Producer
}

func NewRegistry(producers ...Producer) *Registry {
	r := &Registry{byKind: make(map[string]Producer, len(producers))}
	for _, p := range producers {
		r.byKind[p.Kind()] = p
	}
	return r
}

func (r *Registry) Lookup(kind string) (Producer, bool) {
	p, ok := r.byKind[kind]
	return p, ok
}

// --- Example producer: summarize.thread ------------------------------------

// SummarizeThread is a mock implementation that demonstrates the Producer
// shape for the example task kind in templates/contracts/json-schema/.
// Replace with a real implementation when adopting.
type SummarizeThread struct {
	Now func() time.Time
}

func (SummarizeThread) Kind() string { return "summarize.thread" }

func (s SummarizeThread) Produce(_ context.Context, task *primitives.Task) (json.RawMessage, *primitives.Error) {
	var in struct {
		ThreadID string `json:"thread_id"`
	}
	if err := json.Unmarshal(task.Inputs, &in); err != nil {
		return nil, &primitives.Error{Code: "input.invalid", Message: err.Error()}
	}
	out := struct {
		Summary      string `json:"summary"`
		MessageCount int    `json:"message_count"`
	}{
		Summary:      "TODO: replace with real summarization of thread " + in.ThreadID,
		MessageCount: 0,
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, &primitives.Error{Code: "producer.serialize", Message: err.Error()}
	}
	return b, nil
}
