// smoother_sim simulates the TagSmoother drip algorithm so you can observe its
// behaviour at a glance without a live RFID device.
//
// Usage:
//
//	go run ./test/smoother_sim/ [flags]
//
// Flags:
//
//	-scenario  burst | steady | floating | sparse | mixed  (default: floating)
//	-window    smoothing spread window                      (default: 2s)
//	-max-delay max age before force-flush                   (default: 12s)
//	-tick      drain interval                               (default: 10ms)
//	-rate      base tag arrival rate tags/sec               (default: 10)
//	-duration  total simulation run time                    (default: 15s)
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// ── tuneable constants (mirrors tag_smoother.go) ─────────────────────────────

var (
	smoothWindow       = 2 * time.Second
	smoothMaxDelay     = 12 * time.Second
	smoothTickInterval = 10 * time.Millisecond
)

// ── minimal smoother (self-contained, no intSet dependency) ──────────────────

type smoothEntry struct {
	epc       int
	arrivedAt time.Time
}

type smoother struct {
	mu        sync.Mutex
	pending   []smoothEntry
	dripCarry float64 // fractional drip accumulator across ticks

	// counters
	totalTags  atomic.Int64
	uniqueTags atomic.Int64

	// internal unique tracking
	seen   map[int]struct{}
	seenMu sync.Mutex

	// per-tick stats for printing
	lastReleased int
	lastArrived  int
}

func newSmoother() *smoother {
	return &smoother{
		pending: make([]smoothEntry, 0, 256),
		seen:    make(map[int]struct{}),
	}
}

func (s *smoother) push(epc int) {
	e := smoothEntry{epc: epc, arrivedAt: time.Now()}
	s.mu.Lock()
	s.pending = append(s.pending, e)
	s.lastArrived++
	s.mu.Unlock()
}

// drain releases the appropriate batch and returns (pending, released) counts.
func (s *smoother) drain() (pending, released int) {
	s.mu.Lock()
	n := len(s.pending)
	if n == 0 {
		s.lastReleased = 0
		s.lastArrived = 0
		s.mu.Unlock()
		return 0, 0
	}

	now := time.Now()
	ticksPerWindow := int(smoothWindow / smoothTickInterval)
	ticksPerWindow = max(1, ticksPerWindow)

	s.dripCarry += float64(n) / float64(ticksPerWindow)
	release := int(s.dripCarry)
	s.dripCarry -= float64(release)

	deadline := now.Add(-smoothMaxDelay)
	forced := 0
	for forced < n && s.pending[forced].arrivedAt.Before(deadline) {
		forced++
	}
	if forced > release {
		release = forced
		s.dripCarry = 0 // reset carry after a force-flush
	}
	if release > n {
		release = n
	}

	batch := make([]smoothEntry, release)
	copy(batch, s.pending[:release])
	s.pending = s.pending[release:]
	arrived := s.lastArrived
	s.lastArrived = 0
	s.mu.Unlock()

	for _, e := range batch {
		s.totalTags.Add(1)
		s.seenMu.Lock()
		if _, exists := s.seen[e.epc]; !exists {
			s.seen[e.epc] = struct{}{}
			s.uniqueTags.Add(1)
		}
		s.seenMu.Unlock()
	}

	s.mu.Lock()
	s.lastReleased = release
	_ = arrived
	s.mu.Unlock()

	return n - release, release
}

// ── arrival scenarios ─────────────────────────────────────────────────────────

// fracAcc accumulates fractional tag counts across ticks so that e.g.
// 0.1 tags/tick reliably produces 1 tag every 10 ticks rather than 0 forever.
type fracAcc struct{ carry float64 }

func (a *fracAcc) next(tagsPerTick float64) int {
	a.carry += tagsPerTick
	n := int(a.carry)
	a.carry -= float64(n)
	return n
}

type scenario func(elapsed time.Duration, baserate float64, acc *fracAcc) int

// returns how many tags arrive this tick (elapsed = time since sim start)
var scenarios = map[string]scenario{
	// large burst at t=0, then silence
	"burst": func(elapsed time.Duration, baserate float64, _ *fracAcc) int {
		if elapsed <= smoothTickInterval {
			return int(baserate * 50) // 50× burst on first tick only
		}
		return 0
	},

	// constant rate via fractional accumulator
	"steady": func(_ time.Duration, baserate float64, acc *fracAcc) int {
		return acc.next(baserate * smoothTickInterval.Seconds())
	},

	// rate oscillates sinusoidally between 0.25× and 1.75× baserate
	"floating": func(elapsed time.Duration, baserate float64, acc *fracAcc) int {
		secs := elapsed.Seconds()
		// period ~8 s, amplitude ±0.75× baserate
		rate := baserate * (1.0 + 0.75*math.Sin(2*math.Pi*secs/8.0))
		if rate < 0 {
			rate = 0
		}
		return acc.next(rate * smoothTickInterval.Seconds())
	},

	// infrequent individual tags — one every ~3 s; exercises force-flush path
	"sparse": func(elapsed time.Duration, baserate float64, _ *fracAcc) int {
		ms := elapsed.Milliseconds()
		if ms > 0 && ms%3000 < int64(smoothTickInterval.Milliseconds()) {
			return int(baserate)
		}
		return 0
	},

	// burst at t=0, then steady trickle at 40% baserate
	"mixed": func(elapsed time.Duration, baserate float64, acc *fracAcc) int {
		if elapsed < 2*smoothTickInterval {
			return int(baserate * 30)
		}
		return acc.next(baserate * 0.4 * smoothTickInterval.Seconds())
	},
}

// ── main ──────────────────────────────────────────────────────────────────────

func main() {
	scenarioName := flag.String("scenario", "floating", "burst | steady | floating | sparse | mixed")
	windowFlag := flag.Duration("window", 2*time.Second, "smoothing spread window")
	maxDelayFlag := flag.Duration("max-delay", 12*time.Second, "max age before force-flush")
	tickFlag := flag.Duration("tick", 10*time.Millisecond, "drain interval")
	rateFlag := flag.Float64("rate", 10, "base tag arrival rate (tags/sec)")
	durationFlag := flag.Duration("duration", 15*time.Second, "total simulation time")
	noZerosFlag := flag.Bool("no-zeros", false, "suppress ticks where both arrived and released are 0")
	flag.Parse()

	// apply flags to package-level vars used by the smoother
	smoothWindow = *windowFlag
	smoothMaxDelay = *maxDelayFlag
	smoothTickInterval = *tickFlag

	arrivalFn, ok := scenarios[*scenarioName]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown scenario %q; choose: burst | steady | floating | sparse | mixed\n", *scenarioName)
		os.Exit(1)
	}

	fmt.Printf("scenario=%-10s  window=%v  max-delay=%v  tick=%v  base-rate=%.1f/s  duration=%v\n\n",
		*scenarioName, *windowFlag, *maxDelayFlag, *tickFlag, *rateFlag, *durationFlag)
	fmt.Printf("%-10s  %-8s  %-8s  %-8s  %-8s  %-8s\n",
		"elapsed", "arrived", "pending", "released", "total", "unique")
	fmt.Println("----------  --------  --------  --------  --------  --------")

	sim := newSmoother()
	ticker := time.NewTicker(*tickFlag)
	defer ticker.Stop()

	acc := &fracAcc{}
	start := time.Now()
	var (
		totalPushed int
		peakPending int
		maxDripRate float64
		minDripRate = math.MaxFloat64
	)

	for t := range ticker.C {
		elapsed := t.Sub(start)
		if elapsed > *durationFlag {
			break
		}

		// push arrivals for this tick
		arrivals := arrivalFn(elapsed, *rateFlag, acc)
		for i := 0; i < arrivals; i++ {
			sim.push(totalPushed + i) // unique epc per tag
		}
		totalPushed += arrivals

		// drain
		pending, released := sim.drain()

		// stats
		if pending > peakPending {
			peakPending = pending
		}
		tickSecs := tickFlag.Seconds()
		dripRate := float64(released) / tickSecs
		if released > 0 && dripRate > maxDripRate {
			maxDripRate = dripRate
		}
		if released > 0 && dripRate < minDripRate {
			minDripRate = dripRate
		}

		total := sim.totalTags.Load()
		unique := sim.uniqueTags.Load()

		if !*noZerosFlag || arrivals > 0 || released > 0 {
			fmt.Printf("t=%8.3fs  arr=%5d  pend=%5d  rel=%5d  tot=%6d  uniq=%6d\n",
				elapsed.Seconds(), arrivals, pending, released, total, unique)
		}
	}

	// final drain — empty anything still pending
	for {
		pending, released := sim.drain()
		if pending == 0 && released == 0 {
			break
		}
		total := sim.totalTags.Load()
		unique := sim.uniqueTags.Load()
		fmt.Printf("t=%8s  arr=%5d  pend=%5d  rel=%5d  tot=%6d  uniq=%6d  (drain)\n",
			"end", 0, pending, released, total, unique)
	}

	if minDripRate == math.MaxFloat64 {
		minDripRate = 0
	}

	fmt.Println("\n── summary ─────────────────────────────────────────────────")
	fmt.Printf("  tags pushed   : %d\n", totalPushed)
	fmt.Printf("  tags released : %d\n", sim.totalTags.Load())
	fmt.Printf("  unique tags   : %d\n", sim.uniqueTags.Load())
	fmt.Printf("  peak pending  : %d\n", peakPending)
	fmt.Printf("  drip rate     : %.1f – %.1f tags/sec\n", minDripRate, maxDripRate)
}
