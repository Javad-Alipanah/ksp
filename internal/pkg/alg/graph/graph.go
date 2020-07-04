package graph

type Node struct {
	added bool
	discovered bool
	parent int
}

type Graph struct {
	nodes []Node
	edges map[int][]int
}

// returns a graph with `size` nodes without any edges
func New(size int) *Graph {
	return &Graph{
		nodes: make([]Node, size),
		edges: make(map[int][]int, 0),
	}
}

func (g *Graph) size() int {
	return len(g.nodes)
}

func (g *Graph) AddNode(node int) {
	if node >= g.size() || node < 0 {
		panic("invalid node to add to graph")
	}

	g.nodes[node].added = true
}

func (g *Graph) HasBeenAdded(node int) bool {
	return g.nodes[node].added
}

// Connect adds an edge between src and dst
func (g *Graph) Connect(src, dst int)  {
	if src >= g.size() || dst >= g.size() || src < 0 || dst < 0 {
		panic("invalid src or dst for edge")
	}

	if g.edges[src] == nil {
		g.edges[src] = make([]int, 0)
	}

	g.edges[src] = append(g.edges[src], dst)
}

func (g *Graph) buildPath(src, dst int) []int {
	curr := dst
	path := make([]int, 1)
	path[0] = curr
	for curr != src {
		curr = g.nodes[curr].parent
		path = append(path, curr)
	}

	reverse := make([]int, len(path))
	for i := len(path); i > 0; i-- {
		reverse[len(path) - i] = path[i - 1]
	}
	return reverse
}

type queue struct {
	elems []int
}

func newQueue() *queue {
	return &queue{
		elems: make([]int, 0),
	}
}

func (q *queue) enqueue(elem int) {
	q.elems = append(q.elems, elem)
}

func (q *queue) dequeue() int {
	elem := q.elems[0]
	q.elems = q.elems[1:]
	return elem
}

func (q *queue) isEmpty() bool {
	return len(q.elems) == 0
}

func (g *Graph) BFS(src, dst int) []int {
	q := newQueue()
	g.nodes[src].discovered = true
	q.enqueue(src)
	for !q.isEmpty() {
		v := q.dequeue()
		if v == dst {
			path := make([]int, 0)
			return append(path, src)
		}
		for _, w := range g.edges[v] {
			if !g.nodes[w].discovered {
				g.nodes[w].discovered = true
				g.nodes[w].parent = v
				if w == dst {
					return g.buildPath(src, w)
				}
				q.enqueue(w)
			}
		}
	}
	return nil
}