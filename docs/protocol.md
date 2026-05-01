# Tinyhold Network Protocol

## Transport
- TCP for reliable packets (handshake, chunk data, inventory)
- UDP may be added later for entity snapshots

## Packet Format
All packets start with a 3-byte header:

| Field | Type   | Description                |
|-------|--------|----------------------------|
| Type  | uint8  | PacketType enum            |
| Len   | uint16 | Payload length (big-endian)|

## Packet Types

| Value | Name           | Direction | Description                      |
|-------|----------------|-----------|----------------------------------|
| 0     | Handshake      | C↔S       | Connection setup, player ID      |
| 1     | WorldChunk     | S→C       | Full chunk tile data             |
| 2     | PlayerInput    | C→S       | Keys pressed, mouse aim          |
| 3     | PlayerState    | S→C       | Position, velocity, animation    |
| 4     | EntityUpdate   | S→C       | Mob/item positions and states    |

## TODO
- Define exact payload layouts
- Add compression for chunk packets
- Consider delta compression for entity snapshots
