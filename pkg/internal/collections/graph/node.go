package graph

type connectionDirection int

const (
	OUT  connectionDirection = -1
	IN   connectionDirection = 1
	BOTH connectionDirection = 0
)

func oppositeDirection(dir connectionDirection) connectionDirection {
	// OUT -> IN (-1 * -1 = 1)
	// IN -> OUT (1 * -1 = -1)
	// BOTH -> BOTH (0 * -1 = 0)
	return dir * -1
}

type Node[T any] struct {
	Val         T
	connections map[*Node[T]]connectionDirection
}

func NewNode[T any](v T) *Node[T] {
	return &Node[T]{
		Val:         v,
		connections: make(map[*Node[T]]connectionDirection),
	}
}

func (n Node[T]) GetConnection(node *Node[T]) (connectionDirection, bool) {
	direction, ok := n.connections[node]
	return direction, ok
}

// Deprecated: in favor of GetConnection
func (n Node[T]) IsConnectedWith(node *Node[T]) bool {
	_, ok := n.connections[node]
	return ok
}

// Connect sets a connection between passed node and current node
// setting the direction for each node. When the direction is [BOTH]
func (n *Node[T]) Connect(node *Node[T], dir connectionDirection) {
	oposite := oppositeDirection(dir)
	n.connections[node] = dir
	node.connections[n] = oposite
}

func (n Node[T]) HasConnections() bool {
	return len(n.connections) > 0
}

func (n Node[T]) HasIncomingNodes() bool {
	for _, dir := range n.connections {
		if dir == IN {
			return true
		}
	}

	return false
}

func (n Node[T]) GetIncomingNodes() []*Node[T] {
	nodes := make([]*Node[T], 0)

	for node, dir := range n.connections {
		if dir == IN {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (n Node[T]) HasOutgoingNodes() bool {
	for _, dir := range n.connections {
		if dir == OUT {
			return true
		}
	}

	return false
}

func (n Node[T]) GetOutgoingNodes() []*Node[T] {
	nodes := make([]*Node[T], 0)

	for node, dir := range n.connections {
		if dir == OUT {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (n *Node[T]) Disconnect(node *Node[T]) {
	delete(n.connections, node)
	delete(node.connections, n)
}
