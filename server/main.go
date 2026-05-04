package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const PORT = "42069"

type Player struct {
	ID   uint16
	X, Y int32
	Keys uint8
	Conn net.Conn
}

var (
	playersMu    sync.Mutex
	players             = make(map[uint16]*Player)
	nextPlayerID uint16 = 1
)

type PacketType byte

const (
	PacketHandshake PacketType = iota
	PacketWorldChunk
	PacketPlayerInput
	PacketPlayerState
)

func main() {
	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Printf("Server listening on port %s\n", PORT)

	go gameLoop()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn)
	}
}

func gameLoop() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	const speed int32 = 200

	for range ticker.C {
		playersMu.Lock()

		type snapshot struct {
			ID   uint16
			X    int32
			Y    int32
			Conn net.Conn
		}
		snaps := make([]snapshot, 0, len(players))

		for _, p := range players {
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

			snaps = append(snaps, snapshot{
				ID:   p.ID,
				X:    p.X,
				Y:    p.Y,
				Conn: p.Conn,
			})
		}

		playersMu.Unlock()

		var buf []byte
		for _, s := range snaps {
			pkt := make([]byte, 13)
			pkt[0] = byte(PacketPlayerState)
			binary.LittleEndian.PutUint16(pkt[1:3], 10) // payload length
			binary.LittleEndian.PutUint16(pkt[3:5], s.ID)
			binary.LittleEndian.PutUint32(pkt[5:9], uint32(s.X))
			binary.LittleEndian.PutUint32(pkt[9:13], uint32(s.Y))
			buf = append(buf, pkt...)
		}

		for _, s := range snaps {
			s.Conn.Write(buf)
		}
	}

}

func handleClient(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr()
	fmt.Printf("New connection: %s\n", addr)

	var player *Player

	for {
		header := make([]byte, 3)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			fmt.Printf("Connection closed: %s\n", addr)
			fmt.Printf("Header err: %s\n", err)
			if player != nil {
				removePlayer(player.ID)
			}
			return
		}

		packetType := PacketType(header[0])
		length := binary.LittleEndian.Uint16(header[1:3])

		payload := make([]byte, length)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			fmt.Printf("Connection closed: %s\n", addr)
			fmt.Printf("Payload err: %s\n", err)
			if player != nil {
				removePlayer(player.ID)
			}
			return
		}

		if packetType == PacketHandshake && length == 1 {
			playersMu.Lock()
			id := nextPlayerID
			nextPlayerID++

			player = &Player{
				ID:   id,
				X:    0,
				Y:    0,
				Keys: 0,
				Conn: conn,
			}
			players[id] = player
			playersMu.Unlock()

			response := make([]byte, 3)
			response[0] = 1
			binary.LittleEndian.PutUint16(response[1:3], id)
			respond(conn, PacketHandshake, response)
			continue
		} else if packetType == PacketPlayerInput && length == 5 && player != nil {
			player.Keys = payload[0]
		}
	}
}

func respond(conn net.Conn, ptype PacketType, payload []byte) error {
	buf := make([]byte, 3+len(payload))
	buf[0] = byte(ptype)
	binary.LittleEndian.PutUint16(buf[1:3], uint16(len(payload)))
	copy(buf[3:], payload)
	_, err := conn.Write(buf)
	return err
}

func removePlayer(id uint16) {
	playersMu.Lock()
	delete(players, id)
	playersMu.Unlock()
	fmt.Printf("Removed player %d\n", id)
}
