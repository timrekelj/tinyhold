package game

import (
	"crypto/rand"
	"encoding/binary"
	"sync"
)

type Player struct {
	ID   uint16
	X, Y int32
	Keys uint8
}

type Snapshot struct {
	ID uint16
	X  int32
	Y  int32
}

type BlockUpdate struct {
	X    int32
	Y    int32
	Type uint8
}

const speed int32 = 200

type World struct {
	mu            sync.Mutex
	players       map[uint16]*Player
	nextID        uint16
	Seed          int64
	placedBlocks  map[TileCoord]uint8
	pendingBlocks []BlockUpdate
}

func NewWorld(seed int64) *World {
	if seed == 0 {
		seed = randomSeed()
	}
	return &World{
		players:      make(map[uint16]*Player),
		nextID:       1,
		Seed:         seed,
		placedBlocks: make(map[TileCoord]uint8),
	}
}

func randomSeed() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0x1A2B3C4D5E6F7081
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
}

func (w *World) AddPlayer() *Player {
	w.mu.Lock()
	defer w.mu.Unlock()
	p := &Player{ID: w.nextID}
	w.nextID++
	w.players[p.ID] = p
	return p
}

func (w *World) RemovePlayer(id uint16) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.players, id)
}

func (w *World) UpdateInput(id uint16, keys uint8) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if p, ok := w.players[id]; ok {
		p.Keys = keys
	}
}

func (w *World) Tick() []Snapshot {
	w.mu.Lock()
	defer w.mu.Unlock()

	snaps := make([]Snapshot, 0, len(w.players))
	for _, p := range w.players {
		if p.Keys&1 != 0 {
			p.Y -= speed
		}
		if p.Keys&2 != 0 {
			p.Y += speed
		}
		if p.Keys&4 != 0 {
			p.X -= speed
		}
		if p.Keys&8 != 0 {
			p.X += speed
		}

		snaps = append(snaps, Snapshot{
			ID: p.ID,
			X:  p.X,
			Y:  p.Y,
		})
	}
	return snaps
}

func (w *World) WorldTileAt(x, y int32) uint8 {
	if b, ok := w.placedBlocks[TileCoord{x, y}]; ok {
		return b
	}
	return TileAt(x, y, w.Seed)
}

func (w *World) PlaceBlock(x, y int32) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	coord := TileCoord{x, y}
	if _, exists := w.placedBlocks[coord]; exists {
		return false
	}
	if TileAt(x, y, w.Seed) == TilePath {
		return false
	}

	w.placedBlocks[coord] = TileGrass + 2
	w.pendingBlocks = append(w.pendingBlocks, BlockUpdate{X: x, Y: y, Type: TileGrass + 2})
	return true
}

func (w *World) GetPendingBlocks() []BlockUpdate {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]BlockUpdate, len(w.pendingBlocks))
	copy(out, w.pendingBlocks)
	w.pendingBlocks = w.pendingBlocks[:0]
	return out
}

func (w *World) GetPlacedBlocks() []BlockUpdate {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]BlockUpdate, 0, len(w.placedBlocks))
	for coord, typ := range w.placedBlocks {
		out = append(out, BlockUpdate{X: coord.X, Y: coord.Y, Type: typ})
	}
	return out
}

func (w *World) GenerateChunk(cx, cy int32) []byte {
	tiles := make([]byte, ChunkSize*ChunkSize)
	w.mu.Lock()
	defer w.mu.Unlock()

	baseX := cx * ChunkSize
	baseY := cy * ChunkSize

	for y := int32(0); y < ChunkSize; y++ {
		for x := int32(0); x < ChunkSize; x++ {
			tx := baseX + x
			ty := baseY + y
			if b, ok := w.placedBlocks[TileCoord{tx, ty}]; ok {
				tiles[y*ChunkSize+x] = b
			} else {
				tiles[y*ChunkSize+x] = TileAt(tx, ty, w.Seed)
			}
		}
	}
	return tiles
}
