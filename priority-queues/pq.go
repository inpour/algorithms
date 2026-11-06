package priority_queues

import (
	"errors"
)

var ErrEmptyPQ = errors.New("priority queue is empty")
var ErrIndexNotExists = errors.New("index is not in the priority queue")
var ErrIndexOutOfRange = errors.New("index out of range")
var ErrIndexAlreadyExists = errors.New("index is already in the priority queue")
var ErrDecreaseWithInvalidKey = errors.New("calling DecreaseKey() with a key equal or greater that the key in the priority queue")
var ErrIncreaseWithInvalidKey = errors.New("calling IncreaseKey() with a key equal or less than the key in the priority queue")
