package communication

import (
	"github.com/otnt/distributed-system-notes/node"
	"net"
)

const (
	TCP = iota
	UDP
)

type messagePair struct {
	node node.Node
	data interface{}
}

var connMap map[string]net.Conn

type Communication struct{
	Protocol string
	//Protocol() string
	//Send(<-chan messagePair)
	//Receive(chan<- messagePair)
	//send(mp messagePair)
	connMapChan *chan string
}

func Init() {
	connMapChan := make(chan string)
	go maintainConnMap(connMapChan)
}

func NewCommunication(protocol int) *Communication {

	switch (protocol) {
	case(TCP):
	return &Communication{Protocol:"tcp", connMapChan:&connMapChan}
	default:
		panic("Unsupported protocol. Supported protocols are: TCP")
	}
}

func maintainConnMap(chan string) {
	for connRequest := range(connMapChan) {
	}

}

func (comm *Communication) Send(mpChan <-chan messagePair) {
	for mp := range(mpChan) {
	}

}

