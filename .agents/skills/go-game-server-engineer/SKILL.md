# Skill: go-game-server-engineer

## Purpose

Use this skill for Go server work in Tinyhold: authoritative simulation, networking, persistence, concurrency, module design, and server-side gameplay systems.

## When To Use

- Implementing or reviewing Go server features.
- Designing authoritative gameplay systems and world simulation logic.
- Working on WebSocket transport, HTTP endpoints, connection lifecycle, or packet handling.
- Debugging concurrency, goroutine ownership, backpressure, or server performance.

## Approach

1. Keep the server authoritative for gameplay decisions and persisted world state.
2. Prefer clear ownership of mutable state; avoid shared state across goroutines without explicit synchronization.
3. Keep networking code separate from simulation rules where practical.
4. Use small interfaces only when they clarify boundaries or improve testability.
5. Verify behavior with Go tests and targeted integration checks when possible.

## References

- Go documentation: https://go.dev/doc/
- net/http package: https://pkg.go.dev/net/http
- coder/websocket: https://github.com/coder/websocket
- gorilla/websocket: https://pkg.go.dev/github.com/gorilla/websocket
