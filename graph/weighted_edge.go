package graph

import (
	"cmp"
	"errors"
	"math"
)

// WeightedEdge class represents a weighted edge in an EdgeWeightedGraph.
// Each edge consists of two integers (naming the two vertices) and a real-value weight.
// The data type provides methods for accessing the two endpoints of the edge and the weight.
// The natural order for this data type is by ascending order of weight.
type WeightedEdge struct {
	v      int
	w      int
	weight float64
}

// NewWeightedEdge Initializes an edge between vertices v and w of the given weight.
func NewWeightedEdge(v, w int, weight float64) (*WeightedEdge, error) {
	if v < 0 || w < 0 {
		return nil, errors.New("vertex index must be a non-negative integer")
	}
	if math.IsNaN(weight) {
		return nil, errors.New("weight is NaN")
	}

	return &WeightedEdge{
		v:      v,
		w:      w,
		weight: weight,
	}, nil
}

// Weight returns the weight of this edge.
func (e *WeightedEdge) Weight() float64 {
	return e.weight
}

// Either returns on of endpoints of this edge.
func (e *WeightedEdge) Either() int {
	return e.v
}

// Other returns on of endpoints of this edge.
func (e *WeightedEdge) Other(vertex int) (int, error) {
	if vertex == e.v {
		return e.w, nil
	} else if vertex == e.w {
		return e.v, nil
	}

	return vertex, errors.New("illegal endpoint")
}

// CompareTo return a negative integer, zero, or positive integer depending on whether
// the weight of this is less than, equal to, or greater than the
// argument edge
func (e *WeightedEdge) CompareTo(that *WeightedEdge) int {
	return cmp.Compare(e.weight, that.Weight())
}
