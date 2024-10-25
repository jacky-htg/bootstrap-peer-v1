package bootstrap

import (
	"os"
	"sync"

	"github.com/bytedance/sonic"
)

// PeerManager mengelola daftar peers dengan thread-safe.
type PeerManager struct {
	peers []Peer
	mu    sync.Mutex
}

type Peer struct {
	Address string `json:"address"`
}

// NewPeerManager membuat instance baru dari PeerManager.
func NewPeerManager() *PeerManager {
	return &PeerManager{
		peers: make([]Peer, 0),
	}
}

// RegisterPeer menambahkan peer baru ke dalam daftar.
func (pm *PeerManager) RegisterPeer(peerAddress string) (bool, string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	peer := Peer{Address: peerAddress}

	for _, p := range pm.peers {
		if p.Address == peerAddress {
			return false, "Peer already registered"
		}
	}

	pm.peers = append(pm.peers, peer)
	return true, "Peer registered successfully"
}

// GetAllPeers mengembalikan daftar semua peers.
func (pm *PeerManager) GetAllPeers() []Peer {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	return pm.peers
}

// RemovePeer menghapus peer dari daftar.
func (pm *PeerManager) RemovePeer(peerAddress string) (bool, string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for i, p := range pm.peers {
		if p.Address == peerAddress {
			pm.peers = append(pm.peers[:i], pm.peers[i+1:]...)
			return true, "Peer removed successfully"
		}
	}

	return false, "Peer not found"
}

func (pm *PeerManager) loadPeers(jsonPath string) error {
	file, err := os.ReadFile(jsonPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return sonic.Unmarshal(file, &pm.peers)
}

func (pm *PeerManager) savePeers(jsonPath string) error {
	file, err := sonic.Marshal(pm.peers)
	if err != nil {
		return err
	}

	return os.WriteFile(jsonPath, file, 0644)
}
