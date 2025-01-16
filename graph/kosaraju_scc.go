package graph

// KosarajuSCC (Kosaraju-Sharir strongly connected components) represents a data type for determining the
// strongly connected components (or strong components for short) in a digraph.
// This implementation uses the Kosaraju-Sharir algorithm.
// The component identifier (id) of a vertex is an integer between 0 and k–1, where k is the number
// of strong components. Two vertices have the same component identifier if and only if they are
// in the same strong component.
// It uses O(V) extra space (not including the graph), where V is the number of vertices.
type KosarajuSCC struct {
	marked []bool // marked[v] = has vertex v been marked?
	id     []int  // id[v] = id of strong component containing v
	size   []int  // size[id] = number of vertices in given strong component
	count  int    // number of strong components
}

// NewKosarajuSCC computes the strong components of the digraph.
// The complexity is O(V + E), where V is the number of vertices and E is the number of edges.
func NewKosarajuSCC(digraph *Digraph) *KosarajuSCC {
	k := &KosarajuSCC{
		marked: make([]bool, digraph.V()),
		id:     make([]int, digraph.V()),
		size:   make([]int, digraph.V()),
		count:  0,
	}
	dfs := NewDepthFirstOrder(digraph.Reverse())
	for v := range dfs.ReversePost() {
		if !k.marked[v] {
			k.dfs(digraph, v)
			k.count++
		}
	}
	return k
}

// dfs (depth first search) from v
func (k *KosarajuSCC) dfs(digraph *Digraph, v int) {
	k.marked[v] = true
	k.id[v] = k.count
	k.size[k.count]++
	adj, _ := digraph.Adj(v)
	for w := range adj {
		if !k.marked[w] {
			k.dfs(digraph, w)
		}
	}
}

// ID returns the component id of the strong component containing vertex v.
// The complexity is O(1).
func (k *KosarajuSCC) ID(v int) (int, error) {
	if err := k.validateVertex(v); err != nil {
		return 0, err
	}
	return k.id[v], nil
}

// Size returns the number of vertices in the strong component containing vertex v.
// The complexity is O(1).
func (k *KosarajuSCC) Size(v int) (int, error) {
	if err := k.validateVertex(v); err != nil {
		return 0, err
	}
	return k.size[k.id[v]], nil
}

// Count returns the number of strong components.
// The complexity is O(1).
func (k *KosarajuSCC) Count() int {
	return k.count
}

// StronglyConnected returns true if vertices v and w are in the same strong component.
// The complexity is O(1).
func (k *KosarajuSCC) StronglyConnected(v, w int) (bool, error) {
	if err := k.validateVertex(v); err != nil {
		return false, err
	}
	if err := k.validateVertex(w); err != nil {
		return false, err
	}
	return k.id[v] == k.id[w], nil
}

func (k *KosarajuSCC) validateVertex(v int) error {
	if v < 0 || v >= len(k.marked) {
		return ErrInvalidVertexIndex
	}
	return nil
}
