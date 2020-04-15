package bengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"

	. "github.com/easychessanimations/gochess/butils"
)

/////////////////////////////////////////////////////////////////////
// types

type pvTableAB struct {
	// uint64 is the type of Zobrist key
	PositionEntries map[uint64]Move
}

type AlphaBetaInfo struct {
	Alpha        int32
	Beta         int32
	MaxDepth     int32
	CurrentDepth int32
}

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (pvt pvTableAB) GetRec(pos *Position, moves []Move, remainingDepth int) []Move {
	if remainingDepth <= 0 {
		return moves
	}

	move, ok := pvt.PositionEntries[pos.Zobrist()]

	if ok {
		pos.DoMove(move)

		moves = pvt.GetRec(pos, append(moves, move), remainingDepth-1)

		pos.UndoMove()
	}

	return moves
}

func (pvt pvTableAB) Get(pos *Position) []Move {
	return pvt.GetRec(pos, []Move{}, 64)
}

func (eng *Engine) alphaBetaRec(abi AlphaBetaInfo) int32 {
	if eng.timeControl.Stopped() {
		eng.stopped = true
	}

	if eng.stopped {
		return abi.Alpha
	}

	eng.Stats.Nodes++
	eng.Stats.SelDepth = eng.Stats.Depth

	if abi.CurrentDepth >= abi.MaxDepth {
		return eng.Score()
	}

	lms := eng.Position.LegalMoves()

	if len(lms) == 0 {
		if eng.Position.IsChecked(eng.Position.Us()) {
			// mate
			return -MateScore + abi.CurrentDepth
		} else {
			// stalemate
			return 0
		}
	}

	pvMove, ok := eng.pvTableAB.PositionEntries[eng.Position.Zobrist()]

	// start search with pv move if any
	if ok {
		newLms := []Move{pvMove}

		for _, testMove := range lms {
			if testMove != pvMove {
				newLms = append(newLms, testMove)
			}
		}

		lms = newLms
	}

	for _, lm := range lms {
		ignored := false
		for _, irm := range eng.ignoreRootMoves {
			if lm == irm {
				ignored = true
				break
			}
		}

		if !ignored {
			newInfo := abi
			newInfo.Alpha = -abi.Beta
			newInfo.Beta = -abi.Alpha
			newInfo.CurrentDepth = abi.CurrentDepth + 1

			eng.Position.DoMove(lm)

			score := -eng.alphaBetaRec(newInfo)

			eng.Position.UndoMove()

			if score >= abi.Beta {
				// beta cut
				return abi.Beta
			}

			if score > abi.Alpha {
				// alpha improvement
				eng.pvTableAB.PositionEntries[eng.Position.Zobrist()] = lm

				abi.Alpha = score
			}
		}
	}

	return abi.Alpha
}

func (eng *Engine) searchAB(depth, estimated int32) int32 {
	if eng.pvTableAB.PositionEntries == nil {
		// initialize pv table
		eng.pvTableAB.PositionEntries = make(map[uint64]Move)
	} else if len(eng.pvTableAB.PositionEntries) > 1e7 {
		// reset pv table if grown too large
		eng.pvTableAB.PositionEntries = make(map[uint64]Move)
	}
	fmt.Println("info string position entries", len(eng.pvTableAB.PositionEntries))
	// delete root move
	delete(eng.pvTableAB.PositionEntries, eng.Position.Zobrist())
	return eng.alphaBetaRec(AlphaBetaInfo{
		Alpha:        -1e5,
		Beta:         1e5,
		MaxDepth:     depth,
		CurrentDepth: 0,
	})
}

/////////////////////////////////////////////////////////////////////
