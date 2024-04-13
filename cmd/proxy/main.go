package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/mworzala/kite/internal/pkg/handler"
	"github.com/mworzala/kite/pkg/proto/packet"
	"github.com/mworzala/kite/pkg/proxy"
)

func main() {
	println("Hello, World!")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	p := proxy.New(func(p *proxy.Player) {
		p.SetState(packet.Handshake, handler.NewServerboundHandshakeHandler(p))
	})
	if err := p.Start(ctx); err != nil {
		panic(err)
	}

	<-ctx.Done()

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()

	if err := p.Stop(stopCtx); err != nil {
		panic(err)
	}

	println("Goodbye, World!")
}
