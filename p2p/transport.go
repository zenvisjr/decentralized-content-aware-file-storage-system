package p2p

import "net"

//How nodes talk (the interface for communication)

// Peer is an interface that represents a peer in the P2P network
type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

// Transport is anything that handles the communication
// between the nodes in the network. This can be of the
// form (TCP, UDP, websockets, ...)
type Transport interface {
	// net.Conn
	ListenAddr() string
	ListenAndAccept() error
	Consume() <-chan RPC
	ConsumeStream() <-chan Peer
	Dial(string) error
	Close() error
	CheckFileHashMap(string) bool
	AddFileHashMap(string, string)
}
