package main

import (
	"context"
	"errors"
	"github.com/mworzala/kite/pkg/mojang"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/packet"
)

func main() {
	println("Hello, World!")

	keyPair, err := mojang.GenerateKeyPair()
	if err != nil {
		panic(err)
	}

	proxy := &Proxy{
		MojKeyPair:     keyPair,
		VelocitySecret: "abcdef",
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	listener, err := net.Listen("tcp", "localhost:25577")
	if err != nil {
		panic(err)
	}
	go clientListenLoop(proxy, listener)

	<-ctx.Done()

	if err := listener.Close(); err != nil {
		panic(err)
	}

	println("Goodbye, World!")
}

func clientListenLoop(proxy *Proxy, listener net.Listener) {
	for {
		cc, err := listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			break
		} else if err != nil {
			log.Printf("Failed to accept connection: %s\n", err)
			continue
		}

		p := Player{proxy: proxy}
		p.conn = kite.NewConn(packet.Serverbound, cc, p.handleClientPacket)
		go p.conn.ReadLoop()
	}
}
