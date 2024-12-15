package graph

// ConnectedComponents represents a data type for determining the connected components in a graph.
// This implementation uses depth-first search.
// The component identifier (id) of a vertex is an integer between 0 and k–1, where k is the number
// of connected components. Two vertices have the same component identifier if and only if they are
// in the same connected component.
type ConnectedComponents struct {
	marked []bool // marked[v] = has vertex v been marked?
	id     []int  // id[v] = id of connected component containing v
	size   []int  // size[id] = number of vertices in given component
	count  int    // number of connected components
}

// NewConnectedComponents computes the connected components of the graph
// It uses Θ(v) extra space (not including the graph) and takes Θ(v+e) time, where v is the number
// of vertices and e is the number of edges.
func NewConnectedComponents(graph *Graph) *ConnectedComponents {
	c := &ConnectedComponents{
		marked: make([]bool, graph.V()),
		id:     make([]int, graph.V()),
		size:   make([]int, graph.V()),
		count:  0,
	}
	for v := 0; v < graph.V(); v++ {
		if !c.marked[v] {
			c.dfs(graph, v)
			c.count++
		}
	}
	return c
}

// dfs (depth first search) from v
func (c *ConnectedComponents) dfs(graph *Graph, v int) {
	c.marked[v] = true
	c.id[v] = c.count
	c.size[c.count]++
	adj, _ := graph.Adj(v)
	for w := range adj {
		if !c.marked[w] {
			c.dfs(graph, w)
		}
	}
}

// ID returns the component id of the connected component containing vertex v.
// The complexity is O(1).
func (c *ConnectedComponents) ID(v int) (int, error) {
	if err := c.validateVertex(v); err != nil {
		return 0, err
	}
	return c.id[v], nil
}

// Size returns the number of vertices in the connected component containing vertex v.
// The complexity is O(1).
func (c *ConnectedComponents) Size(v int) (int, error) {
	if err := c.validateVertex(v); err != nil {
		return 0, err
	}
	return c.size[c.id[v]], nil
}

// Count returns the number of connected components in the graph.
// The complexity is O(1).
func (c *ConnectedComponents) Count() int {
	return c.count
}

// Connected returns true if vertices v and w are in the same connected component.
// The complexity is O(1).
func (c *ConnectedComponents) Connected(v, w int) (bool, error) {
	if err := c.validateVertex(v); err != nil {
		return false, err
	}
	if err := c.validateVertex(w); err != nil {
		return false, err
	}
	return c.id[v] == c.id[w], nil
}

func (c *ConnectedComponents) validateVertex(v int) error {
	if v < 0 || v >= len(c.marked) {
		return ErrInvalidVertexIndex
	}
	return nil
}
