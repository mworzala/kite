package proxy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type Proxy struct {
	listener net.Listener
	ctx      context.Context

	clientInitializer func(p *Player)
}

func New(clientInitializer func(p *Player)) *Proxy {
	return &Proxy{
		listener: nil, ctx: nil,

		clientInitializer: clientInitializer,
	}
}

func (p *Proxy) Start(ctx context.Context) error {
	const listenAddr = "localhost:25577"

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	p.listener = listener
	go p.clientListenLoop()

	println("started listening on", listenAddr)

	return nil
}

func (p *Proxy) Stop(ctx context.Context) error {
	if err := p.listener.Close(); err != nil {
		panic(err) //todo
	}
	return nil
}

func (p *Proxy) clientListenLoop() {
	for {
		cc, err := p.listener.Accept()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("Failed to accept connection: %s\n", err)
			continue
		}

		conn, readLoop := proto.NewConn(packet.Serverbound, cc)
		player := &Player{conn}
		p.clientInitializer(player)
		go readLoop()
	}
}
