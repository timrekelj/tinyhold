package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const PORT = "42069"

type PacketType byte
const (
	PacketHandshake PacketType = iota
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr()
	fmt.Printf("New connection: %s\n", addr)

	for {
		header := make([]byte, 3)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			fmt.Printf("Connection closed: %s\n", addr)
			fmt.Printf("Header err: %s\n", err)
			return
		}

		packetType := PacketType(header[0])
		length := binary.BigEndian.Uint16(header[1:3])

		payload := make([]byte, length)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			fmt.Printf("Connection closed: %s\n", addr)
			fmt.Printf("Payload err: %s\n", err)
			return
		}

		fmt.Printf("Received: type=%d, len=%d, payload=%v\n", packetType, length, payload)

		if packetType == PacketHandshake && length == 1 {
			// payload[0] = client version
			// payload[1:2] = player ID
			response := []byte{1, 0, 1}
			respond(conn, PacketHandshake, response)
			continue
		}
	}
}

func respond(conn net.Conn, ptype PacketType, payload []byte) error {
	buf := make([]byte, 3+len(payload))
	buf[0] = byte(ptype)
	binary.BigEndian.PutUint16(buf[1:3], uint16(len(payload)))
	copy(buf[3:], payload)
	_, err := conn.Write(buf)
	return err
}
