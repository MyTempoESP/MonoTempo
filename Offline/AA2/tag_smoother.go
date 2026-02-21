package main

import (
	"aa2/intSet"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// smoothWindow is the target duration over which a burst is spread.
	smoothWindow = 1 * time.Second

	// smoothMaxDelay is the age at which buffered tags are force-flushed.
	smoothMaxDelay = 2 * time.Second

	// smoothTickInterval is the drain goroutine's wake interval.
	smoothTickInterval = 10 * time.Millisecond
)

// smoothEntry is one buffered tag waiting to be committed to display counters.
type smoothEntry struct {
	epc       int
	arrivedAt time.Time
}

// TagSmoother buffers incoming RFID tag events and releases them at a steady
// drip rate into the caller-owned PCData counters and intSets, so that the
// display counter increments smoothly instead of jumping in chunks.
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

// Push enqueues one tag for smoothed delivery.
func (s *TagSmoother) Push(epc int) {
	e := smoothEntry{epc: epc, arrivedAt: time.Now()}
	s.mu.Lock()
	s.pending = append(s.pending, e)
	s.mu.Unlock()
}

// Clear discards all buffered tags and resets counters. Called on infoAction.
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

// drain releases the appropriate number of buffered entries per tick.
//
// Drip rate: ceil(N / ticksPerWindow) entries per tick, where
// ticksPerWindow = smoothWindow / smoothTickInterval = 100.
// This spreads N tags over ~1 second.
//
// Any entry older than smoothMaxDelay is released unconditionally.
func (s *TagSmoother) drain() {
	s.mu.Lock()
	n := len(s.pending)
	if n == 0 {
		s.mu.Unlock()
		return
	}

	now := time.Now()

	const ticksPerWindow = int(smoothWindow / smoothTickInterval) // 100

	// Ceiling division: always release at least 1 per tick.
	release := (n + ticksPerWindow - 1) / ticksPerWindow

	// Force-flush entries that have exceeded the max delay.
	deadline := now.Add(-smoothMaxDelay)
	forced := 0
	for forced < n && s.pending[forced].arrivedAt.Before(deadline) {
		forced++
	}
	if forced > release {
		release = forced
	}
	if release > n {
		release = n
	}

	batch := make([]smoothEntry, release)
	copy(batch, s.pending[:release])
	s.pending = s.pending[release:]
	s.mu.Unlock()

	for _, e := range batch {
		s.pcTags.Add(1)
		s.tagSet.Insert(e.epc)
		s.permanentSet.Insert(e.epc)
	}

	s.pcUniqueTags.Store(int32(s.tagSet.Count()))
	s.pcPermanentUnique.Store(int32(s.permanentSet.Count()))
}
