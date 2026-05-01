package main

import (
	"log"
	"tinyhold/server/internal/net"
)

func main() {
	log.Println("Tinyhold server starting...")
	// TODO: load config, init world, start network loop
	if err := net.Start(":7777"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
