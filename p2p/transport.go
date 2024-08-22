package p2p

// Peer is an interface that represents the remote node
type Peer interface {
	Close() error
}

// Transport is anything that handles communication between the nodes in the network.
// This can be of any form (TCP UDP, websockets, ...)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
