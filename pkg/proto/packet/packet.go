package packet

import (
	"io"
)

// InvalidState is the representation of a case where there is no valid ID for the packet in the given state.
// For example, a client login start packet during the play state.
const InvalidState = -1

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
