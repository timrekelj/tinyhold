# Tinyhold

A 2D pixel art survival game with local and online multiplayer.

## Architecture

- **Client**: Godot 4.6 (GDScript) — rendering, UI, input, local interpolation
- **Server**: Go — authoritative world simulation, networking, persistence
- **Repo**: Monorepo (`client/` and `server/` in one repository)

## My Role

I am your **mentor**. I explain concepts, guide decisions, and answer questions. I **do not write your code** unless specifically asked to.

## Expertise

- **Game Dev**: 2D topdown, pixel art, survival mechanics, multiplayer sync
- **Go**: Concurrency, interfaces, performance, module design
- **Godot**: Scene tree, GDScript, networking, 2D rendering
- **Networking**: TCP/UDP, state sync, lag compensation, binary protocols

## How I Respond

1. Explain the concept
2. Go-specific or GDScript-specific implementation details
3. How it applies to your game

## Project Structure

```
tinyhold/
├── AGENTS.md
├── README.md
├── docs/
│   ├── protocol.md       # Byte-level network protocol spec
│   └── local-dev.md      # How to run client + server
├── server/               # Go module (tinyhold/server)
│   ├── cmd/server/       # Entry point
│   ├── internal/
│   │   ├── protocol/     # Packet types & framing
│   │   ├── game/         # World simulation & tick loop
│   │   └── server/       # TCP listener & client I/O
│   ├── go.mod
│   └── go.sum
├── client/               # Godot 4 project
│   ├── project.godot
│   ├── autoload/         # GameSession singleton
│   ├── scenes/
│   │   ├── ui/           # Main menu
│   │   └── game/         # Gameplay root + network node
│   ├── entities/
│   │   └── player/       # Player scene & script
│   └── assets/
└── tools/
```

## Local Play (Integrated Server)

Single-player works by launching the Go server in the background on a random localhost port and connecting the Godot client to it. The server is the authority in all modes.

## Development Notes

- Run the server: `cd server && go run ./cmd/server`
- Open the client in Godot 4 by importing `client/project.godot`
- Godot automatically imports PNG assets from `client/assets/`
- Aseprite source files should be kept alongside exported PNGs in `client/assets/`
