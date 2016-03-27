package consistentHashing

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/otnt/distributed-system-how-to/node"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

//channels are used when tesing async methods
func createChannels() (task chan *node.Node, complete chan *node.Node) {
	task = make(chan *node.Node)
	complete = make(chan *node.Node)
	return
}

//create three nodes for testing
//keys are
//[42f91ea3437c635edb4b30d3628141e3d6a0ebc7 d9443c1098545a2f7388ae40e3aade9828dc588e b5f7f3116902fa6f4a5819be2015b7972cb93e0a]
//[97c1af2272de15532b1483651b715129332f8406 7260a48008fb01d884067d8e50b64ac56b9c3221 eb102fa9386db4715c2cfc93d019ca21c194b767]
//[bd5b206633d9f79501860b0c03559379b89435ff c98a2ee624a8306f46ef5e01f6ea5dcce0b7ac52 98273171fa35b563b7519cef47d48005ee391f2c]
func createNodes() []*node.Node {
	vnodeNum := 3
	nodes := []*node.Node{
		node.NewNode("127.0.0.1", 0, vnodeNum),
		node.NewNode("127.0.0.1", 1, vnodeNum),
		node.NewNode("127.0.0.1", 2, vnodeNum),
	}

	return nodes
}

//after modification, get the keys that are actually
//int redblack tree
func getKeysActuallyAre(tree *rbt.Tree) []string {
	interfaceKeysInTree := tree.Keys()
	keysInTree := make([]string, 0)
	for _, key := range interfaceKeysInTree {
		keysInTree = append(keysInTree, key.(string))
	}

	return keysInTree
}

//use sorting method to create expected keys
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

//add nodes to tree in async method
func addNodesToTreeAsync(task chan *node.Node, complete chan *node.Node, ring *Ring, nodes []*node.Node) {
	go ring.AddAsync((<-chan *node.Node)(task), (chan<- *node.Node)(complete))
	for _, node := range nodes {
		task <- node
		<-complete
	}
}

//remove nodes from tree in async method
func removeNodesFromTreeAsync(task chan *node.Node, complete chan *node.Node, ring *Ring, nodes []*node.Node) {
	go ring.RemoveAsync((<-chan *node.Node)(task), (chan<- *node.Node)(complete))
	for _, node := range nodes {
		task <- node
		<-complete
	}
}

//add nodes to tree in sync method
func addNodesToTree(ring *Ring, nodes []*node.Node) {
	for _, node := range nodes {
		ring.AddSync(node)
	}
}

//remove nodes from tree in sync method
func removeNodesFromTree(ring *Ring, nodes []*node.Node) {
	for _, node := range nodes {
		ring.RemoveSync(node)
	}
}

func TestAddNodeAsync(t *testing.T) {
	ring := NewRing()
	task, complete := createChannels()
	nodes := createNodes()

	addNodesToTreeAsync(task, complete, ring, nodes)

	keysShouldBe := getKeysShouldBe(nodes)
	keysInTree := getKeysActuallyAre(ring.Tree)

	assert.Equal(t, keysShouldBe, keysInTree, "Keys should be the same")
}

func TestRemoveNodeAsync(t *testing.T) {
	ring := NewRing()
	task, complete := createChannels()
	nodes := createNodes()

	nodesToRemove := nodes[:1]
	nodesAfterRemove := nodes[1:]

	addNodesToTreeAsync(task, complete, ring, nodes)
	removeNodesFromTreeAsync(task, complete, ring, nodesToRemove)

	keysShouldBe := getKeysShouldBe(nodesAfterRemove)
	keysInTree := getKeysActuallyAre(ring.Tree)

	assert.Equal(t, keysShouldBe, keysInTree, "Keys should be the same")
}

func TestAddNodeSync(t *testing.T) {
	ring := NewRing()
	nodes := createNodes()

	addNodesToTree(ring, nodes)

	keysShouldBe := getKeysShouldBe(nodes)
	keysInTree := getKeysActuallyAre(ring.Tree)

	assert.Equal(t, keysShouldBe, keysInTree, "Keys should be the same")
}

func TestRemoveNodeSync(t *testing.T) {
	ring := NewRing()
	nodes := createNodes()

	nodesToRemove := nodes[:1]
	nodesAfterRemove := nodes[1:]

	addNodesToTree(ring, nodes)
	removeNodesFromTree(ring, nodesToRemove)

	keysShouldBe := getKeysShouldBe(nodesAfterRemove)
	keysInTree := getKeysActuallyAre(ring.Tree)

	assert.Equal(t, keysShouldBe, keysInTree, "Keys should be the same")
}

func contains(list []string, element string) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}

//Looking up a key will give a Node whose
//Keys contain this key
func TestLookUp1(t *testing.T) {
	ring := NewRing()
	nodes := createNodes()
	addNodesToTree(ring, nodes)
	keysShouldBe := getKeysShouldBe(nodes)

	for _, key := range keysShouldBe {
		node, err := ring.LookUp(SubOne(key))

		assert.Nil(t, err)
		assert.True(t, contains(node.Keys, key))
	}
}

//Looking for a key that is same as the node,
//and a key that is smaller than this node by one,
//and a key that is larger than this node by one,
//the first two Nodes should be the same,
//the third Node should contains the next key
func TestLookUp2(t *testing.T) {
	ring := NewRing()
	nodes := createNodes()
	addNodesToTree(ring, nodes)
	keysShouldBe := getKeysShouldBe(nodes)

	for _, key := range keysShouldBe {
		//lookup a key smaller than an node by one
		node0, err0 := ring.LookUp(SubOne(key))
		//lookup a key at same position as node
		node1, err1 := ring.LookUp(key)
		//lookup a key larger than an node by one
		node2, err2 := ring.LookUp(AddOne(key))

		assert.Nil(t, err0)
		assert.Nil(t, err1)
		assert.Nil(t, err2)

		assert.Equal(t, node0, node1)

		index := sort.SearchStrings(keysShouldBe, key)
		index = (index + 1) % len(keysShouldBe)
		assert.True(t, contains(node2.Keys, keysShouldBe[index]))
	}
}

//NextOf(key) is equivallent to
//LookUp(newKey) where newKey = key+1
func TestNextOf(t *testing.T) {
	ring := NewRing()
	nodes := createNodes()
	addNodesToTree(ring, nodes)
	keysShouldBe := getKeysShouldBe(nodes)

	for _, key := range keysShouldBe {
		node, _ := ring.LookUp(AddOne(key))
		node2, _ := ring.NextOf(key)
		assert.Equal(t, node, node2)
	}
}
