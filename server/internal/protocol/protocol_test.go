package protocol

import (
	"bytes"
	"io"
	"testing"
)

func TestWriteReadRoundtrip(t *testing.T) {
	tests := []struct {
		name    string
		ptype   PacketType
		payload []byte
	}{
		{"empty payload", PacketHandshake, []byte{}},
		{"single byte", PacketPlayerInput, []byte{0x01}},
		{"small payload", PacketPlayerState, []byte{0x01, 0x02, 0x03}},
		{"known handshake", PacketHandshake, []byte{0x01}},
		{"known input 5 bytes", PacketPlayerInput, []byte{0x01, 0x00, 0x00, 0x00, 0x00}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			if err := Write(&buf, tt.ptype, tt.payload); err != nil {
				t.Fatalf("Write failed: %v", err)
			}

			ptype, payload, err := Read(&buf)
			if err != nil {
				t.Fatalf("Read failed: %v", err)
			}

			if ptype != tt.ptype {
				t.Errorf("type = %d, want %d", ptype, tt.ptype)
			}
			if len(payload) != len(tt.payload) {
				t.Errorf("payload len = %d, want %d", len(payload), len(tt.payload))
			}
			for i := range tt.payload {
				if payload[i] != tt.payload[i] {
					t.Errorf("payload[%d] = %02x, want %02x", i, payload[i], tt.payload[i])
				}
			}
		})
	}
}

func TestWriteReadMaxPayload(t *testing.T) {
	payload := make([]byte, 65535)
	for i := range payload {
		payload[i] = byte(i % 256)
	}

	var buf bytes.Buffer
	if err := Write(&buf, PacketPlayerState, payload); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	_, result, err := Read(&buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(result) != 65535 {
		t.Errorf("len = %d, want 65535", len(result))
	}
	if result[0] != 0x00 || result[65534] != 254 {
		t.Errorf("payload corrupted at boundaries")
	}
}

func TestReadEOF(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"empty stream", []byte{}},
		{"header truncated", []byte{0x00, 0x01}},
		{"payload truncated", []byte{0x00, 0x04, 0x00, 0x01, 0x02}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewReader(tt.data)
			_, _, err := Read(buf)
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				t.Errorf("expected EOF error, got: %v", err)
			}
		})
	}
}

func TestPacketTypeConstants(t *testing.T) {
	if PacketHandshake != 0 {
		t.Errorf("PacketHandshake = %d, want 0", PacketHandshake)
	}
	if PacketWorldChunk != 1 {
		t.Errorf("PacketWorldChunk = %d, want 1", PacketWorldChunk)
	}
	if PacketPlayerInput != 2 {
		t.Errorf("PacketPlayerInput = %d, want 2", PacketPlayerInput)
	}
	if PacketPlayerState != 3 {
		t.Errorf("PacketPlayerState = %d, want 3", PacketPlayerState)
	}
}
