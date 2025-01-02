package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	multiaddr "github.com/multiformats/go-multiaddr"
)

var peerid = flag.String("id", "", "-id <Peer ID>")
var ip = flag.String("ip", "", "-ip <IP>")
var port = flag.String("port", "", "-port <Port>")

const (
	addrformat = "/ip4/%s/tcp/%s/p2p/%s"
)

func usage() {
	fmt.Println("Usage: pinged -ip <IP> -port <Port> -id <Peer ID>")
	os.Exit(1)
}

func main() {

	flag.Parse()

	if *peerid == "" || *ip == "" || *port == "" {
		usage()
	}

	address := fmt.Sprintf(addrformat, *ip, *port, *peerid)

	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	fmt.Println(os.Args[1])

	addr, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		panic(err)
	}

	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}

	if err := node.Connect(context.Background(), *peer); err != nil {
		panic(err)
	}

	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	fmt.Println("sending 5 ping messages to", addr)
	ch := pingService.Ping(context.Background(), peer.ID)

	for i := 0; i < 5; i++ {
		res := <-ch
		fmt.Println("got ping response!", "RTT:", res.RTT)
	}

	// shut the node down
	if err := node.Close(); err != nil {
		panic(err)
	}

}
