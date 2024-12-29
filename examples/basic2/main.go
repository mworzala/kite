package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/proto/packet"
)

func main() {
	println("Hello, World!")

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

	publicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}

	verifyToken := make([]byte, 16)
	_, err = rand.Read(verifyToken)
	if err != nil {
		panic(fmt.Errorf("failed to generate random verify token: %w", err))
	}

	proxy := &Proxy{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		VerifyToken: verifyToken,

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
