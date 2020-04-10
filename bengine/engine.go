package bengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	. "github.com/easychessanimations/gochess/butils"
)

/////////////////////////////////////////////////////////////////////

// package bengine implements board, move generation and position searching
//
// the package can be used as a general library for chess tool writing and
// provides the core functionality for the zurichess chess engine
//
// position (basic.go, position.go) uses:
//
//   * bitboards for representation - https://chessprogramming.wikispaces.com/Bitboards
//   * magic bitboards for sliding move generation - https://chessprogramming.wikispaces.com/Magic+Bitboards
//
// search (engine.go) features implemented are:
//
//   * aspiration window - https://chessprogramming.wikispaces.com/Aspiration+Windows
//   * check extension - https://chessprogramming.wikispaces.com/Check+Extensions
//   * fail soft - https://chessprogramming.wikispaces.com/Fail-Soft
//   * futility Pruning - https://chessprogramming.wikispaces.com/Futility+pruning
//   * history leaf pruning - https://chessprogramming.wikispaces.com/History+Leaf+Pruning
//   * killer move heuristic - https://chessprogramming.wikispaces.com/Killer+Heuristic
//   * late move redution (LMR) - https://chessprogramming.wikispaces.com/Late+Move+Reductions
//   * mate distance pruning - https://chessprogramming.wikispaces.com/Mate+Distance+Pruning
//   * negamax framework - http://chessprogramming.wikispaces.com/Alpha-Beta#Implementation-Negamax%20Framework
//   * null move prunning (NMP) - https://chessprogramming.wikispaces.com/Null+Move+Pruning
//   * principal variation search (PVS) - https://chessprogramming.wikispaces.com/Principal+Variation+Search
//   * quiescence search - https://chessprogramming.wikispaces.com/Quiescence+Search
//   * razoring - https://chessprogramming.wikispaces.com/Razoring
//   * static Single Evaluation - https://chessprogramming.wikispaces.com/Static+Exchange+Evaluation
//   * zobrist hashing - https://chessprogramming.wikispaces.com/Zobrist+Hashing
//
// move ordering (move_ordering.go) consists of:
//
//   * hash move heuristic
//   * captures sorted by MVVLVA - https://chessprogramming.wikispaces.com/MVV-LVA
//   * killer moves - https://chessprogramming.wikispaces.com/Killer+Move
//   * history Heuristic - https://chessprogramming.wikispaces.com/History+Heuristic
//   * countermove Heuristic - https://chessprogramming.wikispaces.com/Countermove+Heuristic
//
// evaluation (material.go) consists of
//
//   * material and mobility
//   * piece square tables
//   * king pawn shield - https://chessprogramming.wikispaces.com/King+Safety
//   * king safery ala Toga style - https://chessprogramming.wikispaces.com/King+Safety#Attacking%20King%20Zone
//   * pawn structure: connected, isolated, double, passed, rammed. Evaluation is cached (see cache.go)
//   * attacks on minors and majors
//   * rooks on open and semiopenfiles - https://chessprogramming.wikispaces.com/Rook+on+Open+File
//   * tapered evaluation - https://chessprogramming.wikispaces.com/Tapered+Eval

/////////////////////////////////////////////////////////////////////
// member functions

// ply returns the ply from the beginning of the search
func (eng *Engine) ply() int32 {
	return int32(eng.Position.Ply - eng.rootPly)
}

// SetPosition sets current position
// if pos is nil, the starting position is set
func (eng *Engine) SetPosition(pos *Position) {
	if pos != nil {
		eng.Position = pos
	} else {
		eng.Position, _ = PositionFromFEN(FENStartPos)
	}
}

// DoMove executes a move.
func (eng *Engine) DoMove(move Move) {
	eng.Position.DoMove(move)
	GlobalHashTable.prefetch(eng.Position)
}

// UndoMove undoes the last move
func (eng *Engine) UndoMove() {
	eng.Position.UndoMove()
}

// Score evaluates current position from current player's POV
func (eng *Engine) Score() int32 {
	return Evaluate(eng.Position).GetCentipawnsScore() * eng.Position.Us().Multiplier()
}

// cachedScore implements a cache on top of Score
// the cached static evaluation is stored in the hashEntry
func (eng *Engine) cachedScore(e *hashEntry) int32 {
	if e.kind&hasStatic == 0 {
		e.kind |= hasStatic
		e.static = int16(eng.Score())
	}
	return int32(e.static)
}

// endPosition determines whether the current position is an end game
// returns score and a bool if the game has ended
func (eng *Engine) endPosition() (int32, bool) {
	pos := eng.Position // shortcut
	// trivial cases when kings are missing
	if Kings(pos, White) == 0 {
		if Kings(pos, Black) == 0 {
			return 0, true // both kings are missing
		}
		return pos.Us().Multiplier() * (MatedScore + eng.ply()), true
	}
	if Kings(pos, Black) == 0 {
		return pos.Us().Multiplier() * (MateScore - eng.ply()), true
	}
	// neither side can mate
	if pos.InsufficientMaterial() {
		return 0, true
	}
	// fifty full moves without a capture or a pawn move
	if pos.FiftyMoveRule() {
		return 0, true
	}
	// repetition is a draw
	// at root we need to continue searching even if we saw two repetitions already,
	// however we can prune deeper search only at two repetitions
	if r := pos.ThreeFoldRepetition(); eng.ply() > 0 && r >= 2 || r >= 3 {
		return 0, true
	}
	return 0, false
}

/////////////////////////////////////////////////////////////////////
