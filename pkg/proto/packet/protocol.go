package packet

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
