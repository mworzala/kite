package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/mworzala/kite"
	"github.com/mworzala/kite/internal/pkg/handler"
	"github.com/mworzala/kite/pkg/proto/packet"
)

func main() {
	println("Hello, World!")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	p := &kite.Proxy{
		ListenAddr: "localhost:25577",
		ClientInitializer: func(p *kite.Player) {
			p.SetState(packet.Handshake, handler.NewServerboundHandshakeHandler(p))
		},
	}
	if err := p.Start(); err != nil {
		panic(err)
	}

	<-ctx.Done()

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()

	_ = stopCtx
	if err := p.Stop(); err != nil {
		panic(err)
	}

	println("Goodbye, World!")
}
