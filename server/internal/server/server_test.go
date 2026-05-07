package server

import (
	"encoding/binary"
	"net"
	"testing"
	"time"

	"tinyhold/server/internal/protocol"
)

func TestStartStop(t *testing.T) {
	srv := New("127.0.0.1:0")
	if err := srv.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	srv.Stop()
}

func TestHandshake(t *testing.T) {
	srv := New("127.0.0.1:0")
	if err := srv.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer srv.Stop()

	addr := srv.listener.Addr().String()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	if err := protocol.Write(conn, protocol.PacketHandshake, []byte{0x01}); err != nil {
		t.Fatalf("Write handshake failed: %v", err)
	}

	ptype, payload, err := protocol.Read(conn)
	if err != nil {
		t.Fatalf("Read response failed: %v", err)
	}

	if ptype != protocol.PacketHandshake {
		t.Errorf("response type = %d; want %d", ptype, protocol.PacketHandshake)
	}
	if len(payload) < 3 {
		t.Fatalf("payload len = %d; want >= 3", len(payload))
	}

	version := payload[0]
	playerID := binary.LittleEndian.Uint16(payload[1:3])

	if version != 1 {
		t.Errorf("version = %d; want 1", version)
	}
	if playerID == 0 {
		t.Error("player ID should be non-zero")
	}
}

func TestHandshakeWrongPayload(t *testing.T) {
	srv := New("127.0.0.1:0")
	if err := srv.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer srv.Stop()

	addr := srv.listener.Addr().String()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	if err := protocol.Write(conn, protocol.PacketHandshake, []byte{0x01, 0x02}); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	conn.SetDeadline(time.Now().Add(500 * time.Millisecond))
	_, _, err = protocol.Read(conn)
	if err == nil {
		t.Error("expected no response for invalid handshake payload")
	}
}

func TestInputBroadcast(t *testing.T) {
	srv := New("127.0.0.1:0")
	if err := srv.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer srv.Stop()

	addr := srv.listener.Addr().String()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	if err := protocol.Write(conn, protocol.PacketHandshake, []byte{0x01}); err != nil {
		t.Fatalf("Write handshake failed: %v", err)
	}

	if _, _, err := protocol.Read(conn); err != nil {
		t.Fatalf("Read handshake response failed: %v", err)
	}

	if err := protocol.Write(conn, protocol.PacketPlayerInput, []byte{0x08, 0x00, 0x00, 0x00, 0x00}); err != nil {
		t.Fatalf("Write input failed: %v", err)
	}

	conn.SetDeadline(time.Now().Add(500 * time.Millisecond))
	ptype, payload, err := protocol.Read(conn)
	if err != nil {
		t.Fatalf("Read state failed: %v", err)
	}

	if ptype != protocol.PacketPlayerState {
		t.Errorf("packet type = %d; want %d", ptype, protocol.PacketPlayerState)
	}
	if len(payload) != 10 {
		t.Errorf("payload len = %d; want 10", len(payload))
	}
}

func TestDisconnectRemovesPlayer(t *testing.T) {
	srv := New("127.0.0.1:0")
	if err := srv.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer srv.Stop()

	addr := srv.listener.Addr().String()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}

	if err := protocol.Write(conn, protocol.PacketHandshake, []byte{0x01}); err != nil {
		t.Fatalf("Write handshake failed: %v", err)
	}
	if _, _, err := protocol.Read(conn); err != nil {
		t.Fatalf("Read handshake response failed: %v", err)
	}

	conn.Close()
	time.Sleep(100 * time.Millisecond)

	srv.mu.Lock()
	numConns := len(srv.conns)
	srv.mu.Unlock()

	if numConns != 0 {
		t.Errorf("conns = %d; want 0 after disconnect", numConns)
	}
}

func TestStopWithConnectedClient(t *testing.T) {
	srv := New("127.0.0.1:0")
	if err := srv.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	addr := srv.listener.Addr().String()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}

	if err := protocol.Write(conn, protocol.PacketHandshake, []byte{0x01}); err != nil {
		t.Fatalf("Write handshake failed: %v", err)
	}
	if _, _, err := protocol.Read(conn); err != nil {
		t.Fatalf("Read handshake response failed: %v", err)
	}

	conn.Close()
	time.Sleep(100 * time.Millisecond)
	srv.Stop()
}

func TestMultipleConnections(t *testing.T) {
	srv := New("127.0.0.1:0")
	if err := srv.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer srv.Stop()

	addr := srv.listener.Addr().String()

	for i := 0; i < 3; i++ {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatalf("Dial %d failed: %v", i, err)
		}
		defer conn.Close()

		if err := protocol.Write(conn, protocol.PacketHandshake, []byte{0x01}); err != nil {
			t.Fatalf("Write handshake %d failed: %v", i, err)
		}

		_, payload, err := protocol.Read(conn)
		if err != nil {
			t.Fatalf("Read response %d failed: %v", i, err)
		}

		id := binary.LittleEndian.Uint16(payload[1:3])
		if id == 0 {
			t.Errorf("connection %d got player ID 0", i)
		}
	}

	time.Sleep(100 * time.Millisecond)

	srv.mu.Lock()
	numConns := len(srv.conns)
	srv.mu.Unlock()

	if numConns != 3 {
		t.Errorf("conns = %d; want 3", numConns)
	}
}
