package bengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	"math/rand"

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

// tryMove descends on the search tree; this function
// is called from searchTree after the move is executed
// and it will undo the move
//
// α, β represent lower and upper bounds
// depth is the remaining depth (decreasing)
// lmr is how much to reduce a late move. Implies non-null move
// nullWindow indicates whether to scout first. Implies non-null move
//
// returns the score from the deeper search
func (eng *Engine) tryMove(α, β, depth, lmr int32, nullWindow bool) int32 {
	depth--

	score := α + 1
	if lmr > 0 { // reduce late moves
		score = -eng.searchTree(-α-1, -α, depth-lmr)
	}

	if score > α { // if late move reduction is disabled or has failed
		if nullWindow {
			score = -eng.searchTree(-α-1, -α, depth)
			if α < score && score < β {
				score = -eng.searchTree(-β, -α, depth)
			}
		} else {
			score = -eng.searchTree(-β, -α, depth)
		}
	}

	eng.UndoMove()
	return score
}

// isIgnoredRootMove returns true if move should be ignored at root
func (eng *Engine) isIgnoredRootMove(move Move) bool {
	if eng.ply() != 0 {
		return false
	}
	for _, m := range eng.ignoreRootMoves {
		if m == move {
			return true
		}
	}
	for _, m := range eng.onlyRootMoves {
		if m == move {
			return false
		}
	}
	return len(eng.onlyRootMoves) != 0
}

// searchTree implements searchTree framework
//
// searchTree fails soft, i.e. the score returned can be outside the bounds
//
// α, β represent lower and upper bounds
// depth is the search depth (decreasing)
//
// returns the score of the current position up to depth (modulo reductions/extensions)
// the returned score is from current player's POV
//
// invariants:
//   if score <= α then the search failed low and the score is an upper bound
//   else if score >= β then the search failed high and the score is a lower bound
//   else score is exact
//
// assuming this is a maximizing nodes, failing high means that a
// minimizing ancestor node already has a better alternative
func (eng *Engine) searchTree(α, β, depth int32) int32 {
	ply := eng.ply()
	pvNode := α+1 < β
	pos := eng.Position
	us := pos.Us()

	// update statistics
	eng.Stats.Nodes++
	if !eng.stopped && eng.Stats.Nodes >= eng.checkpoint {
		eng.checkpoint = eng.Stats.Nodes + checkpointStep
		if eng.timeControl.Stopped() {
			eng.stopped = true
		}
	}
	if eng.stopped {
		return α
	}
	if pvNode && ply > eng.Stats.SelDepth {
		eng.Stats.SelDepth = ply
	}

	// verify that this is not already an endgame
	if score, done := eng.endPosition(); done && (ply != 0 || score != 0) {
		// at root we ignore draws because some GUIs don't properly detect
		// theoretical draws; e.g. cutechess doesn't detect that kings and
		// bishops when all bishops are on the same color; if the position
		// is a theoretical draw, keep searching for a move
		return score
	}

	// mate pruning: if an ancestor already has a mate in ply moves then
	// the search will always fail low so we return the lowest winning score
	if MateScore-ply <= α {
		return KnownWinScore
	}

	// stop searching when the maximum search depth is reached
	// depth can be < 0 due to aggressive LMR
	if depth <= 0 {
		return eng.searchQuiescence(α, β)
	}

	// check the transposition table
	// entry will store the cached static evaluation which may be computed later
	entry := eng.retrieveHash()
	hash := entry.move
	if eng.isIgnoredRootMove(hash) {
		entry = hashEntry{}
		hash = NullMove
	}
	if score := int32(entry.score); depth <= int32(entry.depth) &&
		isInBounds(entry.kind, α, β, score) &&
		(ply != 0 || !eng.isIgnoredRootMove(hash)) {
		if pvNode {
			// update the pv table, otherwise we risk not having a node at root
			// if the pv entry was overwritten
			eng.pvTable.Put(pos, hash)
		}
		if score >= β && hash != NullMove {
			// if this is a CUT node, update the killer like in the regular move loop
			eng.stack.SaveKiller(hash)
		}
		return score
	}

	sideIsChecked := pos.IsChecked(us)

	// do a null move; if the null move fails high then the current
	// position is too good, so opponent will not play it
	// verification that we are not in check is done by tryMove
	// which bails out if after the null move we are still in check
	if !sideIsChecked && // nullmove is illegal when in check
		MinorsAndMajors(pos, us) != 0 && // at least one minor/major piece.
		KnownLossScore < α && β < KnownWinScore && // disable in lost or won positions
		(entry.kind&hasStatic == 0 || int32(entry.static) >= β) {
		eng.DoMove(NullMove)
		reduction := 1 + depth/3
		score := eng.tryMove(β-1, β, depth-reduction, 0, false)
		if score >= β && score < KnownWinScore {
			return score
		}
	}

	// razoring at very low depth: if QS is under a considerable margin
	// we return that score
	if depth == 1 &&
		!sideIsChecked && // disable in check
		!pvNode && // disable in pv nodes
		KnownLossScore < α && β < KnownWinScore { // disable when searching for a mate
		rα := α - futilityMargin
		if score := eng.searchQuiescence(rα, rα+1); score <= rα {
			return score
		}
	}

	// futility and history pruning at frontier nodes
	// based on Deep Futility Pruning http://home.hccnet.nl/h.g.muller/deepfut.html
	// based on History Leaf Pruning https://chessprogramming.wikispaces.com/History+Leaf+Pruning
	// statically evaluates the position. Use static evaluation from hash if available
	static := int32(0)
	allowLeafsPruning := false
	if depth <= futilityDepthLimit && // enable when close to the frontier
		!sideIsChecked && // disable in check
		!pvNode && // disable in pv nodes
		KnownLossScore < α && β < KnownWinScore { // disable when searching for a mate
		allowLeafsPruning = true
		static = eng.cachedScore(&entry)
	}

	// principal variation search: search with a null window if there is already a good move
	bestMove, localα := NullMove, int32(-InfinityScore)
	// dropped true if not all moves were searched
	// mate cannot be declared unless all moves were tested
	dropped := false
	numMoves := int32(0)

	eng.stack.GenerateMoves(Violent|Quiet, hash)
	for move := eng.stack.PopMove(); move != NullMove; move = eng.stack.PopMove() {
		if ply == 0 {
			if eng.isIgnoredRootMove(move) {
				continue
			}
			eng.Log.CurrMove(int(depth), move, int(numMoves+1))
		}

		givesCheck := pos.GivesCheck(move)
		critical := move == hash || eng.stack.IsKiller(move)
		history := eng.history.get(move)
		newDepth := depth
		numMoves++

		if allowLeafsPruning && !critical && !givesCheck && localα > KnownLossScore {
			// prune moves that do not raise alphas and moves that performed bad historically
			// prune bad captures moves that performed bad historically
			if isFutile(pos, static, α, depth*futilityMargin, move) ||
				history < -10 && move.IsQuiet() ||
				see(pos, move) < -futilityMargin {
				dropped = true
				continue
			}
		}

		// extend good moves that also gives check
		// see discussion: http://www.talkchess.com/forum/viewtopic.php?t=56361
		// when the move gives check, history pruning and futility pruning are also disabled
		if givesCheck && !seeSign(pos, move) {
			newDepth += checkDepthExtension
			critical = true
		}

		// late move reduction: search best moves with full depth, reduce remaining moves
		lmr := int32(0)
		if !sideIsChecked && depth > lmrDepthLimit && !critical {
			// reduce quiet moves and bad captures more at high depths and after many quiet moves
			// large numMoves means it's likely not a CUT node.  Large depth means reductions are less risky
			if move.IsQuiet() {
				if history <= 0 {
					lmr = 2 + min(depth, numMoves)/6
				} else {
					lmr = 1 + min(depth, numMoves)/6
				}
			} else if see := see(pos, move); see < -futilityMargin {
				lmr = 2 + min(depth, numMoves)/6
			} else if see < 0 {
				lmr = 1 + min(depth, numMoves)/6
			}
		}

		// skip illegal moves that leave the king in check
		eng.DoMove(move)
		if pos.IsChecked(us) {
			eng.UndoMove()
			continue
		}

		score := eng.tryMove(max(α, localα), β, newDepth, lmr, numMoves > 1)

		if score >= β {
			// fail high, cut node
			eng.history.add(move, 5+5*depth)
			eng.stack.SaveKiller(move)
			eng.updateHash(failedHigh|(entry.kind&hasStatic), depth, score, move, int32(entry.static))
			return score
		}
		if score > localα {
			bestMove, localα = move, score
		}
		eng.history.add(move, -1)
	}

	bound := getBound(α, β, localα)
	if !dropped && bestMove == NullMove {
		// if no move was found then the game is over
		bound = exact
		if sideIsChecked {
			localα = MatedScore + ply
		} else {
			localα = 0
		}
	}

	eng.updateHash(bound|(entry.kind&hasStatic), depth, localα, bestMove, int32(entry.static))
	return localα
}

// search starts the search up to depth depth
// the returned score is from current side to move POV
// estimated is the score from previous depths
func (eng *Engine) search(depth, estimated int32) int32 {
	// this method only implements aspiration windows
	//
	// the gradual widening algorithm is the one used by RobboLito
	// and Stockfish and it is explained here:
	// http://www.talkchess.com/forum/viewtopic.php?topic_view=threads&p=499768&t=46624
	γ, δ := estimated, int32(initialAspirationWindow)
	α, β := max(γ-δ, -InfinityScore), min(γ+δ, InfinityScore)
	score := estimated

	if depth < 4 {
		// disable aspiration window for very low search depths
		α, β = -InfinityScore, +InfinityScore
	}

	for !eng.stopped {
		// at root a non-null move is required, cannot prune based on null-move
		score = eng.searchTree(α, β, depth)
		if score <= α {
			α = max(α-δ, -InfinityScore)
			δ += δ / 2
		} else if score >= β {
			β = min(β+δ, InfinityScore)
			δ += δ / 2
		} else {
			return score
		}
	}

	return score
}

// searchMultiPV searches eng.options.MultiPV principal variations from current position
// returns score and the moves of the highest scoring pv line (possible empty)
// if a pv is not found (e.g. search is stopped during the first ply), return 0, nil
func (eng *Engine) searchMultiPV(depth, estimated int32) (int32, []Move) {
	type pv struct {
		score int32
		moves []Move
	}

	multiPV := eng.Options.MultiPV
	searchMultiPV := (eng.Options.HandicapLevel+4)/5 + 1
	if multiPV < searchMultiPV {
		multiPV = searchMultiPV
	}

	pvs := make([]pv, 0, multiPV)
	eng.ignoreRootMoves = eng.ignoreRootMoves[:0]
	for p := 0; p < multiPV; p++ {
		if eng.UseAB {
			// search using naive alphabeta
			estimated = eng.searchAB(depth, estimated)
		} else {
			estimated = eng.search(depth, estimated)
		}
		if eng.stopped {
			break // if eng has been stopped then this is not a legit pv
		}

		var moves []Move
		if eng.UseAB {
			// get pev from naive alphabeta's pv table
			moves = eng.pvTableAB.Get(eng.Position)
		} else {
			moves = eng.pvTable.Get(eng.Position)
		}
		hasPV := len(moves) != 0 && !eng.isIgnoredRootMove(moves[0])
		if p == 0 || hasPV { // at depth 0 we might not get a PV
			pvs = append(pvs, pv{estimated, moves})
		}
		if !hasPV {
			break
		}
		// if there is PV ignore the first move for the next PVs
		eng.ignoreRootMoves = append(eng.ignoreRootMoves, moves[0])
	}

	// sort PVs by score
	if len(pvs) == 0 {
		return 0, nil
	}
	for i := range pvs {
		for j := i; j >= 0; j-- {
			if j == 0 || pvs[j-1].score > pvs[i].score {
				tmp := pvs[i]
				copy(pvs[j+1:i+1], pvs[j:i])
				pvs[j] = tmp
				break
			}
		}
	}

	for i := range pvs {
		eng.Log.PrintPV(eng.Stats, i+1, pvs[i].score, pvs[i].moves)
	}

	// for best play return the PV with highest score
	if eng.Options.HandicapLevel == 0 || len(pvs) <= 1 {
		return pvs[0].score, pvs[0].moves
	}

	// PVs are sorted by score. Pick one PV at random
	// and if the score is not too far off, return it
	s := int32(eng.Options.HandicapLevel)
	d := s*s/2 + s*10 + 5
	n := rand.Intn(len(pvs))
	for pvs[n].score+d < pvs[0].score {
		n--
	}
	return pvs[n].score, pvs[n].moves
}

// Play evaluates current position. See PlayMoves for the returned values
func (eng *Engine) Play(tc *TimeControl) (score int32, moves []Move) {
	return eng.PlayMoves(tc, nil)
}

// PlayMoves evaluates current position searching only moves specifid by rootMoves
//
// returns the principal variation, that is
//      moves[0] is the best move found and
//      moves[1] is the pondering move
//
// if rootMoves is nil searches all root moves
//
// returns a nil pv if no move was found because the game is already finished
// returns empty pv array if it's valid position, but no pv was found (e.g. search depth is 0)
//
// Time control, tc, should already be started
func (eng *Engine) PlayMoves(tc *TimeControl, rootMoves []Move) (score int32, moves []Move) {
	if !initialized {
		initEngine()
	}

	eng.Log.BeginSearch()
	eng.Stats = Stats{Depth: -1}

	eng.rootPly = eng.Position.Ply
	eng.timeControl = tc
	eng.stopped = false
	eng.checkpoint = checkpointStep
	eng.stack.Reset(eng.Position)
	eng.history.newSearch()
	eng.onlyRootMoves = rootMoves

	for depth := int32(0); depth < 64; depth++ {
		if !tc.NextDepth(depth) {
			// stop if tc control says we are done
			// search at least one depth, otherwise a move cannot be returned
			break
		}

		eng.Stats.Depth = depth
		if s, m := eng.searchMultiPV(depth, score); len(moves) == 0 || len(m) != 0 {
			score, moves = s, m
		}
	}

	eng.Log.EndSearch()
	if len(moves) == 0 && !eng.Position.HasLegalMoves() {
		return 0, nil
	} else if moves == nil {
		return score, []Move{}
	}
	return score, moves
}

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

// searchQuiescence evaluates the position after solving all captures
//
// this is a very limited search which considers only some violent moves
// depth is ignored, so hash uses depth 0; search continues until
// stand pat or no capture can improve the score
func (eng *Engine) searchQuiescence(α, β int32) int32 {
	eng.Stats.Nodes++

	entry := eng.retrieveHash()
	if score := int32(entry.score); isInBounds(entry.kind, α, β, score) {
		return score
	}

	static := eng.cachedScore(&entry)
	if static >= β {
		// stand pat if the static score is already a cut-off
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
		// prune futile moves that would anyway result in a stand-pat at that next depth
		if !inCheck && isFutile(pos, static, α, futilityMargin, move) ||
			!inCheck && seeSign(pos, move) {
			continue
		}

		// discard illegal or losing captures
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

func initEngine() {
	var fens = [FigureArraySize]string{
		Pawn:   "rnbqkbnr/ppp1pppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR w - - 0 1",
		Knight: "r1bqkbnr/pppppppp/8/8/3N4/8/PPPPPPPP/R1BQKBNR w - - 0 1",
		Bishop: "rn1qkbnr/pppppppp/8/8/3B4/8/PPPPPPPP/RN1QKBNR w - - 0 1",
		Rook:   "rnbqkbn1/pppppppp/8/8/3R4/8/PPPPPPPP/RNBQKBN1 w - - 0 1",
		Queen:  "rnb1kbnr/pppppppp/8/8/3Q4/8/PPPPPPPP/RNB1KBNR w - - 0 1",
	}

	for f, fen := range fens {
		if fen != "" {
			pos, _ := PositionFromFEN(fen)
			futilityFigureBonus[f] = Evaluate(pos).GetCentipawnsScore()
		}
	}

	initialized = true
}

/////////////////////////////////////////////////////////////////////
