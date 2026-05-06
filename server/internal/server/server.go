package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"tinyhold/server/internal/game"
	"tinyhold/server/internal/protocol"
)

// Server wraps the TCP listener and game world.
type Server struct {
	addr     string
	listener net.Listener
	world    *game.World
	conns    map[uint16]net.Conn
	mu       sync.Mutex
	ticker   *time.Ticker
	quit     chan struct{}
	wg       sync.WaitGroup
}

// New creates a server that will listen on addr.
func New(addr string) *Server {
	return &Server{
		addr:  addr,
		world: game.NewWorld(),
		conns: make(map[uint16]net.Conn),
		quit:  make(chan struct{}),
	}
}

// Start opens the listener and spawns accept & broadcast goroutines.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = ln
	fmt.Printf("Server listening on %s\n", s.addr)

	s.ticker = time.NewTicker(50 * time.Millisecond)
	s.wg.Add(2)
	go s.acceptLoop()
	go s.broadcastLoop()
	return nil
}

// Stop signals shutdown and waits for goroutines to finish.
func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
	s.ticker.Stop()
	s.wg.Wait()
}

func (s *Server) acceptLoop() {
	defer s.wg.Done()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				continue
			}
		}
		s.wg.Add(1)
		go s.handleClient(conn)
	}
}

func (s *Server) broadcastLoop() {
	defer s.wg.Done()
	for {
		select {
		case <-s.ticker.C:
			s.broadcast()
		case <-s.quit:
			return
		}
	}
}

func (s *Server) broadcast() {
	snaps := s.world.Tick()

	var buf []byte
	for _, snap := range snaps {
		pkt := make([]byte, 13)
		pkt[0] = byte(protocol.PacketPlayerState)
		binary.LittleEndian.PutUint16(pkt[1:3], 10) // payload length
		binary.LittleEndian.PutUint16(pkt[3:5], snap.ID)
		binary.LittleEndian.PutUint32(pkt[5:9], uint32(snap.X))
		binary.LittleEndian.PutUint32(pkt[9:13], uint32(snap.Y))
		buf = append(buf, pkt...)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, conn := range s.conns {
		conn.Write(buf)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	addr := conn.RemoteAddr()
	fmt.Printf("New connection: %s\n", addr)

	var player *game.Player

	for {
		ptype, payload, err := protocol.Read(conn)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Connection closed: %s (%v)\n", addr, err)
			} else {
				fmt.Printf("Connection closed: %s\n", addr)
			}
			if player != nil {
				s.removePlayer(player.ID)
			}
			return
		}

		switch ptype {
		case protocol.PacketHandshake:
			if len(payload) == 1 {
				player = s.world.AddPlayer()
				s.mu.Lock()
				s.conns[player.ID] = conn
				s.mu.Unlock()

				resp := make([]byte, 3)
				resp[0] = 1
				binary.LittleEndian.PutUint16(resp[1:3], player.ID)
				if err := protocol.Write(conn, protocol.PacketHandshake, resp); err != nil {
					fmt.Printf("Failed to send handshake: %v\n", err)
				}
			}
		case protocol.PacketPlayerInput:
			if len(payload) == 5 && player != nil {
				keys := payload[0]
				s.world.UpdateInput(player.ID, keys)
			}
		}
	}
}

func (s *Server) removePlayer(id uint16) {
	s.world.RemovePlayer(id)
	s.mu.Lock()
	delete(s.conns, id)
	s.mu.Unlock()
	fmt.Printf("Removed player %d\n", id)
}
