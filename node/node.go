package node

import (
	"crypto/sha1"
	"fmt"
)

type Node struct {
	//each Node is globally identified by ip:port
	//the 'ip:port' string is the uuid of this Node
	Ip   string
	Port int
	Uuid string

	//Key is the position on consistent hashing, because key is usually extremely
	//large(e.g. 2^160) to reduce confliction, it should be a string.
	//Also, notice we are using virtual node, so one node is actually mapped
	//to several places on the ring. To do this, we append several number(depends on
	//user configuration) at end of uuid
	//e.g. If replication number is 3, uuid is 127.0.0.1:8080
	//then we have tree keys:
	//sha1('127.0.0.1:8080-0')
	//sha1('127.0.0.1:8080-1')
	//sha1('127.0.0.1:8080-2')
	Keys []string

	//when the membership protocol knows consistent hashing has
	//saved this Node, then it's marked as saved
	Saved bool

	//State shows whether the node is alive, failure or any other states.
	//The concrete type of it is up to the failure detector.
	State interface{}
}

func NewNode(ip string, port int, vnodeNum int) (node *Node) {
	node = &Node{}
	node.Ip = ip
	node.Port = port
	node.Uuid = fmt.Sprintf("%s:%d", ip, port)

	node.Keys = make([]string, vnodeNum)
	for i := 0; i < vnodeNum; i++ {
		bytes := []byte(fmt.Sprintf("%s-%d", node.Uuid, i))
		h := sha1.New()
		h.Write(bytes)
		node.Keys[i] = fmt.Sprintf("%x", h.Sum(nil))
	}
	node.Saved = false
	return
}
