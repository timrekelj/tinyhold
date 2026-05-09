package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tinyhold/server/internal/server"
)

func main() {
	local := flag.Bool("local", false, "Bind to localhost only (127.0.0.1)")
	port := flag.Int("port", 42069, "Port to listen on")
	seed := flag.Int64("seed", 0, "World seed (0 = random)")
	flag.Parse()

	bind := ""
	if *local {
		bind = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%d", bind, *port)

	srv := server.New(addr, *seed)
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
