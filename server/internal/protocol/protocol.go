package protocol

import (
	"encoding/binary"
	"io"
)

// PacketType identifies the kind of network packet.
type PacketType byte

const (
	PacketHandshake    PacketType = iota // Client -> Server: join request
	PacketWorldChunk                      // Server -> Client: world data
	PacketPlayerInput                     // Client -> Server: key state
	PacketPlayerState                     // Server -> Client: position snapshot
)

// Write sends a framed packet to w: [type:1][length:2][payload:N].
func Write(w io.Writer, ptype PacketType, payload []byte) error {
	buf := make([]byte, 3+len(payload))
	buf[0] = byte(ptype)
	binary.LittleEndian.PutUint16(buf[1:3], uint16(len(payload)))
	copy(buf[3:], payload)
	_, err := w.Write(buf)
	return err
}

// Read reads a framed packet from r and returns its type and payload.
func Read(r io.Reader) (PacketType, []byte, error) {
	header := make([]byte, 3)
	if _, err := io.ReadFull(r, header); err != nil {
		return 0, nil, err
	}
	ptype := PacketType(header[0])
	length := binary.LittleEndian.Uint16(header[1:3])

	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return 0, nil, err
	}
	return ptype, payload, nil
}
