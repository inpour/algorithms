package graph

import (
	"github.com/inpour/algorithms/fundamental"
	"iter"
)

// DirectedEulerian represents a data type for finding an Eulerian cycle or path in a digraph.
// Eulerian Path is a path in a digraph that visits every edge exactly once. Eulerian cycle is an Eulerian Path that
// starts and ends on the same vertex.
// This implementation uses a non-recursive depth-first search.
// It uses O(V + E) extra space (not including the graph), where V is the number of vertices and E is the number of edges.
type DirectedEulerian struct {
	status      EulerianStatus
	pathOrCycle *fundamental.Stack[int]
}

// NewDirectedEulerian computes an Eulerian path or cycle in the specified digraph, if one exists.
// The complexity is O(V + E), where V is the number of vertices and E is the number of edges.
func NewDirectedEulerian(digraph *Digraph) *DirectedEulerian {
	e := &DirectedEulerian{
		status:      HasEulerianCycle,
		pathOrCycle: fundamental.NewStack[int](),
	}

	// If there are no edges in the digraph, it is Eulerian (has cycle with length zero)
	if digraph.E() == 0 {
		return e
	}

	// find vertex from which to start potential Eulerian path (a vertex v with outdegree(v) > indegree(v) if it exits),
	// if digraph has Eulerian cycle only non isolated vertex is enough
	s := e.nonIsolatedVertex(digraph)
	deficit := 0
	for v := 0; v < digraph.V(); v++ {
		outDegree, _ := digraph.OutDegree(v)
		inDegree, _ := digraph.InDegree(v)
		// digraph can't have an Eulerian cycle
		if outDegree > inDegree {
			e.status = HasEulerianPath
			deficit += outDegree - inDegree
			s = v
			// digraph can't have an Eulerian path
			if deficit > 1 {
				e.status = NotEulerian
				return e
			}
		}
	}

	// create local view of adjacency lists, to iterate one vertex at a time
	adjQueue := make([]*fundamental.Queue[int], digraph.V())
	for v := 0; v < digraph.V(); v++ {
		adjQueue[v] = fundamental.NewQueue[int]()
		adj, _ := digraph.Adj(v)
		for w := range adj {
			adjQueue[v].Enqueue(w)
		}
	}

	// initialize stack for non-recursive depth-first search (dfs)
	dfsStack := fundamental.NewStack[int]()
	dfsStack.Push(s)
	// greedily search through edges in iterative DFS style
	for !dfsStack.IsEmpty() {
		v, _ := dfsStack.Pop()
		for !adjQueue[v].IsEmpty() {
			dfsStack.Push(v)
			v, _ = adjQueue[v].Dequeue()
		}
		// push vertex with no more leaving edges
		e.pathOrCycle.Push(v)
	}

	// check if all edges are used
	if e.pathOrCycle.Size() != digraph.E()+1 {
		e.status = NotEulerian
	}

	return e
}

// nonIsolatedVertex returns any non-isolated vertex, -1 if no such vertex.
func (e *DirectedEulerian) nonIsolatedVertex(digraph *Digraph) int {
	for v := 0; v < digraph.V(); v++ {
		outDegree, _ := digraph.OutDegree(v)
		if outDegree > 0 {
			return v
		}
	}
	return -1
}

// EulerianStatus returns EulerianStatus of the specified digraph:
// The complexity is O(1).
func (e *DirectedEulerian) EulerianStatus() EulerianStatus {
	return e.status
}

// PathOrCycle returns the sequence of vertices on an Eulerian path or cycle.
// The complexity is O(1).
func (e *DirectedEulerian) PathOrCycle() iter.Seq[int] {
	return e.pathOrCycle.Iterator()
}
