package p2p

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
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

var peerPublicKey map[string]*rsa.PublicKey = make(map[string]*rsa.PublicKey)

func GetPeerPublicKey(id string) (*rsa.PublicKey, bool) {
    key, ok := peerPublicKey[id]
    return key, ok
}


func DefensiveHandshakeFunc(peer Peer) error {
	const (
		maxIDLength      = 1024
		handshakeTimeout = 5 * time.Second
	)

	myID := peer.LocalAddr().String()

	_ = peer.SetDeadline(time.Now().Add(handshakeTimeout))
	defer peer.SetDeadline(time.Time{})

	reader := bufio.NewReader(io.LimitReader(peer, maxIDLength))

	myPubKey, err := LoadPublicKey()
	if err != nil {
		return fmt.Errorf("loading my public key failed: %w", err)
	}

	if tcpPeer, ok := peer.(*TCPPeer); ok && tcpPeer.outbound {
		return outboundHandshake(peer, myID, myPubKey, reader)
	}

	return inboundHandshake(peer, myID, myPubKey, reader)
}

func inboundHandshake(peer Peer, myID string, myPubKey *rsa.PublicKey, reader *bufio.Reader) error {
		theirID, err := readPeerID(reader)
		if err != nil {
			return fmt.Errorf("receiving ID failed: %w", err)
		}

		if theirID == myID {
			return fmt.Errorf("self-connection detected: %s", theirID)
		}

		theirPubKey, err := readAndDecodePublicKey(reader)
		if err != nil {
			return fmt.Errorf("receiving public key failed: %w", err)
		}

		if err := sendID(peer, myID); err != nil {
			return fmt.Errorf("sending ID failed: %w", err)
		}

		if err := sendPeerPublicKey(myPubKey, peer); err != nil {
			return fmt.Errorf("sending public key failed: %w", err)
		}

		peerPublicKey[theirID] = theirPubKey

		fmt.Printf("Handshake success: [%s] ← [%s]\n", myID, theirID)

		return nil
	}




func outboundHandshake(peer Peer, myID string, myPubKey *rsa.PublicKey, reader *bufio.Reader) error {
	if err := sendID(peer, myID); err != nil {
		return fmt.Errorf("sending ID failed: %w", err)
	}

	if err := sendPeerPublicKey(myPubKey, peer); err != nil {
		return fmt.Errorf("sending public key failed: %w", err)
	}

	theirID, err := readPeerID(reader)
	if err != nil {
		return fmt.Errorf("receiving ID failed: %w", err)
	}

	if theirID == myID {
		return fmt.Errorf("self-connection detected: %s", theirID)
	}

	theirPubKey, err := readAndDecodePublicKey(reader)
	if err != nil {
		return fmt.Errorf("receiving public key failed: %w", err)
	}

	peerPublicKey[theirID] = theirPubKey

	fmt.Printf("Handshake success: [%s] → [%s]\n", myID, theirID)

	return nil
}
	



func sendPeerPublicKey(pub *rsa.PublicKey, peer Peer) error {
	// Marshal to DER (no PEM)
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	// Encode DER to base64 string
	b64 := base64.StdEncoding.EncodeToString(pubBytes)
	_, err = fmt.Fprintf(peer, "%s\n", b64)
	return err
}

func readAndDecodePublicKey(reader *bufio.Reader) (*rsa.PublicKey, error) {
	b64, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("reading base64 public key failed: %w", err)
	}
	b64 = strings.TrimSpace(b64)

	pubBytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("key parse failed: %w", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPubKey, nil
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
