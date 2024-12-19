package search

import (
	"errors"
	"iter"
)

type SymbolTable[K, V any] interface {
	Put(key K, val V)          // put key-value pair into the table (overwrite value if key already exists)
	Get(key K) (V, error)      // get value paired with key, ErrAbsentKey error if key is absent
	Delete(key K) error        // remove key (and its value) from table, ErrAbsentKey error if key is absent
	Contains(key K) bool       // is there a value paired with key?
	IsEmpty() bool             // is the table empty?
	Size() int                 // number of key-value pairs in the table
	Iterator() iter.Seq2[K, V] // returns an iterator that iterates over all the key-value pairs in the table
}

type OrderedSymbolTable[K, V any] interface {
	SymbolTable[K, V]
	Min() K                                 // smallest key
	max() K                                 // largest key
	Floor(key K) (K, error)                 // largest key less than or equal to key, ErrAbsentKey error if key is absent
	Ceiling(key K) (K, error)               // smallest key greater than or equal to key, ErrAbsentKey error if key is absent
	Rank(key K) (int, error)                // number of keys less than key, ErrAbsentKey error if key is absent
	Select(k int) (K, error)                // key of rank k, ErrInvalidRank error if rank is out of range
	DelMin()                                // delete smallest key
	DelMax()                                // delete largest key
	RangeSize(lo, hi K) int                 // number of keys in [lo:hi]
	RangeIterator(lo, hi K) iter.Seq2[K, V] // key-value pairs where keys in [lo:hi], in sorted order
}

var ErrAbsentKey = errors.New("key is absent")
var ErrInvalidRank = errors.New("rank is out of range")
