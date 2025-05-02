package intSet

import (
	"sync"
	"sync/atomic"
)

// Implements a thread-safe set of integers
type IntSet struct {
	data      map[int]struct{}
	dataMutex sync.RWMutex

	// Item count
	count int64
}

// Initialize a new IntSet
func New() (s IntSet) {

	s.dataMutex.Lock()
	s.data = make(map[int]struct{})
	s.dataMutex.Unlock()

	return
}

// Checks if integer 'n' exists in the set
func (s *IntSet) Exists(n int) bool {

	s.dataMutex.RLock()
	_, ok := s.data[n]
	s.dataMutex.RUnlock()

	return ok
}

// Inserts 'n' in the set if it's not already, returns true if n is new in the Set
func (s *IntSet) Insert(n int) bool {

	if s.Exists(n) {

		return false
	}

	s.dataMutex.Lock()
	s.data[n] = struct{}{}
	s.dataMutex.Unlock()

	atomic.AddInt64(&s.count, 1)

	return true
}

// Clears the set by literally freeing the data and reallocating
func (s *IntSet) Clear() {

	s.dataMutex.Lock()
	clear(s.data)
	s.data = make(map[int]struct{})
	s.dataMutex.Unlock()

	atomic.StoreInt64(&s.count, 0)
}

// Returns the number of elements in the set
func (s *IntSet) Count() int64 {

	return atomic.LoadInt64(&s.count)
}
