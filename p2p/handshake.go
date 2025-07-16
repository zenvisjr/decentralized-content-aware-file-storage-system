package p2p

// HandshakeFunc represents a function that performs a handshake with a peer
// It takes any type of data and returns an error if the handshake fails
type HandshakeFunc func(Peer) error

// NOHandshakeFunc is a no-op handshake function that always succeeds
func NOHandshakeFunc(peer Peer) error {
	return nil
}
