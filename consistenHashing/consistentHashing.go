package consistentHashing

import (
	"errors"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/otnt/distributed-system-how-to/node"
	"sync"
)

//The abstruct structure of consistent hash ring.
//It consists of a Red-Black-Tree serve as the ring.
type Ring struct {
	Tree *rbt.Tree
	mux  sync.Mutex
}

//Create a new consistent hashing ring with default
//value setting.
//
//@return: pointer to new created consistent hashing ring
func NewRing() (ring *Ring) {
	ring = &Ring{}
	ring.Tree = rbt.NewWithStringComparator()
	return
}

// Add a new Node to consistent hashing ring.
//
// @param addChan: incoming channel of Node pointer
// @param complete: outgoing channel of Node pointer, indicating
//                   which Node has been added
func (ring *Ring) Add(addChan <-chan *node.Node, complete chan<- *node.Node) {
	for {
		node := <-addChan
		keys := node.Keys

		ring.mux.Lock()
		for _, key := range keys {
			ring.Tree.Put(key, node)
		}
		ring.mux.Unlock()

		complete <- node
	}
}

// Remove a Node from consistent hashing ring.
//
// @param removeChan: incoming channel of Node pointer
// @param complete: outgoing channel of Node pointer, indicating
//                  which Node has been removed
func (ring *Ring) Remove(removeChan <-chan *node.Node, complete chan<- *node.Node) {
	for {
		node := <-removeChan
		keys := node.Keys

		ring.mux.Lock()
		for _, key := range keys {
			ring.Tree.Remove(key)
		}
		ring.mux.Unlock()

		complete <- node
	}
}

// Lookup for a Node in consistent hashing ring.
// This function is actually redirected to NextOf
//
// @param key: string of key
// @return: return the node if such successor founded, otherwise an error is
//          given
func (ring *Ring) Lookup(key string) (node.Node, error) {
	return ring.NextOf(key)
}

// Given a key on consistent hashing ring, it returns the nearest successor
// of this key.
//
// @param key: string of key
// @return: return the node if such successor founded, otherwise an error is
//          given
func (ring *Ring) NextOf(key string) (node.Node, error) {
	value, found := ring.Tree.Get(key)

	if found {
		node := value.(node.Node)
		return node, nil
	}
	return node.Node{}, errors.New("key doesn't exist")
}
