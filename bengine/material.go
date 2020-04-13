package bengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	. "github.com/easychessanimations/gochess/butils"
)

/////////////////////////////////////////////////////////////////////

// material.go implements position evaluation
//
// Zurichess' evaluation is a simple neural network with no hidden layers,
// and one output node y = W_m * x * (1-p) + W_e * x * p where W_m are
// middle game weights, W_e are endgame weights, x is input, p is phase between
// middle game and end game, and y is the score
// the network has |x| = len(Weights) inputs corresponding to features
// extracted from the position; these features are symmetrical wrt colors
// the network is trained using the Texel's Tuning Method
// https://chessprogramming.wikispaces.com/Texel%27s+Tuning+Method
// tuning is done by bitbucket.org/zurichess/tuner tool which uses
// tensorflow.org machine learning framework

/////////////////////////////////////////////////////////////////////
// material

// GetCentipawnsScore returns the current position evalution
// in centipawns
func (e Eval) GetCentipawnsScore() int32 {
	phase := Phase(e.position)
	score := (e.Accum[NoColor].M*(256-phase) + e.Accum[NoColor].E*phase) / 256
	return scaleToCentipawns(score)
}

// EvalExtra calculates additional Eval of extra pieces
func EvalExtra(pos *Position) Eval {
	e := Eval{position: pos}

	e.Accum[White] = evaluateExtra(pos, White)
	e.Accum[Black] = evaluateExtra(pos, Black)

	return e
}

// extra piece values
const LANCER_VALUE_NATIVE_M = 200000
const LANCER_VALUE_NATIVE_E = 200000
const SENTRY_VALUE_NATIVE_M = 100000
const SENTRY_VALUE_NATIVE_E = 100000
const JAILER_VALUE_NATIVE_M = 150000
const JAILER_VALUE_NATIVE_E = 150000

const BASE_LANCER_BONUS = 50000

// PawnStartRank tells the pawn start rank for a given color
func PawnStartRank(color Color) Bitboard {
	mask := BbWhiteSquares
	if color == Black {
		mask = BbBlackSquares
	}
	return BbPawnStartRank & mask
}

// evaluateExtra calculates additional Accum of extra pieces for a side
func evaluateExtra(pos *Position, us Color) Accum {
	var accum Accum

	for ld := 0; ld < NUM_LANCER_DIRECTIONS; ld++ {
		lancers := pos.ByPiece(us, MakeLancer(us, ld).Figure())

		numLancers := lancers.Count()

		accum.M += numLancers * LANCER_VALUE_NATIVE_M
		accum.E += numLancers * LANCER_VALUE_NATIVE_E

		baseLancers := lancers & PawnStartRank(us)

		for bb := baseLancers; bb != 0; {
			sq := bb.Pop()

			lancer := pos.Get(sq).Figure()

			if ((sq.File() <= 5) && (lancer == LancerE)) || ((sq.File() >= 3) && (lancer == LancerW)) {
				// add middle game bonus for lancer protecting their own second rank
				accum.M += BASE_LANCER_BONUS
			}
		}
	}

	sentries := pos.ByPiece(us, Sentry)

	numSentries := sentries.Count()

	accum.M += numSentries * LANCER_VALUE_NATIVE_M
	accum.E += numSentries * LANCER_VALUE_NATIVE_E

	jailers := pos.ByPiece(us, Jailer)

	numJailers := jailers.Count()

	accum.M += numJailers * LANCER_VALUE_NATIVE_M
	accum.E += numJailers * LANCER_VALUE_NATIVE_E

	return accum
}

// Evaluate evaluates the position pos
func Evaluate(pos *Position) Eval {
	e := Eval{position: pos}

	e.Accum[White] = evaluate(pos, White)
	e.Accum[Black] = evaluate(pos, Black)

	wps, bps := pawnsAndShelterCache.load(pos)
	e.Accum[White].merge(wps)
	e.Accum[Black].merge(bps)

	// extra
	ee := EvalExtra(pos)
	e.Accum[White].merge(ee.Accum[White])
	e.Accum[Black].merge(ee.Accum[Black])

	e.Accum[NoColor].merge(e.Accum[White])
	e.Accum[NoColor].deduct(e.Accum[Black])
	return e
}

// evaluatePawnsAndShelter evaluates pawns and shelter
func evaluatePawnsAndShelter(pos *Position, us Color) (accum Accum) {
	evaluatePawns(pos, us, &accum)
	evaluateShelter(pos, us, &accum)
	return accum
}

// evaluatePawns evaluates pawns
func evaluatePawns(pos *Position, us Color, accum *Accum) {
	groupBySquare(fPawnSquare, us, Pawns(pos, us), accum)
	groupByBoard(fBackwardPawns, BackwardPawns(pos, us), accum)
	groupByBoard(fConnectedPawns, ConnectedPawns(pos, us), accum)
	groupByBoard(fDoubledPawns, DoubledPawns(pos, us), accum)
	groupByBoard(fIsolatedPawns, IsolatedPawns(pos, us), accum)
	groupByBoard(fRammedPawns, RammedPawns(pos, us), accum)
	groupByRank(fPassedPawnRank, us, PassedPawns(pos, us), accum)
}

// evaluateShelter evaluates shelter
func evaluateShelter(pos *Position, us Color, accum *Accum) {
	// king's position and mobility
	bb := Kings(pos, us)
	kingSq := bb.AsSquare()
	mobility := KingMobility(kingSq)
	groupByFileSq(fKingFile, us, kingSq, accum)
	groupByRankSq(fKingRank, us, kingSq, accum)
	groupByBoard(fKingAttack, mobility, accum)

	// king's shelter
	ekw := East(bb) | bb | West(bb)
	ourPawns := Pawns(pos, us)
	groupByBoard(fKingShelterNear, ekw|Forward(us, ekw)&ourPawns, accum)
	groupByBoard(fKingShelterFar, ForwardSpan(us, ekw)&ourPawns, accum)
	groupByBoard(fKingShelterFront, ForwardSpan(us, bb)&ourPawns, accum)

	// king passed pawn tropism
	for bb := PassedPawns(pos, us); bb != BbEmpty; {
		if sq := bb.Pop(); sq.POV(us).Rank() >= 4 {
			dist := distance[sq][kingSq]
			groupByBucket(fKingPassedPawnTropism, int(dist), 8, accum)
		}
	}

	for bb := PassedPawns(pos, us.Opposite()); bb != BbEmpty; {
		if sq := bb.Pop(); sq.POV(us.Opposite()).Rank() >= 4 {
			dist := distance[sq][kingSq]
			groupByBucket(fKingEnemyPassedPawnTropism, int(dist), 8, accum)
		}
	}
}

// evaluate evaluates position for a single side
func evaluate(pos *Position, us Color) Accum {
	var accum Accum
	them := us.Opposite()
	all := pos.ByColor(White) | pos.ByColor(Black)
	danger := PawnThreats(pos, them)
	ourPawns := Pawns(pos, us)
	theirPawns := pos.ByPiece(them, Pawn)
	theirKingArea := KingArea(pos, them)

	groupByBoard(fNoFigure, BbEmpty, &accum)
	groupByBoard(fPawn, Pawns(pos, us), &accum)
	groupByBoard(fKnight, Knights(pos, us), &accum)
	groupByBoard(fBishop, Bishops(pos, us), &accum)
	groupByBoard(fRook, Rooks(pos, us), &accum)
	groupByBoard(fQueen, Queens(pos, us), &accum)
	groupByBoard(fKing, BbEmpty, &accum)

	// evaluate various pawn attacks and potential pawn attacks
	// on the enemy pieces
	groupByBoard(fPawnMobility, ourPawns&^Backward(us, all), &accum)
	groupByBoard(fMinorsPawnsAttack, Minors(pos, us)&danger, &accum)
	groupByBoard(fMajorsPawnsAttack, Majors(pos, us)&danger, &accum)
	groupByBoard(fMinorsPawnsPotentialAttack, Minors(pos, us)&Backward(us, danger), &accum)
	groupByBoard(fMajorsPawnsPotentialAttack, Majors(pos, us)&Backward(us, danger), &accum)

	numAttackers := 0
	attacks := PawnThreats(pos, us)

	// knight
	for bb := Knights(pos, us); bb > 0; {
		sq := bb.Pop()
		mobility := KnightMobility(sq) &^ (danger | ourPawns)
		attacks |= mobility
		groupByFileSq(fKnightFile, us, sq, &accum)
		groupByRankSq(fKnightRank, us, sq, &accum)
		groupByBoard(fKnightAttack, mobility, &accum)
		if mobility&theirKingArea&^theirPawns != 0 {
			numAttackers++
		}
	}
	// bishop
	// TODO: fix bishop's attack
	numBishops := 0
	for bb := Bishops(pos, us); bb > 0; {
		sq := bb.Pop()
		mobility := BishopMobility(sq, all)
		attacks |= mobility
		mobility &^= danger | ourPawns
		numBishops++
		groupByFileSq(fBishopFile, us, sq, &accum)
		groupByRankSq(fBishopRank, us, sq, &accum)
		groupByBoard(fBishopAttack, mobility, &accum)
		if mobility&theirKingArea&^theirPawns != 0 {
			numAttackers++
		}
	}
	// rook
	openFiles := OpenFiles(pos, us)
	semiOpenFiles := SemiOpenFiles(pos, us)
	for bb := Rooks(pos, us); bb > 0; {
		sq := bb.Pop()
		mobility := RookMobility(sq, all) &^ (danger | ourPawns)
		attacks |= mobility
		groupByFileSq(fRookFile, us, sq, &accum)
		groupByRankSq(fRookRank, us, sq, &accum)
		groupByBoard(fRookAttack, mobility, &accum)
		groupByBool(fRookOnOpenFile, openFiles.Has(sq), &accum)
		groupByBool(fRookOnSemiOpenFile, semiOpenFiles.Has(sq), &accum)
		if mobility&theirKingArea&^theirPawns != 0 {
			numAttackers++
		}
	}
	// queen
	for bb := Queens(pos, us); bb > 0; {
		sq := bb.Pop()
		mobility := QueenMobility(sq, all) &^ (danger | ourPawns)
		attacks |= mobility
		groupByFileSq(fQueenFile, us, sq, &accum)
		groupByRankSq(fQueenRank, us, sq, &accum)
		groupByBoard(fQueenAttack, mobility, &accum)
		if mobility&theirKingArea&^theirPawns != 0 {
			numAttackers++
		}

		dist := distance[sq][Kings(pos, them).AsSquare()]
		groupByCount(fKingQueenTropism, dist, &accum)
	}

	groupByBoard(fAttackedMinors, attacks&Minors(pos, them), &accum)
	groupByBool(fBishopPair, numBishops == 2, &accum)

	// king's safety is very primitive:
	// - king's shelter is evaluated by evaluateShelter
	// - the following counts the number of attackers
	// TODO: queen tropism which was dropped during the last refactoring
	groupByBucket(fKingAttackers, numAttackers, 4, &accum)
	return accum
}

// phase computes the progress of the game
// 0 is opening, 256 is late end game
func Phase(pos *Position) int32 {
	total := int32(4*1 + 4*1 + 4*3 + 2*6)
	curr := total
	curr -= pos.ByFigure(Knight).Count() * 1
	curr -= pos.ByFigure(Bishop).Count() * 1
	curr -= pos.ByFigure(Rook).Count() * 3
	curr -= pos.ByFigure(Queen).Count() * 6
	curr = max(curr, 0)
	return (curr*256 + total/2) / total
}

// scaleToCentipawns scales a score in the original scale to centipawns
func scaleToCentipawns(score int32) int32 {
	// divides by 128 and rounds to the nearest integer
	return (score + 128 + score>>31) >> 8
}

/////////////////////////////////////////////////////////////////////
