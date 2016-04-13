package communication

import(
	"github.com/otnt/distributed-system-notes/node"
	"testing"
	"github.com/stretchr/testify/assert"
	"net"
	"bufio"
	"encoding/json"
)

func serverListen(t *testing.T) net.Listener{
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		t.Fatal(err)
	}
	return l
}

func serverAccept(t *testing.T, l net.Listener) net.Conn {
	conn, err := l.Accept()
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func initial() (*Communication, *node.Node, *Message) {
	comm := NewCommunication(TCP)
	node := &node.Node{Uuid:"127.0.0.1:3000"}
	msg := &Message{0, 0, "sender", "receiver", 0, "data"}
	return comm, node, msg
}

func TestSend(t *testing.T) {
	l := serverListen(t)
	defer l.Close()

	go func() {
		comm, node, msg := initial()
		comm.Send(node, msg)
	}()

	conn := serverAccept(t, l)
	tmp, _ := bufio.NewReader(conn).ReadBytes(DELIM)
	tmp = tmp[:len(tmp)-1]
	res := string(tmp)
	m, _:= json.Marshal(&Message{0,0,"sender","receiver",0,"data"})
	assert.Equal(t, res, string(m))
}

func TestSendAndReceive(t *testing.T) {
	l := serverListen(t)
	defer l.Close()

	finish:=make(chan bool)
	go func(finish chan bool) {
		comm, node, msg := initial()

		res, _ := comm.SendAndReceive(node, msg)
		m:= &Message{1,1,"sender","receiver",0,"ack"}
		assert.Equal(t, res,m)

		finish <- true
	}(finish)

	conn := serverAccept(t, l)

	data, _ := bufio.NewReader(conn).ReadBytes(DELIM)
	data = data[:len(data)-1]
	var res = new(Message)
	json.Unmarshal(data, res)

	assert.Equal(t, res, &Message{0,0,"sender","receiver",0,"data"})

	m, _:= json.Marshal(&Message{1,1,"sender","receiver",0,"ack"})
	conn.Write(append(m,DELIM))

	<-finish
}
