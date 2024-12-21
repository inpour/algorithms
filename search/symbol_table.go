package search

import (
	"errors"
	"iter"
)

type SymbolTable[K, V any] interface {
	Put(key K, val V)          // put key-value pair into the table (overwrite value if key already exists)
	Get(key K) (V, error)      // get value paired with key, ErrAbsentKey if key is absent
	Delete(key K) error        // remove key and associated value, ErrAbsentKey if key is absent
	Contains(key K) bool       // is there a value paired with key?
	IsEmpty() bool             // is the table empty?
	Size() int                 // number of key-value pairs in the table
	Iterator() iter.Seq2[K, V] // returns an iterator that iterates over all the key-value pairs in the table
}

type OrderedSymbolTable[K, V any] interface {
	SymbolTable[K, V]
	Min() (K, error)                        // smallest key, ErrEmptySymbolTable if symbol table is empty
	Max() (K, error)                        // largest key, ErrEmptySymbolTable if symbol table is empty
	Floor(key K) (K, error)                 // largest key less than or equal to key, ErrTooSmallFloorKey if key to floor is too small
	Ceiling(key K) (K, error)               // smallest key greater than or equal to key, ErrTooLargeCeilingKey if key to ceiling is too large
	Rank(key K) (int, error)                // number of keys less than key, ErrAbsentKey if key is absent
	Select(k int) (K, error)                // key of rank k, ErrInvalidRank if rank is out of range
	DelMin() error                          // delete the smallest key and associated value, ErrEmptySymbolTable if symbol table is empty
	DelMax() error                          // delete the largest key and associated value, ErrEmptySymbolTable if symbol table is empty
	RangeSize(lo, hi K) int                 // number of keys in [lo:hi] range
	RangeIterator(lo, hi K) iter.Seq2[K, V] // key-value pairs where keys in [lo:hi] range, in sorted order
}

var ErrAbsentKey = errors.New("key is absent")
var ErrInvalidRank = errors.New("rank is out of range")
var ErrEmptySymbolTable = errors.New("symbol table is empty")
var ErrTooSmallFloorKey = errors.New("key to floor is too small")
var ErrTooLargeCeilingKey = errors.New("key to ceiling is too large")
