package graph

// DepthFirstSearch (DFS) represents a data type for determining single-source or multiple-source reachability in a graph.
// It uses O(V) extra space (not including the graph), where V is the number of vertices.
type DepthFirstSearch struct {
	marked []bool // marked[v] = is there an s-v path?
	count  int    // number of vertices connected to s
}

// NewDepthFirstSearch computes the vertices in graph that are connected to the source vertex (s).
// The complexity is O(V + E), where V is the number of vertices and E is the number of edges.
func NewDepthFirstSearch(graph UndirectedOrDirectedGraph, s int) (*DepthFirstSearch, error) {
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

// NewDepthFirstSearchMultiSource computes the vertices in graph that are connected to any of the source vertices (sources).
// The complexity is O(V + E), where V is the number of vertices and E is the number of edges.
func NewDepthFirstSearchMultiSource(graph UndirectedOrDirectedGraph, sources []int) (*DepthFirstSearch, error) {
	for _, s := range sources {
		if err := graph.validateVertex(s); err != nil {
			return nil, err
		}
	}
	d := &DepthFirstSearch{
		marked: make([]bool, graph.V()),
		count:  0,
	}
	for _, s := range sources {
		if !d.marked[s] {
			d.dfs(graph, s)
		}
	}
	return d, nil
}

// dfs (depth first search) from v
func (d *DepthFirstSearch) dfs(graph UndirectedOrDirectedGraph, v int) {
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
	if err := d.validateVertex(v); err != nil {
		return false, err
	}
	return d.marked[v], nil
}

// Count returns the number of vertices reachable from the source vertex or source vertices.
// The complexity is O(1).
func (d *DepthFirstSearch) Count() int {
	return d.count
}

func (d *DepthFirstSearch) validateVertex(v int) error {
	if v < 0 || v >= len(d.marked) {
		return ErrInvalidVertexIndex
	}
	return nil
}
