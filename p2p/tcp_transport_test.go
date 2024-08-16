package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listrnAddr := ":4000"

	opts := TCPTransportOpts{
		ListenAddr:    listrnAddr,
		HandShakeFunc: NOPHandShakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tcp := NewTCPTransport(opts)
	assert.Equal(t, tcp.ListenAddr, listrnAddr)

	assert.Nil(t, tcp.ListenAndAccept())
}
