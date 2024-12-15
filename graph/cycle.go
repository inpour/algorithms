package graph

import (
	"github.com/inpour/algorithms/fundamental"
	"iter"
)

// Cycle represents a data type for determining whether an undirected graph has a simple cycle.
// The HasCycle operation determines whether the graph has a cycle and, if so, the Cycle operation returns one.
// This implementation uses depth-first search (DFS).
type Cycle struct {
	marked []bool
	edgeTo []int
	cycle  *fundamental.Stack[int]
}

// NewCycle determines whether the undirected graph has a cycle and, if so, finds such a cycle.
// It uses Θ(v) extra space (not including the graph).
// It takes Θ(v+e) time in the worst case, where v is the number of vertices and e is the number of edges.
// The depth-first search part takes only O(v) time; however, checking for parallel edges takes Θ(v+e) time
// in the worst case.
func NewCycle(graph *Graph) *Cycle {
	c := &Cycle{
		marked: make([]bool, graph.V()),
		edgeTo: make([]int, graph.V()),
		cycle:  fundamental.NewStack[int](),
	}

	//if c.hasSelfLoop(graph) {
	//	return c
	//}

	if c.hasParallelEdges(graph) {
		return c
	}

	for v := 0; v < graph.V(); v++ {
		if !c.marked[v] {
			c.dfs(graph, v, -1)
		}
	}
	return c
}

// hasParallelEdges returns true if the graph have two parallel edges
// side effect: fill cycle to be two parallel edges
func (c *Cycle) hasParallelEdges(graph *Graph) bool {
	for v := 0; v < graph.V(); v++ {

		// check for parallel edges incident to v
		adj, _ := graph.Adj(v)
		for w := range adj {
			if c.marked[w] {
				c.cycle.Push(v)
				c.cycle.Push(w)
				c.cycle.Push(v)
				return true
			}
			c.marked[w] = true
		}

		// reset so marked[v] = false for all v
		adj, _ = graph.Adj(v)
		for w := range adj {
			c.marked[w] = false
		}
	}
	return false
}

// hasSelfLoop returns true if the graph have a self loop
// side effect: fill cycle to be self loop
func (c *Cycle) hasSelfLoop(graph *Graph) bool {
	for v := 0; v < graph.V(); v++ {
		adj, _ := graph.Adj(v)
		for w := range adj {
			if v == w {
				c.cycle.Push(v)
				c.cycle.Push(v)
				return true
			}
		}
	}
	return false
}

// dfs (depth first search)
func (c *Cycle) dfs(graph *Graph, v int, parent int) {
	c.marked[v] = true
	adj, _ := graph.Adj(v)
	for w := range adj {

		// short circuit if cycle already found
		if !c.cycle.IsEmpty() {
			return
		}

		// check for cycle but disregard parent of current vertex
		if !c.marked[w] {
			c.edgeTo[w] = v
			c.dfs(graph, w, v)
		} else if w != parent {
			for x := v; x != w; x = c.edgeTo[x] {
				c.cycle.Push(x)
			}
			c.cycle.Push(w)
			c.cycle.Push(v)
		}
	}
}

// HasCycle returns true if the graph has a cycle.
func (c *Cycle) HasCycle() bool {
	return !c.cycle.IsEmpty()
}

// Cycle returns a cycle in the graph.
func (c *Cycle) Cycle() iter.Seq[int] {
	return c.cycle.Iterator()
}
