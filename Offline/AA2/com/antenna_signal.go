// Package com
package com

import (
	"sync"
	"time"
)

const (
	signalWindowDuration = 100 * time.Millisecond
	signalThresholdLow   = 10.0 // tags/sec
	signalThresholdHigh  = 13.0 // tags/sec
)

type AntennaSignal struct {
	mu     sync.Mutex
	window []time.Time
}

// Record registers a tag detection event at the current time.
func (a *AntennaSignal) Record() {
	now := time.Now()
	a.mu.Lock()
	a.prune(now)
	a.window = append(a.window, now)
	a.mu.Unlock()
}

// Level returns the signal strength level (0-3) based on tag frequency
// in the rolling window.
//
//	0 = no readings
//	1 = Low  (< signalThresholdLow tags/sec)
//	2 = Medium (< signalThresholdHigh tags/sec)
//	3 = High (>= signalThresholdHigh tags/sec)
func (a *AntennaSignal) Level() int {
	now := time.Now()
	a.mu.Lock()
	a.prune(now)
	count := len(a.window)
	a.mu.Unlock()

	if count == 0 {
		return 0
	}

	rate := float64(count) / signalWindowDuration.Seconds()

	if rate >= signalThresholdHigh {
		return 3
	}
	if rate >= signalThresholdLow {
		return 2
	}
	return 1
}

// Clear resets the signal window.
func (a *AntennaSignal) Clear() {
	a.mu.Lock()
	a.window = a.window[:0]
	a.mu.Unlock()
}

// prune removes entries older than the window duration. Must be called with mu held.
func (a *AntennaSignal) prune(now time.Time) {
	cutoff := now.Add(-signalWindowDuration)
	i := 0
	for i < len(a.window) && a.window[i].Before(cutoff) {
		i++
	}
	if i > 0 {
		a.window = a.window[i:]
	}
}
