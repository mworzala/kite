package kite

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/mworzala/kite/pkg/proto"
	"github.com/mworzala/kite/pkg/proto/packet"
)

type Logger func(format string, v ...interface{})

type Proxy struct {
	ListenAddr        string
	ClientInitializer func(p *Player)

	Log      Logger // Defaults to log.Printf
	ErrorLog Logger // Defaults to log.Printf

	listener net.Listener
	ctx      context.Context
}

// Start starts the proxy server on the given ListenAddr.
//
// Start will not block after creating the listener, Stop should be used to stop the proxy.
func (p *Proxy) Start() (err error) {
	if err = p.validateParamsAndSetDefaults(); err != nil {
		return err
	}

	p.listener, err = net.Listen("tcp", p.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	go p.clientListenLoop()

	p.Log("started listening on %s", p.ListenAddr)

	return nil
}

func (p *Proxy) Stop() error {
	if err := p.listener.Close(); err != nil {
		return err
	}
	return nil
}

func (p *Proxy) clientListenLoop() {
	for {
		cc, err := p.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			break
		} else if err != nil {
			log.Printf("Failed to accept connection: %s\n", err)
			continue
		}

		conn, readLoop := proto.NewConn(packet.Serverbound, cc)
		player := &Player{Conn: conn}
		p.ClientInitializer(player)
		go readLoop()
	}
}

func (p *Proxy) validateParamsAndSetDefaults() error {
	if p.ListenAddr == "" {
		return errors.New("listen address is required")
	}

	if p.Log == nil {
		p.Log = log.Printf
	}
	if p.ErrorLog == nil {
		p.ErrorLog = log.Printf
	}

	return nil
}
