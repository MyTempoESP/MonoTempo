package main

import (
	"aa2/intSet"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// smoothWindow is the range over which each tag's display time is randomly
	// spread.  A burst of N tags arriving at once will be committed to the
	// display counters uniformly across this duration, with no jumps or
	// stutters.
	smoothWindow = 1500 * time.Millisecond

	// smoothTickInterval is the drain goroutine's wake interval.
	smoothTickInterval = 50 * time.Millisecond
)

// smoothEntry is one buffered tag waiting to be committed to display counters.
type smoothEntry struct {
	epc       int
	releaseAt time.Time
}

// TagSmoother buffers incoming RFID tag events and releases them to the
// caller-owned PCData counters and intSets at randomly jittered times within
// [now, now+smoothWindow], so the display counter increments smoothly instead
// of jumping in chunks.
//
// Antenna recording is NOT handled here; callers must call
// pcData.Antennas[...].Record() directly.
type TagSmoother struct {
	mu      sync.Mutex
	pending []smoothEntry

	pcTags            *atomic.Int64
	pcUniqueTags      *atomic.Int32
	pcPermanentUnique *atomic.Int32
	tagSet            *intSet.IntSet
	permanentSet      *intSet.IntSet
}

// newTagSmoother constructs a TagSmoother and starts its drain goroutine.
func newTagSmoother(
	pcTags *atomic.Int64,
	pcUniqueTags *atomic.Int32,
	pcPermanentUnique *atomic.Int32,
	tagSet *intSet.IntSet,
	permanentSet *intSet.IntSet,
) *TagSmoother {
	s := &TagSmoother{
		pcTags:            pcTags,
		pcUniqueTags:      pcUniqueTags,
		pcPermanentUnique: pcPermanentUnique,
		tagSet:            tagSet,
		permanentSet:      permanentSet,
		pending:           make([]smoothEntry, 0, 64),
	}
	go s.run()
	return s
}

// Push enqueues one tag for smoothed delivery.  The tag will be committed to
// the display counters at a uniformly-random time within [now, now+smoothWindow].
func (s *TagSmoother) Push(epc int) {
	releaseAt := time.Now().Add(time.Duration(rand.Int63n(int64(smoothWindow))))
	s.mu.Lock()
	s.pending = append(s.pending, smoothEntry{epc: epc, releaseAt: releaseAt})
	s.mu.Unlock()
}

// Clear discards all buffered tags and resets all counters. Called on infoAction.
func (s *TagSmoother) Clear() {
	s.mu.Lock()
	s.pending = s.pending[:0]
	s.pcTags.Store(0)
	s.pcUniqueTags.Store(0)
	s.pcPermanentUnique.Store(0)
	s.mu.Unlock()

	s.tagSet.Clear()
	s.permanentSet.Clear()
}

func (s *TagSmoother) run() {
	ticker := time.NewTicker(smoothTickInterval)
	defer ticker.Stop()
	for range ticker.C {
		s.drain()
	}
}

// drain commits all entries whose release time has passed.
func (s *TagSmoother) drain() {
	s.mu.Lock()
	if len(s.pending) == 0 {
		s.mu.Unlock()
		return
	}

	now := time.Now()

	// Partition in-place: collect due entries into a local batch,
	// compact the rest back into pending.
	var batch []smoothEntry
	j := 0
	for _, e := range s.pending {
		if !e.releaseAt.After(now) {
			batch = append(batch, e)
		} else {
			s.pending[j] = e
			j++
		}
	}
	s.pending = s.pending[:j]
	s.mu.Unlock()

	for _, e := range batch {
		s.pcTags.Add(1)
		s.tagSet.Insert(e.epc)
		s.permanentSet.Insert(e.epc)
	}
	if len(batch) > 0 {
		s.pcUniqueTags.Store(int32(s.tagSet.Count()))
		s.pcPermanentUnique.Store(int32(s.permanentSet.Count()))
	}
}
