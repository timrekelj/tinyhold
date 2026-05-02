# Persistence Strategy

## Current: In-Memory
All player and world state lives in RAM. This is sufficient for the initial development phase — getting handshake, movement, chunk loading, and multiplayer sync working.

## Phase 1: SQLite
When the game loop is functional and ready for persistent saves. Directly integrated into the Go server.

### Why SQLite
- **Embedded** — no external service, runs in-process alongside the Go server
- **Zero config** — single file, no separate database server to manage
- **ACID** — safe concurrent writes, no corruption on crash
- **Lightweight** — perfect for a single server instance handling one world

### Recommended library
`modernc.org/sqlite` — pure Go implementation, no CGO dependency, cross-compiles cleanly.

### What to persist
| Table | Purpose |
|-------|---------|
| `players` | Player ID, name, spawn point, last seen |
| `player_inventory` | Slot index, item ID, stack count |
| `world_chunks` | Chunk coordinates, compressed tile data |
| `world_entities` | Entity type, position, state (drops, placed items) |
| `world_structures` | Player-built bases, interior layouts |

### Save cadence
- **Player inventory** — on pickup/drop/use, and periodic autosave every 60s
- **World chunks** — when modified (block placed/destroyed), not on read
- **Player position** — every 30s or on disconnect

## Not needed
- PostgreSQL/MySQL — overkill for a single-server game
- Redis — no caching layer needed yet
- Object storage — saves are small and local
