package world

import (
	"math"

	perlin "github.com/aquilax/go-perlin"
)

const (
	TileSize  = 16
	ChunkSize = 16
)

const (
	TileGrass = iota
	TilePath
)

// Chunk represents a square region of the world.
type Chunk struct {
	X, Y  int
	Tiles [ChunkSize][ChunkSize]byte
}

// World holds all loaded chunks and generation state.
type World struct {
	Chunks map[[2]int]*Chunk
	Seed   int64
}

// NewWorld creates a world with the given generation seed.
func NewWorld(seed int64) *World {
	return &World{
		Chunks: make(map[[2]int]*Chunk),
		Seed:   seed,
	}
}

// GetChunk returns an existing chunk or generates a new one.
func (w *World) GetChunk(cx, cy int) *Chunk {
	key := [2]int{cx, cy}
	if chunk, ok := w.Chunks[key]; ok {
		return chunk
	}

	chunk := &Chunk{X: cx, Y: cy}
	pPath := perlin.NewPerlin(2, 2, 1, w.Seed+1000)

	for y := range ChunkSize {
		for x := range ChunkSize {
			wx := float64(cx*ChunkSize + x)
			wy := float64(cy*ChunkSize + y)

			path := pPath.Noise2D(wx, wy)
			isPath := math.Abs(path) < 0.1

			if isPath {
				chunk.Tiles[y][x] = TilePath
			} else {
				chunk.Tiles[y][x] = TileGrass
			}
		}
	}

	w.Chunks[key] = chunk
	return chunk
}
