package graph

import (
	"github.com/inpour/algorithms/fundamental"
	"iter"
)

type EulerianStatus int

const (
	NotEulerian      EulerianStatus = iota // no Eulerian path or cycle
	HasEulerianPath                        // Semi-Eulerian
	HasEulerianCycle                       // Eulerian
)

// Eulerian represents a data type for finding an Euler cycle or path in a graph.
// Eulerian Path is a path in a graph that visits every edge exactly once. Eulerian Circuit is an Eulerian Path that
// starts and ends on the same vertex.
// This implementation uses a non-recursive depth-first search.
// This implementation is trickier than the one for digraphs because when we use edge v-w from v's adjacency list,
// we must be careful not to use the second copy of the edge from w's adjacency list.
type Eulerian struct {
	status      EulerianStatus
	pathOrCycle *fundamental.Stack[int]
}

// eulerianEdge helper is an undirected edge, with a field to indicate whether the edge has already been used.
type eulerianEdge struct {
	v      int
	w      int
	isUsed bool
}

func newEulerianEdge(v, w int) *eulerianEdge {
	return &eulerianEdge{
		v:      v,
		w:      w,
		isUsed: false,
	}
}

// otherEdge returns the other vertex of the edge.
func (e *eulerianEdge) otherEdge(vertex int) int {
	if vertex == e.v {
		return e.w
	}
	return e.v
}

// NewEulerian computes an Eulerian path or cycle in the specified graph, if one exists.
// The complexity is O(V + E), where V is the number of vertices and E is the number of edges.
func NewEulerian(graph *Graph) *Eulerian {
	e := &Eulerian{
		status:      HasEulerianCycle,
		pathOrCycle: fundamental.NewStack[int](),
	}

	// If there are no edges in the graph, it is Eulerian (has cycle with length zero)
	if graph.E() == 0 {
		return e
	}

	// find vertex from which to start potential Euler path (a vertex v with odd degree(v) if it exits),
	// if graph has Euler cycle only non isolated vertex is enough
	s := e.nonIsolatedVertex(graph)
	oddDegreeVertices := 0
	for v := 0; v < graph.V(); v++ {
		degree, _ := graph.Degree(v)
		// graph can't have an Euler cycle
		if degree%2 != 0 {
			e.status = HasEulerianPath
			oddDegreeVertices++
			s = v
			// graph can't have an Euler path
			if oddDegreeVertices > 2 {
				e.status = NotEulerian
				return e
			}
		}
	}

	// create local view of adjacency lists, to iterate one vertex at a time
	// the helper eulerianEdge data type is used to avoid exploring both copies of an edge v-w
	adjQueue := make([]*fundamental.Queue[*eulerianEdge], graph.V())
	for v := 0; v < graph.V(); v++ {
		adjQueue[v] = fundamental.NewQueue[*eulerianEdge]()
	}
	for v := 0; v < graph.V(); v++ {
		selfLoopCnt := 0
		adj, _ := graph.Adj(v)
		for w := range adj {
			// careful with self loops
			if v == w {
				if selfLoopCnt%2 == 0 {
					edge := newEulerianEdge(v, w)
					adjQueue[v].Enqueue(edge)
					adjQueue[w].Enqueue(edge)
				}
				selfLoopCnt++
			} else if v < w {
				edge := newEulerianEdge(v, w)
				adjQueue[v].Enqueue(edge)
				adjQueue[w].Enqueue(edge)
			}
		}
	}

	// initialize stack for non-recursive depth-first search (dfs)
	dfsStack := fundamental.NewStack[int]()
	dfsStack.Push(s)
	// greedily search through edges in iterative DFS style
	for !dfsStack.IsEmpty() {
		v, _ := dfsStack.Pop()
		for !adjQueue[v].IsEmpty() {
			edge, _ := adjQueue[v].Dequeue()
			if edge.isUsed {
				continue
			}
			edge.isUsed = true
			dfsStack.Push(v)
			v = edge.otherEdge(v)
		}
		// push vertex with no more leaving edges to cycle
		e.pathOrCycle.Push(v)
	}

	// check if all edges are used
	if e.pathOrCycle.Size() != graph.E()+1 {
		e.status = NotEulerian
	}

	return e
}

// nonIsolatedVertex returns any non-isolated vertex, -1 if no such vertex.
func (e *Eulerian) nonIsolatedVertex(graph *Graph) int {
	for v := 0; v < graph.V(); v++ {
		degree, _ := graph.Degree(v)
		if degree > 0 {
			return v
		}
	}
	return -1
}

// EulerianStatus returns EulerianStatus of the specified graph:
// The complexity is O(1).
func (e *Eulerian) EulerianStatus() EulerianStatus {
	return e.status
}

// PathOrCycle returns the sequence of vertices on an Eulerian path or cycle.
// The complexity is O(1).
func (e *Eulerian) PathOrCycle() iter.Seq[int] {
	return e.pathOrCycle.Iterator()
}
