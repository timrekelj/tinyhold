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
| 5     | Heartbeat      | C↔S       | Keepalive, latency measurement   |
| 6     | Disconnect     | C↔S       | Graceful leave                   |
| 7     | Inventory      | S→C       | Item updates, slot changes       |

## Payload Layouts

### Handshake (0)
**C→S (connect request):**

| Field   | Type   | Description                     |
|---------|--------|---------------------------------|
| Version | uint8  | Protocol version                |

**S→C (connect response):**

| Field    | Type   | Description                     |
|----------|--------|---------------------------------|
| Version  | uint8  | Protocol version                |
| PlayerID | uint16 | Assigned player ID              |

If version mismatch, server sends Version=0, PlayerID=0 and closes connection.

### PlayerInput (2)

| Field    | Type   | Description                     |
|----------|--------|---------------------------------|
| Keys     | uint8  | Bitmask of movement keys        |
| AimX     | int16  | Mouse aim X (relative)          |
| AimY     | int16  | Mouse aim Y (relative)          |

### PlayerState (3)

| Field    | Type   | Description                     |
|----------|--------|---------------------------------|
| PlayerID | uint16 | Player ID                       |
| X        | int32  | World X (sub-pixel)             |
| Y        | int32  | World Y (sub-pixel)             |
| VelX     | int16  | Velocity X                      |
| VelY     | int16  | Velocity Y                      |
| Anim     | uint8  | Current animation frame         |

### EntityUpdate (4)

| Field    | Type   | Description                     |
|----------|--------|---------------------------------|
| Count    | uint8  | Number of entity updates        |
| N ×:     |        |                                 |
| EntityID | uint16 | Entity ID                       |
| Type     | uint8  | Entity type                     |
| X        | int32  | World X                         |
| Y        | int32  | World Y                         |
| State    | uint8  | Entity-specific state           |

### Heartbeat (5)

| Field    | Type   | Description                     |
|----------|--------|---------------------------------|
| Seq      | uint32 | Sequence number                 |

Server echoes the same Seq back for RTT measurement.

### Disconnect (6)

| Field    | Type   | Description                     |
|----------|--------|---------------------------------|
| Reason   | uint8  | 0=client quit, 1=kicked, 2=error |

### Inventory (7)

| Field    | Type   | Description                     |
|----------|--------|---------------------------------|
| Slot     | uint8  | Slot index                      |
| ItemID   | uint16 | Item type                       |
| Count    | uint8  | Stack count                     |

## TODO
- Add compression for chunk packets
- Consider delta compression for entity snapshots
- Define WorldChunk payload layout
- Define Keys bitmask for PlayerInput
