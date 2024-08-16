package main

import (
	"DFS/p2p"
	"fmt"
	"log"
)

func OnPeer(peer p2p.Peer) error {
	// return fmt.Errorf("failed onpeer func")
	peer.Close()
	fmt.Println("doing some logic with the peer outside of tcptransport")
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tcp := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tcp.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tcp.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}

	// kingdom := p2p.Kingdom {
	// 	Animal: p2p.Animal {
	// 		Name: "vinay",
	// 		Age: 32,
	// 		Race: "human",
	// 	},
	// }

	// fmt.Println(kingdom.Name)
}
