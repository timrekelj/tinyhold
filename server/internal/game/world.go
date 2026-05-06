package game

import (
	"sync"
)

// Player represents a single player in the world simulation.
type Player struct {
	ID   uint16
	X, Y int32
	Keys uint8
}

// Snapshot is an authoritative state update for one player.
type Snapshot struct {
	ID uint16
	X  int32
	Y  int32
}

const speed int32 = 200 // units per tick (20 tps => 4000 u/s)

// World holds the authoritative simulation state.
type World struct {
	mu      sync.Mutex
	players map[uint16]*Player
	nextID  uint16
}

// NewWorld creates an empty world.
func NewWorld() *World {
	return &World{
		players: make(map[uint16]*Player),
		nextID:  1,
	}
}

// AddPlayer creates a new player at the origin and returns it.
func (w *World) AddPlayer() *Player {
	w.mu.Lock()
	defer w.mu.Unlock()
	p := &Player{ID: w.nextID}
	w.nextID++
	w.players[p.ID] = p
	return p
}

// RemovePlayer deletes a player from the world.
func (w *World) RemovePlayer(id uint16) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.players, id)
}

// UpdateInput sets the key-state mask for a player.
func (w *World) UpdateInput(id uint16, keys uint8) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if p, ok := w.players[id]; ok {
		p.Keys = keys
	}
}

// Tick advances simulation by one step and returns snapshots for all players.
func (w *World) Tick() []Snapshot {
	w.mu.Lock()
	defer w.mu.Unlock()

	snaps := make([]Snapshot, 0, len(w.players))
	for _, p := range w.players {
		if p.Keys&1 != 0 {
			p.Y -= speed
		} // up
		if p.Keys&2 != 0 {
			p.Y += speed
		} // down
		if p.Keys&4 != 0 {
			p.X -= speed
		} // left
		if p.Keys&8 != 0 {
			p.X += speed
		} // right

		snaps = append(snaps, Snapshot{
			ID: p.ID,
			X:  p.X,
			Y:  p.Y,
		})
	}
	return snaps
}
