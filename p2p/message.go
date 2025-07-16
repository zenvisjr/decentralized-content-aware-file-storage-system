package p2p

//What they talk about (data structures for communication)


const (
	IncommingMessage = 0x1
	IncommingStream  = 0x2
)

// RPC holds any arbitart data that is sent over
// the network between peers
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
