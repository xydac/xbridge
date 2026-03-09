# xbridge

A zero-config build proxy for Mac mini. Single Go binary, no dependencies.

## Quick Reference

- **Language:** Go 1.22+
- **Build:** `go build ./cmd/xbridge` or `make build`
- **Test:** `go test ./...` or `make test`
- **Run:** `go run ./cmd/xbridge serve`

## Project Structure

- `cmd/xbridge/` - Entry point (cobra root command)
- `cli/` - CLI command definitions (one file per command)
- `internal/server/` - HTTP server, routing, middleware
- `internal/server/handlers/` - HTTP handlers (health, build, simulator, git)
- `internal/engine/` - Thin wrappers around xcodebuild, simctl, git (HTTP-agnostic)
- `internal/build/` - Build job lifecycle, queuing, artifact management
- `internal/config/` - Config loading, profile resolution
- `internal/client/` - HTTP client for CLI client mode

## Key Design Patterns

- Engine layer is HTTP-agnostic (returns structs, not JSON)
- Build manager uses channels for log streaming
- SSE is hand-rolled (text/event-stream + fprintf)
- "BOOTED" is a magic UDID that resolves to the currently booted simulator
- Config precedence: CLI flags > xbridge.yaml > auto-detect

## Dependencies

- chi (HTTP router)
- cobra (CLI framework)
- gorilla/websocket (WebSocket)
- yaml.v3 (config parsing)
