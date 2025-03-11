package intSet

import (
	"sync"
	"sync/atomic"
)

type IntSet struct {
	data      map[int]struct{}
	dataMutex sync.RWMutex

	/* Item count */
	count int64
}

func New() (s IntSet) {
	s.dataMutex.Lock()
	s.data = make(map[int]struct{})
	s.dataMutex.Unlock()

	return
}

func (s *IntSet) Exists(n int) (ok bool) {

	s.dataMutex.RLock()
	_, ok = s.data[n]
	s.dataMutex.RUnlock()

	return
}

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

func (s *IntSet) Clear() {

	s.dataMutex.Lock()
	clear(s.data)
	s.data = make(map[int]struct{})
	s.dataMutex.Unlock()

	atomic.StoreInt64(&s.count, 0)
}

func (s *IntSet) Count() int64 {
	return atomic.LoadInt64(&s.count)
}
