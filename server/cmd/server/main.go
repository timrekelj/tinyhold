package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tinyhold/server/internal/server"
)

func main() {
	addr := ":42069"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	srv := server.New(addr)
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\nShutting down...")
	srv.Stop()
}
