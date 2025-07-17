package p2p

import (
	"fmt"
	"log"
	"net"
	"sync"
)


// TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {
	// The underlying TCP connection of the peer.
	net.Conn
	// if we dial and retrieve a conn => outbound == true
	// if we accept and retrieve a conn => outbound == false
	outbound bool

	wg *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

func (t *TCPPeer) Send(b []byte) error {
	_, err := t.Conn.Write(b)
	return err
}

func(t *TCPPeer) CloseStream() {
	t.wg.Done()

}


// TCPTransportOps represents the options for a TCPTransport.
type TCPTransportOps struct {
	ListenerPortAddr string
	ShakeHands       HandshakeFunc
	Decoder          Decoder
	OnPeer           func(Peer) error
}

// TCPTransport represents a TCP transport connection.
type TCPTransport struct {
	TCPTransportOps
	listener net.Listener
	rpcch    chan RPC
	incomingStreamChan chan Peer
}

// NewTCPTransport creates a new TCPTransport instance.
func NewTCPTransport(ops TCPTransportOps) (t *TCPTransport, err error) {
	return &TCPTransport{
		TCPTransportOps: ops,
		rpcch:           make(chan RPC, 1024),
		incomingStreamChan: make(chan Peer, 1024),
	}, nil
}

// ListenAddr implements the Transport interface, which returns the listener address.
func(t *TCPTransport) ListenAddr() string {
	return t.ListenerPortAddr
}

// Consume implements the Tranport interface, which will return read-only channel
// for reading the incoming messages received from another peer in the network.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

// ConsumeStream implements the Tranport interface, which will return read-only channel
// for reading the incoming messages received from another peer in the network.
func (t *TCPTransport) ConsumeStream() <-chan Peer {
	return t.incomingStreamChan
}

// Close implements the Transport interface
// Closes the listener and stops the AcceptAndLoop
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial implements the Transport interface, which dials a connection to the remote node.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	// fmt.Println("error: ", err)
	if err != nil {
		return err
	}
	// defer conn.Close()

	go t.HandleConn(conn, true)
	return nil
}

// ListenAndAccept implements the Transport interface, which listens and accepts connections from remote nodes.
func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.ListenerPortAddr)
	if err != nil {
		return err
	}
	t.listener = ln
	go t.acceptAndLoop()
	log.Println("TCP port listening on", t.ListenerPortAddr)
	return nil
}

//acceptAndLoop accepts connections from remote nodes and spawn a new go routine to handle each connection.
func (t *TCPTransport) acceptAndLoop() {
	for {
		conn, err := t.listener.Accept()
		if net.ErrClosed == err {
			return
		}
		if err != nil {
			fmt.Println("Error accepting tcp connection", err)
			continue
		}

		go t.HandleConn(conn, false)
	}
}

//HandleConn handles the connection from remote node.
func (t *TCPTransport) HandleConn(conn net.Conn, outbound bool) {
	var err error
	defer func() {
		fmt.Println("droping peer connection", err)
		conn.Close()
	}()
	peer := NewTCPPeer(conn, outbound)
	fmt.Printf("New connection from %s to %s\n", peer.RemoteAddr(), peer.LocalAddr())

	if err := t.ShakeHands(peer); err != nil {
		fmt.Println("Error shaking hands", err)
		return
	}

	//if OnPeer is not nil -> a function is provided, call it
	if t.OnPeer != nil {
		//if OnPeer returns an error, drop the peer connection
		if err = t.OnPeer(peer); err != nil {
			return
		}
		//if OnPeer does not return an error, continue
	}

	for {
		// fmt.Println("Waiting for message from peer", peer.RemoteAddr().String())
		rpc := RPC{}
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			fmt.Println("Error decoding", err)
			return
		}
		rpc.From = conn.RemoteAddr().String()
		if rpc.Stream {
			t.incomingStreamChan <- peer
			peer.wg.Add(1)
			fmt.Printf("incomming stream from peer [%s], WAITING....\n", peer.RemoteAddr().String())
			peer.wg.Wait()
			fmt.Printf("stream closed from peer [%s], RESUMING READ LOOP....\n", peer.RemoteAddr().String())
			continue
		}
		// fmt.Printf("Received message %+v\n", rpc)
		//now insted of printing the decoded msg here that we recieved from peer,
		//  we will send it to the rpcch channel
		t.rpcch <- rpc
		
		fmt.Println("Sending again to rpcch channel")
	}

}

