package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"os/signal"
	"time"

	"github.com/mworzala/kite"
)

type UserData struct {
}

func main() {
	println("Hello, World!")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	//ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	//defer cancel()

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
		InitHandler: kite.MakeClientHandshakeHandler(kite.ClientHandshakeHandlerOpts{
			LoginHandlerFunc: kite.MakeClientMojangLoginHandler[UserData](kite.ClientMojangLoginHandlerOpts[UserData]{
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
