package swimProtocol

import (
	"fmt"
)

// Channel registration identifier
const (
	FAILURE_DETECT_GET = iota
	FAILURE_DETECT_PUT
	ADD_NODE
	REMOVE_NODE
)

// Configuration identifier
const (
	WAIT_TIME = iota
	PING_INTERVAL
)

// Default configuration parameters
const (
	WAIT_TIME_DEFAULT     = 1000
	PING_INTERVAL_DEFAULT = 2000
)

// The node in SWIM protocol need to maintain two fields,
// the state is of type State, the update cound is an integer,
// indicating how long has the node not been updated
type node interface {
	State() interface{}
	UpdateCount() int
}

type nodeManager interface {
	AllNodes() []node
	AddNode(<-chan *node)
	RemoveNode(<-chan *node)
}

type SwimProtocol struct {
	// list of all member ndoes
	members []node
	// which node is the next to ping
	index int

	// when get a new failure detection state change, tell node manager to
	// change node state correspondingly
	failureDetectorIncomingChannel <-chan *node
	failureDetectorOutgoingChannel chan<- *node
	addNodeChannel                 chan<- *node
	removeNodeChannel              chan<- *node

	// configuration
	// after ping how many MILLISECONDS with no response is thought as unreachable, default 1000 ms
	waitTime int64
	//the interval gap of two pings, in MILLISECONDS, default 2000 ms
	pingInterval int64
}

func NewSwimProtocol(conf map[int]interface{}) *SwimProtocol {
	var waitTime int64
	var pingInterval int64
	if waitTime = WAIT_TIME_DEFAULT; conf[WAIT_TIME] != nil {
		waitTime = conf[WAIT_TIME].(int64)
	}
	if pingInterval = PING_INTERVAL_DEFAULT; conf[PING_INTERVAL] != nil {
		waitTime = conf[PING_INTERVAL].(int64)
	}

	swim := &SwimProtocol{
		waitTime:     waitTime,
		pingInterval: pingInterval,
	}
	return swim
}

func (swim *SwimProtocol) RegPutChannel(kind int, channel <-chan *node) {
	switch kind {
	case FAILURE_DETECT_PUT:
		swim.failureDetectorIncomingChannel = channel
	default:
		panic(fmt.Sprintf("kind %d doesn't supported", kind))
	}
}

func (swim *SwimProtocol) RegGetChannel(kind int, channel chan<- *node) {
	switch kind {
	case FAILURE_DETECT_GET:
		swim.failureDetectorOutgoingChannel = channel
	case ADD_NODE:
		swim.addNodeChannel = channel
	case REMOVE_NODE:
		swim.removeNodeChannel = channel
	default:
		panic(fmt.Sprintf("kind %d doesn't supported", kind))
	}

}

func (swim *SwimProtocol) Configure(key int, value interface{}) {
	switch key {
	case WAIT_TIME:
		swim.waitTime = value.(int64)
	case PING_INTERVAL:
		swim.pingInterval = value.(int64)
	}
}

func (swim *SwimProtocol) Run(nm nodeManager) {
	// check all channel has been registered
	if swim.failureDetectorIncomingChannel == nil {
		panic("failure detector incoming channel need be registered")
	}
	if swim.failureDetectorOutgoingChannel == nil {
		panic("failure detector outgoing channel need be registered")
	}
	if swim.addNodeChannel == nil {
		panic("add node channel need be registered")
	}
	if swim.removeNodeChannel == nil {
		panic("remove node channel need be registered")
	}

	//

}

func (swim *SwimProtocol) updateMembers(nm nodeManager) {
	swim.members = nm.AllNodes()
}

func (swim *SwimProtocol) pingNext() {

}
