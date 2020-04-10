package bengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	. "github.com/easychessanimations/gochess/butils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

// put puts a new entry in the cache
func (c *pawnsTable) put(lock uint64, white, black Accum) {
	indx := lock & uint64(len(*c)-1)
	c[indx] = pawnsEntry{lock, white, black}
}

// get gets an entry from the cache
func (c *pawnsTable) get(lock uint64) (Accum, Accum, bool) {
	indx := lock & uint64(len(*c)-1)
	return c[indx].white, c[indx].black, c[indx].lock == lock
}

// load evaluates position, using the cache if possible
func (c *pawnsTable) load(pos *Position) (Accum, Accum) {
	h := pawnsHash(pos)
	white, black, ok := c.get(h)
	if !ok {
		white = evaluatePawnsAndShelter(pos, White)
		black = evaluatePawnsAndShelter(pos, Black)
		c.put(h, white, black)
	}
	return white, black
}

/////////////////////////////////////////////////////////////////////
