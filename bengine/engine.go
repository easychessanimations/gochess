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

// retrieveHash gets from GlobalHashTable the current position
func (eng *Engine) retrieveHash() hashEntry {
	entry := GlobalHashTable.get(eng.Position)
	if entry.kind == 0 || entry.move != NullMove && !eng.Position.IsPseudoLegal(entry.move) {
		eng.Stats.CacheMiss++
		return hashEntry{}
	}

	// return mate score relative to root
	// the score was adjusted relative to position before the hash table was updated
	if entry.score < KnownLossScore {
		entry.score += int16(eng.ply())
	} else if entry.score > KnownWinScore {
		entry.score -= int16(eng.ply())
	}

	eng.Stats.CacheHit++
	return entry
}

// updateHash updates GlobalHashTable with the current position
func (eng *Engine) updateHash(flags hashFlags, depth, score int32, move Move, static int32) {
	// if search is stopped then score cannot be trusted
	if eng.stopped {
		return
	}
	// update principal variation table in exact nodes
	if flags&exact != 0 {
		eng.pvTable.Put(eng.Position, move)
	}
	if eng.ply() == 0 && (len(eng.ignoreRootMoves) != 0 || len(eng.onlyRootMoves) != 0) {
		// at root if there are moves to ignore (e.g. because of multipv)
		// then this is an incomplete search, so don't update the hash
		return
	}

	// save the mate score relative to the current position
	// when retrieving from hash the score will be adjusted relative to root
	if score < KnownLossScore {
		score -= eng.ply()
	} else if score > KnownWinScore {
		score += eng.ply()
	}

	GlobalHashTable.put(eng.Position, hashEntry{
		kind:   flags,
		score:  int16(score),
		depth:  int8(depth),
		move:   move,
		static: int16(static),
	})
}

// searchQuiescence evaluates the position after solving all captures.
//
// This is a very limited search which considers only some violent moves.
// Depth is ignored, so hash uses depth 0. Search continues until
// stand pat or no capture can improve the score.
func (eng *Engine) searchQuiescence(α, β int32) int32 {
	eng.Stats.Nodes++

	entry := eng.retrieveHash()
	if score := int32(entry.score); isInBounds(entry.kind, α, β, score) {
		return score
	}

	static := eng.cachedScore(&entry)
	if static >= β {
		// Stand pat if the static score is already a cut-off.
		eng.updateHash(failedHigh|hasStatic, 0, static, entry.move, static)
		return static
	}

	pos := eng.Position
	us := pos.Us()
	inCheck := pos.IsChecked(us)
	localα := max(α, static)
	bestMove := entry.move

	eng.stack.GenerateMoves(Violent, NullMove)
	for move := eng.stack.PopMove(); move != NullMove; move = eng.stack.PopMove() {
		// Prune futile moves that would anyway result in a stand-pat at that next depth.
		if !inCheck && isFutile(pos, static, α, futilityMargin, move) ||
			!inCheck && seeSign(pos, move) {
			continue
		}

		// Discard illegal or losing captures.
		eng.DoMove(move)
		if eng.Position.IsChecked(us) {
			eng.UndoMove()
			continue
		}
		score := -eng.searchQuiescence(-β, -localα)
		eng.UndoMove()

		if score >= β {
			eng.updateHash(failedHigh|hasStatic, 0, score, move, static)
			return score
		}
		if score > localα {
			localα = score
			bestMove = move
		}
	}

	eng.updateHash(getBound(α, β, localα)|hasStatic, 0, localα, bestMove, static)
	return localα
}

/////////////////////////////////////////////////////////////////////
