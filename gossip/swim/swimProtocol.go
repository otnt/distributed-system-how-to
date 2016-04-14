package swimProtocol

import (
	//"fmt"
	"math/rand"
	comm "github.com/otnt/distributed-system-notes/communication"
	"time"
)

// Default configuration parameters
const (
	WAIT_TIME_DEFAULT     = 1000
	PING_INTERVAL_DEFAULT = 2000
	ASSIST_NODE_NUM = 3
	SYNC_NODE_STATE_NUM = 5
)

const (
	ALIVE = iota
	SUSPECT
	FAULTY
)

// The node in SWIM protocol need to maintain two fields,
// the state is of type State, the update count is an integer,
// indicating how long has the node not been updated
type node interface {
	GetSwimState() int
	SetSwimState(s int)
	Uuid() string
}


type nodeManager interface {
	AllNodes() []node
	AddNode(n node)
	RemoveNode(n node)
	GetNode(uuid string) node
}

type SwimProtocol struct {
	// list of all member ndoes
	members []node

	// number of members
	memberNum int

	// which node is the next to ping
	index int

	// after ping how many MILLISECONDS with no response is thought as unreachable, default 1000 ms
	waitTime int64

	//the interval gap of two pings, in MILLISECONDS, default 2000 ms
	pingInterval int64

	// which node to be sync
	selectIndex int

	// node manager
	nm nodeManager
}

func NewSwimProtocol(nm nodeManager) *SwimProtocol {
	swim := &SwimProtocol{
		waitTime: WAIT_TIME_DEFAULT,
		pingInterval: PING_INTERVAL_DEFAULT,
		nm: nm,
	}
	return swim
}

// Send data has field
// ALIVE:[uuid1, uuid2, ...]
// SUSPECT:[uuid1, uuid2, ...]
// FAULTY:[uuid1, uuid2, ...]
// For simplicity, we just select the nodes from start to end
func (swim *SwimProtocol) Run() {
	communication := comm.NewCommunication(comm.TCP)

	//Send
	swim.updateMembers()
	members := swim.members
	for index, member:= range members {
		//Check a ramdon member intervally
		msg,err :=communication.SendAndReceive(member,swim.getSyncMessage())

		//If the member is alive, mark it as alive
		if err == nil {
			member.SetSwimState(ALIVE)
		} else {
			//Ask other nodes to assert the liveness
			recChan := make(chan *comm.Message)
			for i:=0;i<ASSIST_NODE_NUM; i++ {
				go func(){
					// avoid send to same node
					var randIndex int
					for randIndex = index; randIndex == index; randIndex = rand.Intn(len(members)) {}

					msg := <- communication.SendAndReceiveChannel(members[randIndex],swim.getSyncMessage())
					recChan <- msg
				}()
			}

			select{
			//If one of them think it as alive, mark it as alive
			case msg = <-recChan:
				member.SetSwimState(ALIVE)
			//Otherwise, mark it as suspected
			case <-time.After(WAIT_TIME_DEFAULT * time.Millisecond):
				member.SetSwimState(SUSPECT)
			}
		}
		swim.synchronize(msg)
	}

	//Receive
	for {
		//
	}
}

type syncData [SYNC_NODE_STATE_NUM][]string

func (swim *SwimProtocol)synchronize(msg *comm.Message) {
	sd := msg.Data.(syncData)
	var n node
	for state, uuids := range sd {
		for _, uuid := range uuids {
			n = swim.nm.GetNode(uuid)
			n.SetSwimState(state)
			if state == FAULTY {
				swim.nm.RemoveNode(n)
			}
		}
	}
}


func (swim *SwimProtocol) getSyncMessage() (*comm.Message){
	sd := new(syncData)
	for i:=0;i<SYNC_NODE_STATE_NUM;i++ {
		member:=swim.members[swim.selectIndex]
		uuid:=member.Uuid()
		state:=member.GetSwimState()
		if sd[state] == nil {
			sd[state] = make([]string, 0)
		}
		sd[state] = append(sd[state], uuid)
		swim.selectIndex = (swim.selectIndex + 1) % swim.memberNum
	}
	return &comm.Message{Kind:comm.MEMBERSHIP, Data:sd}
}

// Refresh node members and shuffle them
func (swim *SwimProtocol) updateMembers() {
	members := swim.nm.AllNodes()
	shuffle := rand.Perm(len(members))
	swim.members = make([]node, len(members))
	for val, index := range shuffle {
		swim.members[val] = members[index]
	}
	swim.memberNum = len(members)
}

