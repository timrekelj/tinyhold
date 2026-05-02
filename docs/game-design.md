# Tinyhold Game Design

## Overview
A 2D topdown pixel art survival game inspired by Minecraft's core loop — gather, craft, build, survive — but viewed from above. Supports single-player and local/online multiplayer.

## Core Loop
1. **Gather** materials from the world (wood, stone, ore, plants)
2. **Craft** tools, weapons, and building components
3. **Build** bases and structures for shelter and storage
4. **Survive** against environmental threats and enemies

## World

### Generation
- **Large procedurally generated world** built from square chunks (Minecraft-style chunking adapted to 2D)
- Chunks load dynamically around the player as they explore
- Biome-based generation (forest, desert, tundra, plains, etc.)
- Seed-based — same seed produces identical world

### Chunk System
- Server owns chunk data, sends `WorldChunk` packets to clients
- Chunks contain tile data (ground type, resources, structures)
- Compression planned for chunk transmission to reduce bandwidth
- Server tracks which chunks each client has loaded for delta updates

### Later: Multiple Worlds
- Portal structures scattered across the main world
- Each portal leads to a unique dimension with its own biome, enemies, and loot tables
- Boss encounters tied to specific worlds
- Players can return to the main world between runs

## Player

### Movement & Controls
- Topdown WASD movement
- Mouse aiming and attack direction
- Inventory management (grid-based, like Minecraft)

### Survival Mechanics
- Health, hunger, stamina (details to be designed)
- Day/night cycle affecting visibility and enemy behavior
- Crafting table, furnace, and other stations for advanced recipes

### Progression
- Better tools unlock faster gathering and new materials
- Base building expands storage and comfort
- Portal exploration provides high-tier items and challenges

## Multiplayer

### Modes
- **Single-player** — local server instance, seamless experience
- **Co-op** — join a friend's world, shared progression
- **Online multiplayer** — dedicated server, persistent world

### Server Authority
- Server is authoritative for all game state
- Clients send input, server simulates and broadcasts state
- Client-side interpolation for smooth movement

## Development Phases

### Phase 1: Foundation
- Connection, handshake, heartbeat
- Chunk loading and rendering
- Player movement (input → server → state broadcast)
- Basic world rendering from chunks

### Phase 2: Gathering & Crafting
- Resource nodes on the map (trees, rocks, ore veins)
- Tool system (axe, pickaxe, sword tiers)
- Inventory UI and item management
- Crafting recipes and crafting station

### Phase 3: Building & Interior
- Placeable blocks (walls, floors, doors, windows)
- Base structures with interior rooms
- Storage containers, beds, decoration items
- Player-owned claims/territories

### Phase 4: Survival Threats
- Enemies (hostile mobs, nocturnal creatures)
- Day/night cycle with visual and gameplay impact
- Combat system (melee, ranged, dodging)
- Health/respawn mechanics

### Phase 5: Multiple Worlds
- Portal structures and activation
- New dimensions with unique biomes, enemies, bosses
- Dimension-specific items and resources
- Boss encounters and loot tables

## Technical

### Engine & Language
- **Client**: Godot 4 (GDScript) — rendering, UI, input, interpolation
- **Server**: Go — authoritative simulation, networking, persistence
- **Networking**: TCP now, UDP for state sync later
- **Persistence**: In-memory → JSON → SQLite

### Key Systems
| System | Owner | Notes |
|--------|-------|-------|
| World generation | Server | Seeded procedural |
| Chunk streaming | Server → Client | Dynamic load/unload |
| Entity simulation | Server | Mobs, items, projectiles |
| Input processing | Client → Server | Key states, mouse aim |
| Rendering | Client | 2D pixel art, tile-based |
| Persistence | Server | SQLite (planned) |
