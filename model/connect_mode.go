package model

type CONNECT_MODE int

const (
	CONNECT_MODE_PROXY = iota
	CONNECT_MODE_DIRECT
)

func (c CONNECT_MODE) String() string {
	switch c {
	case CONNECT_MODE_PROXY:
		return "PROXY"
	case CONNECT_MODE_DIRECT:
		return "DIRECT"
	default:
		return "Unknown"
	}
}