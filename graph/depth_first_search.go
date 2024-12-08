package graph

// DepthFirstSearch represents a data type for determining the vertices connected to a given source vertex (s) in a graph.
// It uses Θ(v) extra space (not including the graph).
type DepthFirstSearch struct {
	marked []bool
	count  int
}

// NewDepthFirstSearch computes the vertices in graph that are connected to the source vertex (s).
// It takes Θ(v+e) time in the worst case, where v is the number of vertices and e is the number of edges.
func NewDepthFirstSearch(graph *Graph, s int) (*DepthFirstSearch, error) {
	if err := graph.validateVertex(s); err != nil {
		return nil, err
	}
	d := &DepthFirstSearch{
		marked: make([]bool, graph.V()),
		count:  0,
	}
	d.dfs(graph, s)
	return d, nil
}

// dfs (depth first search) from v
func (d *DepthFirstSearch) dfs(graph *Graph, v int) {
	d.count++
	d.marked[v] = true
	adj, _ := graph.Adj(v)
	for w := range adj {
		if !d.marked[w] {
			d.dfs(graph, w)
		}
	}
}

// Marked returns true if there is a path between the source vertex (s) and vertex v.
// The complexity is O(1).
func (d *DepthFirstSearch) Marked(v int) (bool, error) {
	if v < 0 || v >= len(d.marked) {
		return false, ErrInvalidVertexIndex
	}
	return d.marked[v], nil
}

// Count returns the number of vertices connected to the source vertex (s).
// The complexity is O(1).
func (d *DepthFirstSearch) Count() int {
	return d.count
}
