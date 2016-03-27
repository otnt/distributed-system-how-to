package consistentHashing

import (
	"errors"
	//"fmt"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/otnt/distributed-system-how-to/node"
	"sort"
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

//Add a new Node to consistent hashing ring.
//This function block the running routine until the Node
//is successfully added.
//
//@param node: Node to be added
func (ring *Ring) AddSync(node *node.Node) {
	ring.mux.Lock()
	for _, key := range node.Keys {
		ring.Tree.Put(key, node)
	}
	ring.mux.Unlock()

	return
}

//Remove a new Node to consistent hashing ring.
//This function block the running routine until the Node
//is successfully removed.
//
//@param node: Node to be removed
func (ring *Ring) RemoveSync(node *node.Node) {
	ring.mux.Lock()
	for _, key := range node.Keys {
		ring.Tree.Remove(key)
	}
	ring.mux.Unlock()

	return
}

// Add a new Node to consistent hashing ring.
//
// @param addChan: incoming channel of Node pointer
// @param complete: outgoing channel of Node pointer, indicating
//                   which Node has been added
func (ring *Ring) AddAsync(addChan <-chan *node.Node, complete chan<- *node.Node) {
	for node := range addChan {
		keys := node.Keys

		ring.mux.Lock()
		for _, key := range keys {
			ring.Tree.Put(key, node)
		}
		ring.mux.Unlock()

		complete <- node
	}
	close(complete)
}

// Remove a Node from consistent hashing ring.
//
// @param removeChan: incoming channel of Node pointer
// @param complete: outgoing channel of Node pointer, indicating
//                  which Node has been removed
func (ring *Ring) RemoveAsync(removeChan <-chan *node.Node, complete chan<- *node.Node) {
	for node := range removeChan {
		keys := node.Keys

		ring.mux.Lock()
		for _, key := range keys {
			ring.Tree.Remove(key)
		}
		ring.mux.Unlock()

		complete <- node
	}
	close(complete)
}

// Lookup for a Node in consistent hashing ring.
// This function is actually redirected to NextOf
//
// @param key: string of key
// @return: return the node if such successor founded, otherwise an error is
//          given
func (ring *Ring) LookUp(key string) (*node.Node, error) {
	size := ring.Tree.Size()
	if size == 0 {
		return &node.Node{}, errors.New("No node alive")

	}

	interfaceSlice := ring.Tree.Keys()
	keys := make([]string, size)
	for i := range keys {
		keys[i] = interfaceSlice[i].(string)
	}

	index := sort.SearchStrings(keys, key)
	if index == size {
		index = 0
	}

	return ring.Tree.Values()[index].(*node.Node), nil
}

// Given a key on consistent hashing ring, it returns the nearest successor
// of this key.
//
// @param key: string of key
// @return: return the node if such successor founded, otherwise an error is
//          given
func (ring *Ring) NextOf(key string) (*node.Node, error) {
	return ring.LookUp(AddOne(key))
}

//Add key by one, the key is usually a very large number, mixing
//number & alphabetic letter. This method is used specifically for
//adding SHA1 result by 1
func AddOne(key string) string {
	next := map[rune]rune{
		'0': '1',
		'1': '2',
		'2': '3',
		'3': '4',
		'4': '5',
		'5': '6',
		'6': '7',
		'7': '8',
		'8': '9',
		'9': 'a',
		'a': 'b',
		'b': 'c',
		'c': 'd',
		'd': 'e',
		'e': 'f',
		'f': '0',
	}

	newKey := []rune(key)
	for i := len(newKey) - 1; i >= 0; i-- {
		newKey[i] = next[newKey[i]]
		if newKey[i] != '0' {
			break
		}
	}

	return string(newKey)
}

//Substract key by one, the key is usually a very large number, mixing
//number & alphabetic letter. This method is used specifically for
//substractign SHA1 result by 1
func SubOne(key string) string {
	next := map[rune]rune{
		'0': 'f',
		'1': '0',
		'2': '1',
		'3': '2',
		'4': '3',
		'5': '4',
		'6': '5',
		'7': '6',
		'8': '7',
		'9': '8',
		'a': '9',
		'b': 'a',
		'c': 'b',
		'd': 'c',
		'e': 'd',
		'f': 'e',
	}

	newKey := []rune(key)
	for i := len(newKey) - 1; i >= 0; i-- {
		newKey[i] = next[newKey[i]]
		if newKey[i] != 'f' {
			break
		}
	}

	return string(newKey)
}
