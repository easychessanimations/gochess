package bengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	"time"

	. "github.com/easychessanimations/gochess/butils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

// thinkingTime calculates how much time to think this round
// t is the remaining time, i is the increment
func (tc *TimeControl) thinkingTime() time.Duration {
	// the formula allows engine to use more of time in the begining
	// and rely more on the increment later
	tmp := time.Duration(tc.MovesToGo)
	tt := (tc.time + (tmp-1)*tc.inc) / tmp

	if tt < 0 {
		return 0
	}
	if tc.predicted {
		tt = tt * 4 / 3
	}
	if tt < tc.limit {
		return tt
	}
	return tc.limit
}

// start starts the timer
// should start as soon as possible to set the correct time
func (tc *TimeControl) Start(ponder bool) {
	if tc.sideToMove == White {
		tc.time, tc.inc = tc.WTime, tc.WInc
	} else {
		tc.time, tc.inc = tc.BTime, tc.BInc
	}

	// calcuates the last moment when the search should be stopped
	if tc.time > 2*overhead {
		tc.limit = tc.time - overhead
	} else if tc.time > overhead {
		tc.limit = overhead
	} else {
		tc.limit = tc.time
	}

	// if there are still many moves to go, don't use all the time
	tc.limit /= time.Duration(min(tc.MovesToGo, 5))

	// increase the branchFactor a bit to be on the
	// safe side when there are only a few moves left
	for i := int32(4); i > 0; i /= 2 {
		if tc.MovesToGo <= i {
			tc.branch += 16
		}
	}

	tc.stopped = atomicFlag{flag: false}
	tc.ponderhit = atomicFlag{flag: !ponder}

	tc.searchTime = tc.thinkingTime()
	tc.updateDeadlines() // deadlines are ignored while pondering (ponderHit == false)
}

func (tc *TimeControl) updateDeadlines() {
	now := time.Now()
	tc.searchDeadline = now.Add(tc.searchTime / time.Duration(tc.branch/16))

	// stopDeadline is when to abort the search in case of an explosion
	// we give a large overhead here so the search is not aborted very often
	deadline := tc.searchTime * 4
	if deadline > tc.limit {
		deadline = tc.limit
	}
	tc.stopDeadline = now.Add(deadline)
}

// NextDepth returns true if search can start at depth
// in any case Stopped() will return false
func (tc *TimeControl) NextDepth(depth int32) bool {
	tc.currDepth = depth
	return tc.currDepth <= tc.Depth && !tc.hasStopped(tc.searchDeadline)
}

// PonderHit switch to our time control
func (tc *TimeControl) PonderHit() {
	tc.updateDeadlines()
	tc.ponderhit.set()
}

// Stop marks the search as stopped
func (tc *TimeControl) Stop() {
	tc.stopped.set()
}

func (tc *TimeControl) hasStopped(deadline time.Time) bool {
	if tc.currDepth <= 2 {
		// run for at few depths at least otherwise mates can be missed
		return false
	}
	if tc.stopped.get() {
		// use a cached value if available
		return true
	}
	if tc.ponderhit.get() && time.Now().After(deadline) {
		// stop search if no longer pondering and deadline as passed
		return true
	}
	return false
}

// Stopped returns true if the search has stopped because
// Stop() was called or the time has ran out
func (tc *TimeControl) Stopped() bool {
	if !tc.hasStopped(tc.stopDeadline) {
		return false
	}
	// time has ran out so flip the stopped flag
	tc.stopped.set()
	return true
}

/////////////////////////////////////////////////////////////////////
