package bootstrap

import (
	"bootstrap-server/pkg"
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bytedance/sonic"
)

// handleConnection menangani koneksi dan menentukan jenis request.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Baca request dari client
	data, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	var req pkg.Request
	if err := sonic.Unmarshal(data, &req); err != nil {
		fmt.Println("Invalid request format:", err)
		return
	}

	// Handle sesuai tipe request
	switch req.Type {
	case "REGISTER":
		s.registerPeer(req.Payload, conn)
	case "GET_PEERS":
		s.getAllPeers(conn)
	case "REMOVE":
		s.removePeer(req.Payload, conn)
	default:
		fmt.Println("Invalid request type")
	}
}

func (s *Server) handleShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Shutdown initiated, saving peers...")

	if err := s.pm.savePeers(s.jsonPath); err != nil {
		fmt.Println("Gagal menyimpan peers:", err)
	} else {
		fmt.Println("Peers saved successfully.")
	}
	os.Exit(0)
}

// registerPeer menambahkan peer baru ke daftar.
func (s *Server) registerPeer(peer string, conn net.Conn) {
	success, msg := s.pm.RegisterPeer(peer)
	fmt.Fprintln(conn, msg)

	if !success {
		fmt.Println("Failed to register peer:", peer)
	}
}

func (s *Server) getAllPeers(conn net.Conn) {
	peers := s.pm.GetAllPeers()

	data, err := sonic.Marshal(peers)
	if err != nil {
		fmt.Println("Error encoding peers:", err)
		return
	}
	conn.Write(append(data, '\n'))
}

func (s *Server) removePeer(peer string, conn net.Conn) {
	success, msg := s.pm.RemovePeer(peer)
	fmt.Fprintln(conn, msg)

	if !success {
		fmt.Println("Failed to remove peer:", peer)
	}
}
