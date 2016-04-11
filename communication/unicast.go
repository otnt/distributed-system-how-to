package communication

import (
//	"fmt"
//	"net"
)


func NewCommunication(protocol int) *Communication {
	if protocol != TCP {
		panic("Unsupported protocol, currently support: TCP")
	}
	return &Communication{protocol:protocol}
}

