// Command cli is a thin client for the service defined in
// templates/contracts/openapi/service.yaml. It builds a Task, POSTs it, and
// renders the Result. All domain logic lives behind the contract — this
// binary only does argument parsing, HTTP plumbing, and rendering.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Minimal Task/Result shapes. In a real project, generate these from
// templates/agent-primitives/schemas/ (or a shared types package) rather
// than hand-writing them here.
type task struct {
	TaskID         string          `json:"task_id"`
	Kind           string          `json:"kind"`
	Version        string          `json:"version,omitempty"`
	Inputs         json.RawMessage `json:"inputs"`
	IdempotencyKey string          `json:"idempotency_key"`
	Provenance     provenance      `json:"provenance"`
}

type provenance struct {
	ProducedBy struct {
		AgentID string `json:"agent_id"`
	} `json:"produced_by"`
	ProducedAt time.Time `json:"produced_at"`
	StepID     string    `json:"step_id"`
	TraceID    string    `json:"trace_id,omitempty"`
	Ring       int       `json:"ring"`
}

type result struct {
	ResultID string          `json:"result_id"`
	TaskID   string          `json:"task_id"`
	Status   string          `json:"status"`
	Output   json.RawMessage `json:"output,omitempty"`
	Error    *struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		Retryable bool   `json:"retryable,omitempty"`
	} `json:"error,omitempty"`
}

func main() {
	var (
		baseURL  = flag.String("server", envOr("SERVER_URL", "http://localhost:8080/v1"), "base URL of the service")
		kind     = flag.String("kind", "summarize.thread", "task kind to invoke")
		threadID = flag.String("thread", "", "thread id (input for summarize.thread)")
	)
	flag.Parse()

	if *threadID == "" {
		fail(2, "cli.usage", "--thread is required")
	}

	token := os.Getenv("CAPABILITY_TOKEN")
	if token == "" {
		fail(2, "cli.usage", "CAPABILITY_TOKEN env var is required")
	}

	inputs, _ := json.Marshal(map[string]string{"thread_id": *threadID})
	t := task{
		Kind:    *kind,
		Version: "1.0.0",
		Inputs:  inputs,
		Provenance: provenance{
			ProducedAt: time.Now().UTC(),
			Ring:       0,
			StepID:     newID("step"),
			TraceID:    newID("trace"),
		},
	}
	t.Provenance.ProducedBy.AgentID = "cli"
	t.TaskID = hashTask(t)
	t.IdempotencyKey = t.TaskID

	res, err := submit(*baseURL, token, &t)
	if err != nil {
		fail(1, "cli.transport", err.Error())
	}

	if res.Status == "error" && res.Error != nil {
		fmt.Fprintf(os.Stderr, "error: %s: %s\n", res.Error.Code, res.Error.Message)
		os.Exit(exitCodeFor(res.Error.Code, res.Error.Retryable))
	}

	// Default rendering: pretty-print the output. Real CLIs render
	// per-kind — but the shape is the same.
	out, _ := json.MarshalIndent(json.RawMessage(res.Output), "", "  ")
	fmt.Println(string(out))
}

func submit(baseURL, token string, t *task) (*result, error) {
	body, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, baseURL+"/tasks", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", t.IdempotencyKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r result
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, fmt.Errorf("response body was not a Result: %w", err)
	}
	if resp.StatusCode >= 500 && r.Error == nil {
		return nil, errors.New(resp.Status)
	}
	return &r, nil
}

func hashTask(t task) string {
	h := sha256.New()
	h.Write([]byte(t.Kind))
	h.Write([]byte(t.Version))
	h.Write(t.Inputs)
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}

func newID(prefix string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(time.Now().UTC().Format(time.RFC3339Nano)))
	return prefix + ":" + hex.EncodeToString(h.Sum(nil))[:16]
}

func exitCodeFor(code string, retryable bool) int {
	switch code {
	case "capability.denied", "capability.expired", "capability.bad_signature", "capability.wrong_subject", "capability.missing":
		return 3
	case "input.invalid", "cli.usage":
		return 2
	}
	if retryable {
		return 75 // EX_TEMPFAIL
	}
	return 1
}

func fail(exit int, code, message string) {
	fmt.Fprintf(os.Stderr, "error: %s: %s\n", code, message)
	os.Exit(exit)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
