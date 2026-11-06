package priority_queues

import (
	"errors"
)

var ErrEmptyPQ = errors.New("priority queue is empty")
var ErrInexNotExists = errors.New("index is not in the priority queue")
