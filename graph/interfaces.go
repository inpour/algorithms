package graph

import (
	"errors"
	"iter"
)

type UndirectedOrDirectedGraph interface {
	V() int                           // returns the number of vertices
	E() int                           // returns the number of edges
	validateVertex(v int) error       // validate given vertex index (v)
	AddEdge(v, w int) error           // adds the undirected or directed edge v-w
	Adj(v int) (iter.Seq[int], error) // returns an iterator that iterates over vertices adjacent to vertex v
}

var ErrInvalidVertices = errors.New("number of vertices in a Graph must be non-negative")
var ErrInvalidVertexIndex = errors.New("invalid vertex index")
