package p2p

import (
	"encoding/gob"
	"io"
)

// used to decode the incoming data
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

//anyone implementing this interface must provide a
// method that reads from a stream and decodes into v

type GOBDecoder struct{}

func (g *GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct{}

func (g *DefaultDecoder) Decode(r io.Reader, msg *RPC) error {

	peekBuf := make([]byte, 1)

	if _, err := r.Read(peekBuf); err != nil {
		return err
	}

	// In case of a stream we are not decoding what is being sent over the network.
	// We are just setting Stream true so we can handle that in our logic.
	stream := peekBuf[0] == IncommingStream
	if stream {
		msg.Stream = stream
		return nil
	}

	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]
	return nil

}
