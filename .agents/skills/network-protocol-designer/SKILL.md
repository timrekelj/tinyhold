# Skill: network-protocol-designer

## Purpose

Use this skill for Tinyhold network protocol design: packet structure, message flow, versioning, serialization, validation, state sync, and client/server compatibility.

## When To Use

- Designing or changing packets, opcodes, binary layouts, or protocol docs.
- Reviewing client/server compatibility and protocol versioning.
- Deciding what data belongs in snapshots, deltas, commands, acknowledgements, or events.
- Evaluating bandwidth, latency, ordering, reliability, and replay implications.

## Approach

1. Start from the authoritative server model and define what the client is allowed to request versus observe.
2. Keep packet formats deterministic, documented, and easy to validate.
3. Design explicit versioning or migration behavior before changing shipped packet formats.
4. Separate transport concerns from game semantics where practical.
5. Prefer compact formats only after correctness, debuggability, and compatibility are clear.

## References

- Tinyhold protocol documentation: ../../../docs/protocol.md
- Godot high-level multiplayer: https://docs.godotengine.org/en/latest/tutorials/networking/high_level_multiplayer.html
- Go documentation: https://go.dev/doc/
