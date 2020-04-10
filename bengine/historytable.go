package bengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	. "github.com/easychessanimations/gochess/butils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

// newSearch updates the stats before a new search
func (ht *historyTable) newSearch() {
	for i := range ht {
		ht[i].stat /= 8
	}
}

// get returns the stats for a move m
// if the move is not in the table, returns 0
func (ht *historyTable) get(m Move) int32 {
	h := historyHash(m)
	if ht[h].move != m {
		return 0
	}
	return ht[h].stat
}

// inc increments the counters for m
// evicts an old move if necessary
func (ht *historyTable) add(m Move, delta int32) {
	h := historyHash(m)
	if ht[h].move != m {
		ht[h] = historyEntry{stat: delta, move: m}
	} else {
		ht[h].stat += delta
	}
}

/////////////////////////////////////////////////////////////////////
