package consistentHashing

import (
	"errors"
	//"fmt"
	rbte "github.com/emirpasic/gods/examples"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/otnt/distributed-system-notes/node"
	"sync"
)

//The abstruct structure of consistent hash ring.
//It consists of a Red-Black-Tree serve as the ring.
type Ring struct {
	tree *rbte.RedBlackTreeExtended
	mux  sync.Mutex
}

//Create a new consistent hashing ring with default
//value setting.
//
//@return: pointer to new created consistent hashing ring
func NewRing() (ring *Ring) {
	ring = &Ring{
		tree: &rbte.RedBlackTreeExtended{
			rbt.NewWithStringComparator(),
		},
	}
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
		ring.tree.Put(key, node)
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
		ring.tree.Remove(key)
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
			ring.tree.Put(key, node)
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
			ring.tree.Remove(key)
		}
		ring.mux.Unlock()

		complete <- node
	}
	close(complete)
}

// Lookup for a Node in consistent hashing ring given a key.
// If the input key is the same as some Node's key, then the result is
// that exact Node.
//
// @param key: string of key
// @return: return the node if such successor founded, otherwise an error is
//          given
func (ring *Ring) LookUp(key string) (*node.Node, error) {
	return ring.getCeilingOf(key)
}

// Get successor of a given key(this key is supposed to belong to a Node)
//
// @param key: string of key
// @return: return the node if such successor founded, otherwise an error is
//          given
func (ring *Ring) Successor(key string) (*node.Node, error) {
	return ring.getCeilingOf(AddOne(key))
}

// Get predecessor of a given key(this key is supposed to belong to a Node)
//
// @param key: string of key
// @return: return the node if such predecessor founded, otherwise an error is
//          given
func (ring *Ring) Predecessor(key string) (*node.Node, error) {
	return ring.getFloorOf(SubOne(key))
}

// Get a Node given a specific key.
//
// @param key: string of key
// @return: if key exists in tree, return the node and true, otherwise return
//          return nil and false
func (ring *Ring) Get(key string) (*node.Node, bool) {
	if treeNodeValue, found := ring.tree.Get(key); found {
		return treeNodeValue.(*node.Node), found
	} else {
		return nil, found
	}
}

// Check if the ring is empty/no node alive
//
// @return: true if ring is empty otherwise false
func (ring *Ring) Empty() bool {
	return ring.tree.Empty()
}

// Get all keys in the ring in ascending order
//
// @return: all keys in ascending order
func (ring *Ring) Keys() []string {
	interfaceKeys := ring.tree.Keys()
	keys := make([]string, len(interfaceKeys))
	for i, _ := range keys {
		keys[i] = interfaceKeys[i].(string)
	}
	return keys
}

// Get all values in the ring corresponding to the keys in ascending order
//
// @return: all values in the ring corresponding to the keys in ascending order
func (ring *Ring) Values() []*node.Node {
	interfaceValues := ring.tree.Values()
	values := make([]*node.Node, len(interfaceValues))
	for i, _ := range values {
		values[i] = interfaceValues[i].(*node.Node)
	}
	return values
}

// Get Ceiling of a key. A Ceiling is defined as the smallest element of all
// elements that are larger than or equal to the key. In consistent hashing, the ceiling
// of an element that is larger than any exist key, is the smallest key, indicating
// the ring property.
//
// @param key: string of key
// @return: return the node if tree is not empty, otherwise an error is
//          given
func (ring *Ring) getCeilingOf(key string) (*node.Node, error) {
	size := ring.tree.Size()
	if size == 0 {
		return &node.Node{}, errors.New("No node alive")
	}

	// find ceiling of this key
	if treeNode, found := ring.tree.Ceiling(key); found {
		return treeNode.Value.(*node.Node), nil
	} else {
		// Or it may be winded back to zero
		minTreeNodeValue, _ := ring.tree.GetMin()
		return minTreeNodeValue.(*node.Node), nil
	}
}

// Get Floor of a key. A Floor is defined as the largest element of all
// elements that are smaller than or equal to the key. In consistent hashing, the floor
// of an element that is smaller than any exist key, is the largest key, indicating
// the ring property.
//
// @param key: string of key
// @return: return the node if tree is not empty, otherwise an error is
//          given
func (ring *Ring) getFloorOf(key string) (*node.Node, error) {
	size := ring.tree.Size()
	if size == 0 {
		return &node.Node{}, errors.New("No node alive")
	}

	// find floor of this key
	if treeNode, found := ring.tree.Floor(key); found {
		return treeNode.Value.(*node.Node), nil
	} else {
		// Or it may be winded back to zero
		maxTreeNodeValue, _ := ring.tree.GetMax()
		return maxTreeNodeValue.(*node.Node), nil
	}
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
