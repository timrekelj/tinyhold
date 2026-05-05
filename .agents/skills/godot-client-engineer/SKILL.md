# Skill: godot-client-engineer

## Purpose

Use this skill for Godot 4 client work in Tinyhold: scene structure, GDScript, rendering, input, local prediction/interpolation, UI, and multiplayer client integration.

## When To Use

- Implementing or reviewing Godot client features.
- Working with GDScript, nodes, scenes, signals, resources, or UI.
- Integrating the Godot client with the authoritative Go server.
- Debugging Godot multiplayer, WebSocket, or replication behavior.

## Approach

1. Preserve existing Godot scene and script conventions.
2. Keep the server authoritative; client state should be visual, predictive, or interpolated unless explicitly designed otherwise.
3. Prefer simple node composition and clear signal flow over deeply coupled scene logic.
4. Validate networking behavior against Tinyhold's protocol docs and server implementation.

## References

- Godot stable documentation: https://docs.godotengine.org/en/stable/
- Godot high-level multiplayer: https://docs.godotengine.org/en/latest/tutorials/networking/high_level_multiplayer.html
- Godot WebSocketPeer: https://docs.godotengine.org/en/stable/classes/class_websocketpeer.html
