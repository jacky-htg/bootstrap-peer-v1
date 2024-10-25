package bootstrap

import (
	"fmt"
	"net"
)

// Server menyimpan data peer dan menangani koneksi.
type Server struct {
	port     string
	pm       *PeerManager
	jsonPath string
}

// NewServer membuat instance baru server.
func NewServer(port string, jsonPath string) *Server {
	return &Server{
		port:     port,
		pm:       NewPeerManager(),
		jsonPath: jsonPath,
	}
}

// ListenAndServe memulai server untuk menerima koneksi dari client.
func (s *Server) ListenAndServe() error {
	if err := s.pm.loadPeers(s.jsonPath); err != nil {
		return fmt.Errorf("gagal membaca file peers: %v", err)
	}

	ln, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", s.port, err)
	}
	defer ln.Close()

	go s.handleShutdown()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}
