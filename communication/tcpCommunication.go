package communication

import (
	//"github.com/otnt/distributed-system-notes/node"
)

type tcpCommunication struct {
} 

func (*tcpCommunication) Protocol() string {
	return "TCP"
}

func (*tcpCommunication) Send(<-chan messagePair) {
}

func (*tcpCommunication) Receive(chan<- messagePair) {
}

func (*tcpCommunication) send(mp messagePair) {
}
