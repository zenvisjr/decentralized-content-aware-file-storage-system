package p2p

import (
	"net"
	"testing" // Go's built-in testing package
	"time"

	"github.com/stretchr/testify/assert"
	// Adds powerful assertions (true, equal, nil, etc.)
)

// func TestTCPTransport(t *testing.T) {
// 	listenAddr := ":3000"
// 	tr, err := NewTCPTransport(listenAddr)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.Equal(t, listenAddr, tr.listenerPortAddr)

// 	assert.Nil(t, tr.ListenAndAccept())
// }

func TestTCPConnection(t *testing.T) {
	tcpOps := TCPTransportOps{
		ListenerPortAddr: ":3000",
		ShakeHands:       NOHandshakeFunc,
		Decoder:          &DefaultDecoder{},
	}
	tr, err := NewTCPTransport(tcpOps)
	assert.NoError(t, err)
	assert.NoError(t, tr.ListenAndAccept())

	conn, err := net.Dial("tcp", ":3000")
	assert.NoError(t, err)
	defer conn.Close()
	_, err = conn.Write([]byte("hello Zenvis"))
	assert.NoError(t, err)

	// Wait a little to let AcceptAndLoop process
	time.Sleep(500 * time.Millisecond)

}
