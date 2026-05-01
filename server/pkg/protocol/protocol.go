package protocol

// PacketType identifies the kind of packet.
type PacketType byte

const (
	PacketHandshake PacketType = iota
	PacketWorldChunk
	PacketPlayerInput
	PacketPlayerState
	PacketEntityUpdate
)

// Header is the common prefix for every packet.
type Header struct {
	Type PacketType
	Len  uint16 // payload length following header
}
