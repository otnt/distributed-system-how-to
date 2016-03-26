package consistentHashing

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/otnt/distributed-system-how-to/node"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func initChannels() (task chan *node.Node, complete chan *node.Node) {
	task = make(chan *node.Node)
	complete = make(chan *node.Node)
	return
}

func getKeysActuallyAre(tree *rbt.Tree) []string {
	interfaceKeysInTree := tree.Keys()
	keysInTree := make([]string, 0)
	for _, key := range interfaceKeysInTree {
		keysInTree = append(keysInTree, key.(string))
	}

	return keysInTree
}

func getKeysShouldBe(nodes []*node.Node) []string {
	keys := make([]string, 0)
	for _, node := range nodes {
		for _, key := range node.Keys {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	return keys
}

func addJob(task chan *node.Node, complete chan *node.Node, ring *Ring, nodes []*node.Node) {
	go ring.AddAsync((<-chan *node.Node)(task), (chan<- *node.Node)(complete))
	for _, node := range nodes {
		task <- node
		<-complete
	}
}

func removeJob(task chan *node.Node, complete chan *node.Node, ring *Ring, nodes []*node.Node) {
	go ring.RemoveAsync((<-chan *node.Node)(task), (chan<- *node.Node)(complete))
	for _, node := range nodes {
		task <- node
		<-complete
	}
}

func TestAddNode(t *testing.T) {
	ring := NewRing()
	task, complete := initChannels()

	vnodeNum := 3
	nodes := []*node.Node{
		node.NewNode("127.0.0.1", 0, vnodeNum),
		node.NewNode("127.0.0.1", 1, vnodeNum),
		node.NewNode("127.0.0.1", 2, vnodeNum),
	}

	addJob(task, complete, ring, nodes)

	keys := getKeysShouldBe(nodes)
	keysInTree := getKeysActuallyAre(ring.Tree)

	assert.Equal(t, keys, keysInTree, "Keys should be the same")
}

func TestRemoveNode(t *testing.T) {
	ring := NewRing()

	task := make(chan *node.Node)
	complete := make(chan *node.Node)

	vnodeNum := 3
	nodes := []*node.Node{
		node.NewNode("127.0.0.1", 0, vnodeNum),
		node.NewNode("127.0.0.1", 1, vnodeNum),
		node.NewNode("127.0.0.1", 2, vnodeNum),
	}

	nodesToRemove := nodes[:1]
	nodesAfterRemove := nodes[1:]

	addJob(task, complete, ring, nodes)
	removeJob(task, complete, ring, nodesToRemove)

	keys := getKeysShouldBe(nodesAfterRemove)
	keysInTree := getKeysActuallyAre(ring.Tree)

	assert.Equal(t, keys, keysInTree, "Keys should be the same")
}
