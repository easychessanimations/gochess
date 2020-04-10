package bengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	. "github.com/easychessanimations/gochess/butils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// hash table

// prefetch prefetches the hash entry into lower caches
func (ht *HashTable) prefetch(pos *Position) {
	_, key0, _ := split(pos.Zobrist(), ht.mask)
	prefetch(&ht.table[key0&^1])
}

// put puts a new entry in the database
func (ht *HashTable) put(pos *Position, entry hashEntry) {
	lock, key0, key1 := split(pos.Zobrist(), ht.mask)
	entry.lock = lock

	if e := &ht.table[key0]; e.lock == lock || e.kind == 0 || e.depth >= entry.depth {
		ht.table[key0] = entry
	} else {
		ht.table[key1] = entry
	}
}

// get returns the hash entry for position
//
// observation: due to collision errors, the hashEntry returned might be
// from a different table; however, these errors are not common because
// we use 32-bit lock + log_2(len(ht.table)) bits to avoid collisions
func (ht *HashTable) get(pos *Position) hashEntry {
	lock, key0, key1 := split(pos.Zobrist(), ht.mask)
	if ht.table[key0].lock == lock {
		return ht.table[key0]
	}
	if ht.table[key1].lock == lock {
		return ht.table[key1]
	}
	return hashEntry{}
}

// Clear removes all entries from hash
func (ht *HashTable) Clear() {
	for i := range ht.table {
		ht.table[i] = hashEntry{}
	}
}

/////////////////////////////////////////////////////////////////////
