package graph

// SymbolDigraph represents a directed graph, where the vertex names are arbitrary strings.
// By providing mappings between string vertex names and integers, it serves as a wrapper around the Graph,
// which assumes the vertex names are integers between 0 and v - 1.
// This implementation uses a symbol table (st) to map from strings to integers, an array to map from integers
// to strings, and a Graph to store the underlying graph.
type SymbolDigraph struct {
	st      map[string]int // a symbol table that maps names to indices
	keys    []string       // a slice that maps indices to names
	digraph *Digraph       // the underlying digraph
}

// NewSymbolDigraph initializes a SymbolDigraph from vertexNames.
// The complexity is O(V), where V is the number of vertices (length of vertexNames).
func NewSymbolDigraph(vertexNames []string) *SymbolDigraph {
	v := len(vertexNames)
	st := make(map[string]int, v)
	keys := make([]string, v)
	digraph, _ := NewDigraph(v)

	for i, name := range vertexNames {
		st[name] = i
		keys[i] = name
	}

	return &SymbolDigraph{
		st:      st,
		keys:    keys,
		digraph: digraph,
	}
}

// Contains returns true if the graph contain the vertex name.
// The complexity is O(1).
func (s *SymbolDigraph) Contains(name string) bool {
	_, ok := s.st[name]
	return ok
}

// IndexOf returns the integer associated with the vertex name.
// The complexity is O(1).
func (s *SymbolDigraph) IndexOf(name string) (int, error) {
	index, ok := s.st[name]
	if !ok {
		return index, ErrInvalidName
	}
	return index, nil
}

// NameOf returns the name of the vertex associated with the integer v
// The complexity is O(1).
func (s *SymbolDigraph) NameOf(v int) (string, error) {
	if err := s.digraph.validateVertex(v); err != nil {
		var name string
		return name, err
	}
	return s.keys[v], nil
}

// Digraph returns the digraph associated with the symbol graph.
// The complexity is O(1).
func (s *SymbolDigraph) Digraph() *Digraph {
	return s.digraph
}

// AddEdge adds the edge v-w.
// The complexity is O(1).
func (s *SymbolDigraph) AddEdge(v, w string) error {
	vi, err := s.IndexOf(v)
	if err != nil {
		return err
	}
	wi, err := s.IndexOf(w)
	if err != nil {
		return err
	}
	s.digraph.AddEdge(vi, wi)
	return nil
}
