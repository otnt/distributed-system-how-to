package communication

import (
	"github.com/otnt/distributed-system-notes/node"
)

const (
	TCP = iota
	UDP
)

type messagePair struct {
	node node.Node
	data interface{}
}

type communicationBase interface{
	Protocol() string
	Send(<-chan messagePair)
	Receive(chan<- messagePair)
}
