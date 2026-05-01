package net

import (
	"fmt"
	"net"
)

// Start listens on addr and accepts TCP connections.
func Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer ln.Close()

	fmt.Printf("Server listening on %s\n", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	// TODO: handshake, packet read loop
}
