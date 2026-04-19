// Command server is the MCP server entrypoint. The body here is a stub —
// bind to your chosen MCP Go SDK (official modelcontextprotocol/go-sdk,
// mark3labs/mcp-go, or equivalent) and call handler.Dispatch for each
// tool call the SDK delivers.
package main

import (
	"log"
	"time"

	"github.com/example/project/mcp-server/internal/capability"
	"github.com/example/project/mcp-server/internal/handler"
	"github.com/example/project/mcp-server/internal/idempotency"
	"github.com/example/project/mcp-server/internal/producer"
)

func main() {
	h := &handler.Handler{
		AgentID:     "{{AGENT_ID}}",
		Producers:   producer.NewRegistry(producer.SummarizeThread{Now: time.Now}),
		Capability:  capability.StubVerifier{},
		Idempotency: idempotency.NewMemoryCache(),
		Now:         time.Now,
	}

	// TODO: bind an MCP SDK here. Pseudocode:
	//
	//   srv := mcp.NewServer(...)
	//   srv.RegisterToolHandler(func(ctx context.Context, name string, args json.RawMessage) (json.RawMessage, error) {
	//       result := h.Dispatch(ctx, handler.ToolCall{Tool: name, Arguments: args})
	//       return json.Marshal(result)
	//   })
	//   log.Fatal(srv.Serve(os.Stdin, os.Stdout))

	log.Printf("mcp-server skeleton ready (agent=%s); wire an MCP SDK in main.go to serve", h.AgentID)
	_ = h
}
