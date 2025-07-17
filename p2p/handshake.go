package p2p

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

// HandshakeFunc represents a function that performs a handshake with a peer
// It takes any type of data and returns an error if the handshake fails
type HandshakeFunc func(Peer) error

// NOHandshakeFunc is a no-op handshake function that always succeeds
func NOHandshakeFunc(peer Peer) error {
	return nil
}

func SimpleHandshakeFunc(peer Peer) error {
	myID := peer.LocalAddr().String()

	// Send my ID
	_, err := fmt.Fprintf(peer, "%s\n", myID)
	if err != nil {
		return fmt.Errorf("sending ID failed: %w", err)
	}

	// Receive their ID
	reader := bufio.NewReader(peer)
	theirID, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("receiving ID failed: %w", err)
	}
	theirID = strings.TrimSpace(theirID)

	if theirID == myID {
		return fmt.Errorf("self-connection detected — closing")
	}

	fmt.Printf("Handshake success: [%s] ↔ [%s]\n", myID, theirID)
	return nil
}

func DefensiveHandshakeFunc(peer Peer) error {
	const (
		maxIDLength      = 128
		handshakeTimeout = 5 * time.Second
	)

	myID := peer.LocalAddr().String()

	// Defensive timeout
	_ = peer.SetDeadline(time.Now().Add(handshakeTimeout))
	defer peer.SetDeadline(time.Time{}) // clear after

	// Reader safety
	reader := bufio.NewReader(io.LimitReader(peer, maxIDLength+2)) // +2 for newline, etc.

	// If outbound, write first; else read first
	if tcpPeer, ok := peer.(*TCPPeer); ok && tcpPeer.outbound {
		if err := sendID(peer, myID); err != nil {
			return fmt.Errorf("sending ID failed: %w", err)
		}

		theirID, err := readPeerID(reader)
		if err != nil {
			return fmt.Errorf("receiving ID failed: %w", err)
		}

		if theirID == myID {
			return fmt.Errorf("self-connection detected: %s", theirID)
		}

		fmt.Printf("Handshake success: [%s] → [%s]\n", myID, theirID)

	} else {
		theirID, err := readPeerID(reader)
		if err != nil {
			return fmt.Errorf("receiving ID failed: %w", err)
		}

		if theirID == myID {
			return fmt.Errorf("self-connection detected: %s", theirID)
		}

		if err := sendID(peer, myID); err != nil {
			return fmt.Errorf("sending ID failed: %w", err)
		}

		fmt.Printf("Handshake success: [%s] ← [%s]\n", myID, theirID)
	}

	return nil
}

func sendID(peer Peer, id string) error {
	if len(id) > 128 {
		return fmt.Errorf("ID too long: %d", len(id))
	}
	_, err := fmt.Fprintf(peer, "%s\n", id)
	return err
}

func readPeerID(reader *bufio.Reader) (string, error) {
	id, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return "", fmt.Errorf("empty peer ID")
	}
	if len(id) > 128 {
		return "", fmt.Errorf("peer ID too long")
	}
	return id, nil
}
