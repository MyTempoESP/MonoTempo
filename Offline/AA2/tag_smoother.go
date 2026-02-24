package main

import (
	"aa2/intSet"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// smoothWindow is the duration over which a batch of tags is linearly
	// interpolated into the display counters.
	smoothWindow = 2 * time.Second

	// smoothMaxDelay is the age at which buffered tags are force-flushed.
	smoothMaxDelay = 12 * time.Second

	// smoothTickInterval is the drain goroutine's wake interval.
	smoothTickInterval = 50 * time.Millisecond
)

// smoothEntry is one buffered tag waiting to be committed to display counters.
type smoothEntry struct {
	epc       int
	arrivedAt time.Time
}

// TagSmoother buffers incoming RFID tag events and releases them via linear
// interpolation into the caller-owned PCData counters and intSets, so that
// the display counter increments smoothly instead of jumping in chunks.
//
// The smoother tracks how many entries have been "released" (committed to
// display) as a floating-point lerp target.  Each tick it computes:
//
//	target = lerp(lerpFrom, lerpTo, elapsed/smoothWindow)
//
// and releases int(target)-released entries.  When new tags arrive while a
// lerp is in progress the target is extended without restarting the clock,
// so the rate increases but the curve stays linear.
//
// Antenna recording is NOT handled here; callers must call
// pcData.Antennas[...].Record() directly.
type TagSmoother struct {
	mu      sync.Mutex
	pending []smoothEntry

	// Linear interpolation state.
	// released counts total entries committed since the last Clear.
	// queued   counts total entries ever pushed since the last Clear.
	// lerpFrom is the released count at the moment the current lerp started.
	// lerpTo   is the target released count at lerpStart+smoothWindow.
	// lerpStart is zero until the first entry is pushed.
	released  int
	queued    int
	lerpFrom  float64
	lerpTo    float64
	lerpStart time.Time

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
	s.queued++
	s.mu.Unlock()
}

// Clear discards all buffered tags and resets all counters. Called on infoAction.
func (s *TagSmoother) Clear() {
	s.mu.Lock()
	s.pending = s.pending[:0]
	s.released = 0
	s.queued = 0
	s.lerpFrom = 0
	s.lerpTo = 0
	s.lerpStart = time.Time{}
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

// drain uses linear interpolation to decide how many buffered entries to
// release this tick.
//
// Lerp lifecycle:
//   - On the first arrival (or after the previous lerp completed), a new lerp
//     starts: lerpFrom = released, lerpTo = queued, lerpStart = now.
//   - While a lerp is active and new tags arrive, lerpTo is extended so the
//     existing linear ramp covers the new arrivals without restarting the clock.
//   - Entries older than smoothMaxDelay are force-flushed, after which the
//     lerp is restarted from the new released position.
func (s *TagSmoother) drain() {
	s.mu.Lock()

	if len(s.pending) == 0 {
		s.mu.Unlock()
		return
	}

	now := time.Now()

	// Detect new arrivals: queued grew beyond our current lerp target.
	if s.queued > int(s.lerpTo) {
		elapsed := now.Sub(s.lerpStart)
		if s.lerpStart.IsZero() || elapsed >= smoothWindow {
			// No active lerp — start one from the current position.
			s.lerpFrom = float64(s.released)
			s.lerpStart = now
		}
		// Extend the target to include all newly queued entries.
		s.lerpTo = float64(s.queued)
	}

	// Linear interpolation: t ∈ [0, 1] over smoothWindow.
	t := float64(now.Sub(s.lerpStart)) / float64(smoothWindow)
	if t > 1.0 {
		t = 1.0
	}
	lerpTarget := s.lerpFrom + t*(s.lerpTo-s.lerpFrom)
	toRelease := int(lerpTarget) - s.released

	// Force-flush any entries that have waited longer than smoothMaxDelay.
	deadline := now.Add(-smoothMaxDelay)
	forced := 0
	for forced < len(s.pending) && s.pending[forced].arrivedAt.Before(deadline) {
		forced++
	}
	if forced > toRelease {
		toRelease = forced
		// Restart the lerp from the post-flush position.
		s.lerpFrom = float64(s.released + forced)
		s.lerpTo = float64(s.queued)
		s.lerpStart = now
	}

	if toRelease <= 0 {
		s.mu.Unlock()
		return
	}
	if toRelease > len(s.pending) {
		toRelease = len(s.pending)
	}

	batch := make([]smoothEntry, toRelease)
	copy(batch, s.pending[:toRelease])
	s.pending = s.pending[toRelease:]
	s.released += toRelease
	s.mu.Unlock()

	for _, e := range batch {
		s.pcTags.Add(1)
		s.tagSet.Insert(e.epc)
		s.permanentSet.Insert(e.epc)
	}

	s.pcUniqueTags.Store(int32(s.tagSet.Count()))
	s.pcPermanentUnique.Store(int32(s.permanentSet.Count()))
}
