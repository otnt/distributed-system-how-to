package main

import (
	"fmt"
	ch "github.com/otnt/distributed-system-notes/consistentHashing"
	"github.com/otnt/distributed-system-notes/node"
)

func main() {
	//always create ring by using NewRing method
	ring := ch.NewRing()

	//Add/Remove Node to consistent hashing ring in asynchronous method
	nodes := []*node.Node{
		node.NewNode("127.0.0.1", 0, 3),
		node.NewNode("127.0.0.1", 1, 3),
		node.NewNode("127.0.0.1", 2, 3),
		node.NewNode("127.0.0.1", 3, 3),
		node.NewNode("127.0.0.1", 4, 3),
	}

	//need two channels to do this, one serves as task input channel,
	//the other serves as task completion notification channel
	task := make(chan *node.Node, 5)
	complete := make(chan *node.Node, 5)

	//run add node routine
	go ring.AddAsync(task, complete)
	task <- nodes[0]
	task <- nodes[1]
	task <- nodes[2]
	task <- nodes[3]
	task <- nodes[4]
	<-complete //when node is added successfully, complete channel get value back
	<-complete
	<-complete
	<-complete
	<-complete
	//use close to stop routine, the complete channel would be closed as well
	close(task)

	//run remove node routine, only 127.0.0.1:2 and 127.0.0.1:3 would be saved
	task = make(chan *node.Node, 5)
	go ring.RemoveAsync(task, complete)
	task <- nodes[0]
	task <- nodes[3]
	task <- nodes[4]
	<-complete //when node is removed successfully, complete channel get value back
	<-complete
	<-complete
	//use close to stop routine, the complete channel would be closed as well
	close(task)

	//also, you could add or remove node in a synchronize method
	ring.AddSync(nodes[3])
	ring.RemoveSync(nodes[3])

	//look up a key, get the coordinator node for this key
	//127.0.0.1:1
	//[97c1af2272de15532b1483651b715129332f8406 7260a48008fb01d884067d8e50b64ac56b9c3221 eb102fa9386db4715c2cfc93d019ca21c194b767]
	//127.0.0.1:2
	//[bd5b206633d9f79501860b0c03559379b89435ff c98a2ee624a8306f46ef5e01f6ea5dcce0b7ac52 98273171fa35b563b7519cef47d48005ee391f2c]
	n1, _ := ring.LookUp("97c1af2272de15532b1483651b715129332f8406")
	n2, _ := ring.LookUp("7260a48008fb01d884067d8e50b64ac56b9c3221")
	n3, _ := ring.LookUp("eb102fa9386db4715c2cfc93d019ca21c194b767")
	fmt.Println(n1 == n2, n1 == n3, n1 == nodes[1])

	//given a key string, get successor of this node
	succ, _ := ring.Successor("97c1af2272de15532b1483651b715129332f8406")
	fmt.Println(succ == nodes[2])
	succ, _ = ring.Successor("7260a48008fb01d884067d8e50b64ac56b9c3221")
	fmt.Println(succ == nodes[1])
	succ, _ = ring.Successor("eb102fa9386db4715c2cfc93d019ca21c194b767")
	fmt.Println(succ == nodes[1])

	//only when there is no node in the ring, the look up will fail
	//notice you should check this error field, since it is possible
	//for the consistent hashing to be empty
	ring.RemoveSync(nodes[1])
	ring.RemoveSync(nodes[2])
	_, err := ring.LookUp("97c1af2272de15532b1483651b715129332f8406")
	fmt.Println(err)
}
