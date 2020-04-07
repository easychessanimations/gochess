package board

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
	"sort"
	"strings"

	"github.com/easychessanimations/gochess/utils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (b *Board) CreateMoveEvalBuff(moves []Move) MoveEvalBuff {
	meb := MoveEvalBuff{}

	pe, ok := b.PositionHash.PositionEntries[b.Pos]

	if ok {
		for _, move := range moves {
			eval := -INFINITE_SCORE

			me, ok := pe.MoveEntries[move]

			if ok {
				eval = me.Eval
			}

			meb = append(meb, MoveEvalBuffItem{
				Move: move,
				Eval: eval,
			})
		}

		sort.Sort(meb)
	} else {
		for _, move := range moves {
			eval := -INFINITE_SCORE

			if move.IsCapture() {
				eval += CAPTURE_BONUS
			}

			if !move.IsPawnMove() {
				eval += NON_PAWN_MOVE_BONUS
			}

			meb = append(meb, MoveEvalBuffItem{
				Move: move,
				Eval: eval,
			})
		}
	}

	return meb
}

func (b *Board) GetPv(maxDepth int) (string, []Move) {
	b.TestBoard = &Board{}

	b.TestBoard.Init(b.Variant)

	b.TestBoard.Pos = b.Pos

	b.TestBoard.PositionHash = b.PositionHash

	pv := []string{}

	pvMoves := []Move{}

	for i := 0; i < maxDepth; i++ {
		lms := b.TestBoard.LegalMovesForAllPieces()

		meb := b.TestBoard.CreateMoveEvalBuff(lms)

		if len(meb) > 0 {
			if meb[0].Eval > -INFINITE_SCORE {
				pvMove := meb[0].Move

				pv = append(pv, b.TestBoard.MoveToAlgeb(pvMove))

				pvMoves = append(pvMoves, pvMove)

				b.TestBoard.Push(pvMove, !ADD_SAN)
			} else {
				break
			}
		} else {
			break
		}
	}

	return strings.Join(pv, " "), pvMoves
}

// https://www.chessprogramming.org/Alpha-Beta
func (b *Board) AlphaBeta(info AlphaBetaInfo) (Move, int) {
	b.Nodes++

	if info.CurrentDepth > b.SelDepth {
		b.SelDepth = info.CurrentDepth
	}

	bm := Move{}

	if !b.Searching {
		return bm, b.EvalForTurn()
	}

	if info.CurrentDepth >= info.TotalDepth() {
		return bm, b.EvalForTurn()
	}

	plms := b.PslmsForAllPiecesOfColor(b.Pos.Turn)

	meb := b.CreateMoveEvalBuff(plms)

	isNormalSearch := info.CurrentDepth < info.Depth

	numLegals := 0

	for _, mebi := range meb {
		plm := mebi.Move

		if isNormalSearch || plm.IsCapture() {
			b.Push(plm, !ADD_SAN)

			if !b.IsInCheck(b.Pos.Turn.Inverse()) {
				numLegals++

				newInfo := info
				newInfo.Alpha = -info.Beta
				newInfo.Beta = -info.Alpha
				newInfo.CurrentDepth = info.CurrentDepth + 1
				newInfo.Line = append(newInfo.Line, b.MoveToAlgeb(plm))

				_, score := b.AlphaBeta(newInfo)

				b.Pop()

				score *= -1

				if score >= info.Beta {
					b.Betas++

					return bm, info.Beta
				}

				if score > info.Alpha {
					b.Alphas++

					pe := b.PositionHash.GetPositionEntry(b.Pos)

					pe.DecreaseEntries(500)

					me := pe.GetMoveEntry(plm)

					me.Eval = score

					pe.SetMoveEntry(plm, me)

					bm = plm
					info.Alpha = score
				}
			} else {
				b.Pop()
			}
		}
	}

	eval := b.EvalForTurn()

	if isNormalSearch {
		if numLegals <= 0 {
			if b.IsInCheck(b.Pos.Turn) {
				return bm, -(MATE_SCORE - info.CurrentDepth)
			} else {
				return bm, DRAW_SCORE
			}
		}
	} else {
		if numLegals <= 0 {
			return bm, eval
		}
	}

	if info.Alpha > -INFINITE_SCORE {
		return bm, info.Alpha
	}

	return bm, eval
}

func (b *Board) Stop() {
	b.Searching = false
}

func (b *Board) Go(depth int) (Move, int) {
	b.StartPerf()

	b.PositionHash = PositionHash{}

	b.PositionHash.Init()

	bm := Move{}

	score := -INFINITE_SCORE

	bestPv := ""

	pvMoves := []Move{}

	quiescenceDepthUciOption := b.GetUciOptionByNameWithDefault("Quiescence Depth", utils.UciOption{
		ValueInt: 0,
	})

	b.Log(fmt.Sprintf("go depth %d quiescence depth %d", depth, quiescenceDepthUciOption.ValueInt))

	for iterDepth := 1; iterDepth <= depth; iterDepth++ {
		b.SelDepth = 0

		alphaBetaInfo := AlphaBetaInfo{
			Alpha:           -INFINITE_SCORE,
			Beta:            INFINITE_SCORE,
			Depth:           iterDepth,
			QuiescenceDepth: quiescenceDepthUciOption.ValueInt,
			CurrentDepth:    0,
		}

		bm, score = b.AlphaBeta(alphaBetaInfo)

		if !b.Searching {
			break
		}

		nps, elapsed := b.GetNps()

		bestPv, pvMoves = b.GetPv(iterDepth)

		b.LogAnalysisInfo(fmt.Sprintf(
			"depth %d seldepth %d nodes %d time %.0f nps %.0f alphas %d betas %d score cp %d pv %s",
			iterDepth,
			b.SelDepth,
			b.Nodes,
			elapsed,
			nps,
			b.Alphas,
			b.Betas,
			score,
			bestPv,
		))
	}

	bestPvParts := strings.Split(bestPv, " ")

	if len(bestPvParts) > 1 {
		fmt.Println(fmt.Sprintf("bestmove %s ponder %s", bestPvParts[0], bestPvParts[1]))
	} else if bestPv != "" {
		fmt.Println(fmt.Sprintf("bestmove %s ponder null", bestPvParts[0]))
	} else {
		fmt.Println(fmt.Sprintf("bestmove null"))
	}

	b.Searching = false

	if len(pvMoves) > 0 {
		return pvMoves[0], score
	}

	return bm, score
}

/////////////////////////////////////////////////////////////////////
