package main

import (
	"DFS/p2p"
	"log"
	"time"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO: onpeer func

	}
	transport := p2p.NewTCPTransport(tcpOpts)

	fileSeverOpts := FileServerOpts{
		StorageRoot:       "3000_files",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         transport,
	}
	s := NewFileServer(fileSeverOpts)

	go func() {
		time.Sleep(time.Second * 120)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
