package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"os/signal"
	"time"

	"github.com/mworzala/kite"
	"github.com/mworzala/kite/internal/pkg/handler"
)

func main() {
	println("Hello, World!")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var err error
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		panic(err)
	}

	p := &kite.Proxy{
		ListenAddr: "localhost:25577",
		InitHandler: handler.MakeClientHandshakeHandler(handler.ClientHandshakeHandlerOpts{
			LoginHandlerFunc: handler.MakeClientMojangLoginHandler(handler.ClientMojangLoginHandlerOpts{
				PrivateKey: privateKey,
			}),
		}),
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
