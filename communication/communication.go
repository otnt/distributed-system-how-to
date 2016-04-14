package communication

import (
	"net"
	"encoding/json"
	"bufio"
	"sync"
	"time"
)

// communication protocol
const (
	TCP = iota
	UDP
)

const (
	DELIM = byte(0)
)

const (
	MEMBERSHIP = iota
)

type node interface {
	Uuid() string
}

type Message struct {
	SequenceId int64 //beginning from 0
	LastSequenceId int64 //beginning from 0
	Sender string //uuid
	Receiver string //uuid
	Kind int //message type
	Data interface{}
}

type Communication struct{
	Protocol int
	ProtocolName string
	regMap *(map[int](chan *MsgConn))
	regMapLock *(sync.Mutex)
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

type MsgConn struct {
	Msg *Message
	Conn net.Conn
}

func (comm *Communication) Init() {
	l, err := net.Listen(comm.ProtocolName, ":80")
	if err != nil {
		panic("Could not listen on port")
	}

	comm.regMap = new((map[int](chan *MsgConn)))
	comm.regMapLock = new(sync.Mutex)
	regMap := comm.regMap
	regMapLock := comm.regMapLock

	go func() {
		for {
			conn, err :=l.Accept()
			if err != nil {
				continue
			}

			go func() {
				data, _ := bufio.NewReader(conn).ReadBytes(DELIM)
				data = data[:len(data)-1] //remove DELIM

				var res = new(Message)
				err = json.Unmarshal(data, res)
				if err != nil {
					return
				}

				kind := res.Kind
				regMapLock.Lock()
				if val, ok := (*regMap)[kind]; !ok {
					(*regMap)[kind] = make(chan *MsgConn)
					regMapLock.Unlock()
				} else {
					regMapLock.Unlock()
					val <- &MsgConn{res, conn}
				}
			}()
		}
	}()

	return
}

func (comm *Communication) ReceiveRegister(kind int) (rec chan<- *MsgConn, err error){
	go func() {
		var ok bool
		for {
			if rec,ok = (*comm.regMap)[kind]; ok {
				break
			} else {
				<-time.After(time.Second * 5)
			}
		}

	}()
	err = nil
	return
}

func (comm *Communication) Send(n node, msg *Message) error {
	m, err := json.Marshal(msg)
	if err != nil {
		panic("Error when marshaling data")
	}

	conn, err := net.Dial(comm.ProtocolName, n.Uuid())
	if err != nil {
		return err
	}
	defer conn.Close()

	num, err := conn.Write(append(m,DELIM))
	if err != nil || num != len(m) {
		return err
	}

	return nil
}

func (comm *Communication) SendAndReceive(n node, msg *Message) (*Message, error) {
	m, err := json.Marshal(msg)
	if err != nil {
		panic("Error when marshaling data")
	}

	conn, err := net.Dial(comm.ProtocolName, n.Uuid())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(append(m,DELIM))
	if err != nil  {
		return nil, err
	}

	data, _ := bufio.NewReader(conn).ReadBytes(DELIM)
	data = data[:len(data)-1] //remove DELIM

	var res = new(Message)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (comm *Communication) SendAndReceiveChannel(n node, msg *Message) (<-chan *Message) {
	res := make(chan *Message)

	go func() {
		m, err := json.Marshal(msg)
		if err != nil {
			panic("Error when marshaling data")
		}

		conn, err := net.Dial(comm.ProtocolName, n.Uuid())
		defer conn.Close()
		if err != nil {
			close(res)
			return
		}

		_, err = conn.Write(append(m,DELIM))
		if err != nil  {
			close(res)
			return
		}

		data, _ := bufio.NewReader(conn).ReadBytes(DELIM)
		data = data[:len(data)-1] //remove DELIM

		var resMsg = new(Message)
		err = json.Unmarshal(data, resMsg)
		if err != nil {
			close(res)
			return
		}

		res <- resMsg
	}()

	return (<-chan *Message)(res)
}
