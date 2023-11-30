package main

import (
	"log"
	"net"
)

func handleClient(client net.Conn, target string) {
	// Connect to the target server
	server, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("Failed to connect to target: %s\n", err)
		client.Close()
		return
	}

	NewWorker(client, server)
}

func startProxy(listenAddr string, target string) {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %s\n", listenAddr, err)
	}
	defer listener.Close()
	log.Printf("Listening on %s, forwarding to %s\n", listenAddr, target)

	for {
		// Accept new connections
		client, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s\n", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleClient(client, target)
	}
}

func main() {
	// Specify the local address to listen on and the target server address
	var listenAddr = "localhost:25566"
	var target = "localhost:25565"

	startProxy(listenAddr, target)
}
