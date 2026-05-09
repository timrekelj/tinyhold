package game

const (
	TileSize  = 16
	ChunkSize = 16
)

const (
	TileGrass uint8 = iota
	TilePath
)

type TileCoord struct {
	X int32
	Y int32
}

func TileAt(x, y int32, seed int64) uint8 {
	rx := x / 8
	ry := y / 8
	h := int64(rx)*374761393 + int64(ry)*668265263 + seed*1274126177
	h = (h ^ (h >> 13)) * 1274126177
	h = h ^ (h >> 16)

	if uint32(h)%100 < 25 {
		fineH := int64(x)*1619 + int64(y)*31337 + seed*54321
		fineH = (fineH ^ (fineH >> 8)) * 1103515245
		fineH = fineH ^ (fineH >> 16)
		if uint32(fineH)%100 < 80 {
			return TilePath
		}
	}
	return TileGrass
}

func ChunkCoord(tileX, tileY int32) (int32, int32) {
	cx := tileX / ChunkSize
	cy := tileY / ChunkSize
	if tileX < 0 {
		cx = (tileX+1)/ChunkSize - 1
	}
	if tileY < 0 {
		cy = (tileY+1)/ChunkSize - 1
	}
	return cx, cy
}
