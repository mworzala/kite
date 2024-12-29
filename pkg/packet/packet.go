package packet

import (
	"io"
)

type Direction int

const (
	Clientbound Direction = iota
	Serverbound
)

func (d Direction) String() string {
	switch d {
	case Clientbound:
		return "clientbound"
	case Serverbound:
		return "serverbound"
	}
	return "unknown"
}

type State int

const (
	Handshake State = iota
	Status
	Login
	Config
	Play

	// InvalidState is the representation of a case where there is no valid ID for the packet in the given state.
	// For example, a client login start packet during the play state.
	InvalidState = -1
)

func (s State) Validate() bool {
	return s >= Handshake && s <= Play
}

func (s State) String() string {
	switch s {
	case Handshake:
		return "handshake"
	case Status:
		return "status"
	case Login:
		return "login"
	case Config:
		return "config"
	case Play:
		return "play"
	}
	return "unknown"
}

// A Packet is a generic interface for both client and server packets.
// The only required structure is that they have an ID in a game state (could be implemented by custom packets).
type Packet interface {
	Direction() Direction
	ID(state State) int

	Read(r io.Reader) error
	Write(w io.Writer) error
}

// stateId1 is a utility function for making assertions on the sendable state and packet IDs.
// It returns the packet ID if the actual state is the current, or InvalidState if not.
func stateId1(actual, expected State, id int) int {
	if actual != expected {
		return InvalidState
	}
	return id
}

func stateId2(actual, expected1, expected2 State, id1, id2 int) int {
	if actual == expected1 {
		return id1
	}
	if actual == expected2 {
		return id2
	}
	return InvalidState
}
