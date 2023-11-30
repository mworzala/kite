package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
)

// Represents a single client <> server connection
type Worker struct {
	client      net.Conn
	clientState GameState

	server      net.Conn
	serverState GameState
}

func NewWorker(client net.Conn, server net.Conn) *Worker {
	w := &Worker{
		client:      client,
		clientState: Handshake,
		server:      server,
		serverState: Handshake,
	}

	go w.pipe(client, server, Serverbound)
	go w.pipe(server, client, Clientbound)

	return w
}

// Processes a single clientbound packet, returning true to drop the packet, false to forward it
func (w *Worker) processClientboundPacket(p *Packet) (drop bool, err error) {
	return false, nil
}

// Processes a single serverbound packet, returning true to drop the packet, false to forward it
func (w *Worker) processServerboundPacket(p *Packet) (drop bool, err error) {
	return false, nil
}

func (w *Worker) pipe(src net.Conn, dst net.Conn, direction PacketDirection) {
	buffer := make([]byte, 4096)
	var remainder []byte

	for {
		n, err := src.Read(buffer)
		if n == 0 || errors.Is(err, io.EOF) {
			return
		}
		//todo if n!=0 we still need to process those bytes before the error
		if err != nil {
			panic(err)
		}

		// Prepend the remainder from the previous read
		full := buffer[0:n]
		if len(remainder) > 0 {
			full = append(remainder, full...)
			remainder = nil
		}

		for len(full) > 0 {
			reader := bytes.NewReader(full)

			// read the length of the packet
			pLength, err := readVarInt(reader)
			if err != nil {
				panic(err)
			}

			// If there is not enough data left for the entire packet, keep it for the next read
			if pLength > reader.Len() {
				remainder = make([]byte, len(full))
				copy(remainder, full)
				break
			}
			readIndex := uint32(reader.Size() - int64(reader.Len()))

			// Decode the packet data
			p, err := Decode(full[readIndex : readIndex+uint32(pLength)])
			if err != nil {
				panic(err)
			}

			fmt.Printf("Packet: %d\n%s\n", p.PacketID, string(p.Data))

			var drop bool
			if direction == Clientbound {
				drop, err = w.processClientboundPacket(&p)
			} else {
				drop, err = w.processServerboundPacket(&p)
			}
			if err != nil {
				panic(err)
			}

			if !drop {
				_, err = dst.Write(full[0 : readIndex+uint32(pLength)])
				if err != nil {
					panic(err)
				}
			}

			full = full[readIndex+uint32(pLength):]
		}
	}
}
