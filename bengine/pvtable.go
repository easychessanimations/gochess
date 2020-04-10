package bengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	. "github.com/easychessanimations/gochess/butils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// pv table

// put inserts a new entry;  ignores NullMoves
func (pv pvTable) Put(pos *Position, move Move) {
	if move != NullMove {
		zobrist := pos.Zobrist()
		pv[zobrist&pvTableMask] = pvEntry{
			lock: zobrist,
			move: move,
		}
	}
}

// TODO: kook up move in transposition table if none is available
func (pv pvTable) get(pos *Position) Move {
	zobrist := pos.Zobrist()
	if entry := &pv[zobrist&pvTableMask]; entry.lock == zobrist {
		return entry.move
	}
	return NullMove
}

// get returns the principal variation from pos
func (pv pvTable) Get(pos *Position) []Move {
	seen := make(map[uint64]bool)
	var moves []Move
	// extract the moves by following the position
	next := pv.get(pos)
	for next != NullMove && !seen[pos.Zobrist()] {
		seen[pos.Zobrist()] = true
		moves = append(moves, next)
		pos.DoMove(next)
		next = pv.get(pos)
	}
	// undo all moves, so we get back to the initial state
	for range moves {
		pos.UndoMove()
	}
	// add the last repeated move
	if next != NullMove {
		moves = append(moves, next)
	}
	return moves
}

/////////////////////////////////////////////////////////////////////
