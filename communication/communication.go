package communication

import (
	"github.com/otnt/distributed-system-notes/node"
	"net"
	"encoding/json"
	"io"
)

// communication protocol
const (
	TCP = iota
	UDP
)

type Message struct {
	SequenceId int64 //beginning from 0
	Sender string //uuid
	Receiver string //uuid
	Kind int //message type
	Data interface{}
}

type Communication struct{
	Protocol int
	ProtocolName string
}

func NewCommunication(protocol int) *Communication {
	if protocol != TCP {
		panic("Unsupported protocol, currently support: TCP")
	}

	var protocolName string
	if protocol == TCP {
		protocolName = "tcp"
	}

	return &Communication{Protocol:protocol, ProtocolName:protocolName}
}

func (comm *Communication) Send(n *node.Node, msg *Message) error {
	m, err := json.Marshal(msg)
	if err != nil {
		panic("Error when marshaling data")
	}

	conn, err := net.Dial(comm.ProtocolName, n.Uuid)
	if err != nil {
		return err
	}

	num, err := conn.Write(m)
	if err != nil || num != len(m) {
		return err
	}

	return nil
}

func (comm *Communication) SendAndReceive(n *node.Node, msg *Message) (*Message, error) {
	m, err := json.Marshal(msg)
	if err != nil {
		panic("Error when marshaling data")
	}

	conn, err := net.Dial(comm.ProtocolName, n.Uuid)
	if err != nil {
		return nil, err
	}

	num, err := conn.Write(m)
	if err != nil || num != len(m) {
		return nil, err
	}

	buf := make([]byte, 1024)
	data := make([]byte, 1024)
	for n:=1; n > 0; n, err = conn.Read(buf) {
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}
		data = append(data, buf...)
	}

	var res *Message
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}


