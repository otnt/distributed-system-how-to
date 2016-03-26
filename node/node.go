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
	//to several places on the ring.
	Keys []string

	//when the membership protocol knows consistent hashing has
	//saved this Node, then it's marked as saved
	Saved bool
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
