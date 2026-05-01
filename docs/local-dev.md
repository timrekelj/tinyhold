# Local Development

## Running the Server

```bash
cd server
go run ./cmd/server
```

The server listens on `:7777` by default.

## Running the Client

1. Open Godot 4 editor.
2. Import `client/project.godot`.
3. Press **F5** or click the play button.

## Local Single-Player (Integrated Server)

Planned workflow:
- Godot launcher starts `server/cmd/server` on a random localhost port.
- Client connects automatically to `127.0.0.1:<port>`.
- On exit, Godot gracefully shuts down the local server process.

## Asset Pipeline

- Aseprite source files live in `client/assets/`.
- Export PNGs from Aseprite into the same directory for Godot to import.
- Godot automatically imports `.png` files on project open.
