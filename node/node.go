package node

type Node struct {
	//each Node is globally identified by ip:port
	//the 'ip:port' string is the uuid of this Node
	ip   string
	port int
	uuid string

	//key is the position on consistent hashing,
	//because key is usually extremely large(e.g. 2^160) to reduce confliction
	//it should be a string
	key string

	//when the membership protocol knows consistent hashing has
	//saved this Node, then it's marked as saved
	saved bool
}
