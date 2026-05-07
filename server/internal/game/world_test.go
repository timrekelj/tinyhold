package game

import (
	"sync"
	"testing"
)

func TestAddPlayerSequentialIDs(t *testing.T) {
	w := NewWorld()

	p1 := w.AddPlayer()
	p2 := w.AddPlayer()
	p3 := w.AddPlayer()

	if p1.ID != 1 || p2.ID != 2 || p3.ID != 3 {
		t.Errorf("IDs = %d, %d, %d; want 1, 2, 3", p1.ID, p2.ID, p3.ID)
	}
}

func TestAddPlayerSpawnsAtOrigin(t *testing.T) {
	w := NewWorld()
	p := w.AddPlayer()

	if p.X != 0 || p.Y != 0 {
		t.Errorf("spawn position = (%d, %d); want (0, 0)", p.X, p.Y)
	}
	if p.Keys != 0 {
		t.Errorf("initial keys = %d; want 0", p.Keys)
	}
}

func TestRemovePlayer(t *testing.T) {
	w := NewWorld()
	w.AddPlayer()
	w.AddPlayer()
	w.RemovePlayer(1)

	w.mu.Lock()
	_, exists := w.players[1]
	w.mu.Unlock()
	if exists {
		t.Error("player 1 should be removed")
	}

	w.mu.Lock()
	_, exists = w.players[2]
	w.mu.Unlock()
	if !exists {
		t.Error("player 2 should still exist")
	}
}

func TestUpdateInput(t *testing.T) {
	w := NewWorld()
	p := w.AddPlayer()

	w.UpdateInput(p.ID, 1)
	if p.Keys != 1 {
		t.Errorf("keys = %d; want 1", p.Keys)
	}

	w.UpdateInput(p.ID, 0)
	if p.Keys != 0 {
		t.Errorf("keys = %d; want 0", p.Keys)
	}
}

func TestUpdateInputUnknownPlayer(t *testing.T) {
	w := NewWorld()
	w.UpdateInput(999, 1) // should not panic
}

func TestTickNoInput(t *testing.T) {
	w := NewWorld()
	p := w.AddPlayer()

	snaps := w.Tick()

	if len(snaps) != 1 {
		t.Fatalf("snapshots = %d; want 1", len(snaps))
	}
	if snaps[0].ID != p.ID {
		t.Errorf("snapshot ID = %d; want %d", snaps[0].ID, p.ID)
	}
	if snaps[0].X != 0 || snaps[0].Y != 0 {
		t.Errorf("snapshot position = (%d, %d); want (0, 0)", snaps[0].X, snaps[0].Y)
	}
}

func TestTickMovement(t *testing.T) {
	tests := []struct {
		name    string
		keys    uint8
		wantDX  int32
		wantDY  int32
	}{
		{"up", 1, 0, -speed},
		{"down", 2, 0, speed},
		{"left", 4, -speed, 0},
		{"right", 8, speed, 0},
		{"diagonal up-right", 9, speed, -speed},
		{"diagonal down-left", 6, -speed, speed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWorld()
			p := w.AddPlayer()
			w.UpdateInput(p.ID, tt.keys)

			snaps := w.Tick()

			if snaps[0].X != tt.wantDX {
				t.Errorf("X = %d; want %d", snaps[0].X, tt.wantDX)
			}
			if snaps[0].Y != tt.wantDY {
				t.Errorf("Y = %d; want %d", snaps[0].Y, tt.wantDY)
			}
		})
	}
}

func TestTickAccumulatesMovement(t *testing.T) {
	w := NewWorld()
	p := w.AddPlayer()
	w.UpdateInput(p.ID, 8) // right

	w.Tick() // X = 200
	w.Tick() // X = 400
	snaps := w.Tick() // X = 600

	if snaps[0].X != 600 {
		t.Errorf("X = %d; want 600", snaps[0].X)
	}
}

func TestTickMultiplePlayers(t *testing.T) {
	w := NewWorld()
	a := w.AddPlayer()
	b := w.AddPlayer()

	w.UpdateInput(a.ID, 8)  // right
	w.UpdateInput(b.ID, 1)  // up

	snaps := w.Tick()

	if len(snaps) != 2 {
		t.Fatalf("snapshots = %d; want 2", len(snaps))
	}

	if snaps[0].ID == a.ID && (snaps[0].X != speed || snaps[0].Y != 0) {
		t.Errorf("player A snap = (%d, %d); want (%d, 0)", snaps[0].X, snaps[0].Y, speed)
	}
	if snaps[1].ID == b.ID && (snaps[1].X != 0 || snaps[1].Y != -speed) {
		t.Errorf("player B snap = (%d, %d); want (0, %d)", snaps[1].X, snaps[1].Y, -speed)
	}
}

func TestConcurrentAddAndTick(t *testing.T) {
	w := NewWorld()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				w.AddPlayer()
				w.Tick()
			}
		}()
	}
	wg.Wait()
}

func TestConcurrentUpdateAndTick(t *testing.T) {
	w := NewWorld()
	p := w.AddPlayer()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(k uint8) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				w.UpdateInput(p.ID, k)
				w.Tick()
			}
		}(uint8(i % 16))
	}
	wg.Wait()
}
