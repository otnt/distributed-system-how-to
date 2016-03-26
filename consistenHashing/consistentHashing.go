package consistentHashing

import "node"

func Add(addChan <-chan *Node) (complete chan<- *Node) {

	return
}

func Remove(removeChan <-chan *Node) (complete chan<- *Node) {

	return
}

func Lookup(lookUpChan <-chan string) Node {

	return nil
}

func NextOf(nextChan <-chan string) Node {

	return nil
}
