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
	PositionEntries map[uint64][]Move
}

type AlphaBetaInfo struct {
	MaxDepth     uint8
	CurrentDepth uint8
	Alpha        int32
	Beta         int32
}

type HashEntry struct {
	MaxDepth     uint8
	CurrentDepth uint8
	Score        int
}

var ht map[uint64]HashEntry

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (pvt pvTableAB) GetRec(pos *Position, moves []Move, remainingDepth int) []Move {
	if remainingDepth <= 0 {
		return moves
	}

	pvMoves, ok := pvt.PositionEntries[pos.Zobrist()]

	if ok {
		pos.DoMove(pvMoves[0])

		moves = pvt.GetRec(pos, append(moves, pvMoves[0]), remainingDepth-1)

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

	for i := abi.MaxDepth; len(ht) > 1e7 && i >= 0; i-- {
		fmt.Println("info string clear hash at ply", i)
		for k, e := range ht {
			if e.CurrentDepth >= (i - 2) {
				delete(ht, k)
			}
		}
	}

	he, ok := ht[eng.Position.Zobrist()]

	if ok {
		if he.CurrentDepth <= abi.CurrentDepth && he.MaxDepth >= abi.MaxDepth {
			return int32(he.Score)
		}
	}

	if abi.CurrentDepth >= abi.MaxDepth {
		score := eng.Score()

		newHe := HashEntry{
			CurrentDepth: abi.CurrentDepth,
			MaxDepth:     abi.MaxDepth,
			Score:        int(score),
		}

		if ok {
			if abi.CurrentDepth <= he.CurrentDepth && abi.MaxDepth >= he.MaxDepth {
				ht[eng.Position.Zobrist()] = newHe
			}
		} else {
			ht[eng.Position.Zobrist()] = newHe
		}

		return score
	}

	lms := eng.Position.LegalMoves()

	if len(lms) == 0 {
		if eng.Position.IsChecked(eng.Position.Us()) {
			// mate
			return -MateScore + int32(abi.CurrentDepth)
		} else {
			// stalemate
			return 0
		}
	}

	pvMoves, ok := eng.pvTableAB.PositionEntries[eng.Position.Zobrist()]

	// start search with pv move if any
	if ok {
		newLms := pvMoves

		for _, testMove := range lms {
			found := false
			for _, testPvMove := range pvMoves {
				if testMove == testPvMove {
					found = true
					break
				}
			}
			if !found {
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
				pvMoves, ok := eng.pvTableAB.PositionEntries[eng.Position.Zobrist()]

				if ok {
					// sort entry with lm first
					newPvMoves := []Move{lm}
					for _, testPvMove := range pvMoves {
						if testPvMove != lm {
							newPvMoves = append(newPvMoves, testPvMove)
						}
					}

					eng.pvTableAB.PositionEntries[eng.Position.Zobrist()] = newPvMoves
				} else {
					// add entry with single move lm
					eng.pvTableAB.PositionEntries[eng.Position.Zobrist()] = []Move{lm}
				}

				abi.Alpha = score
			}
		}
	}

	return abi.Alpha
}

func (eng *Engine) searchAB(depth, estimated int32) int32 {
	if eng.pvTableAB.PositionEntries == nil {
		// initialize pv table
		eng.pvTableAB.PositionEntries = make(map[uint64][]Move)
	} else if len(eng.pvTableAB.PositionEntries) > 1e7 {
		// reset pv table if grown too large
		eng.pvTableAB.PositionEntries = make(map[uint64][]Move)
	}
	// delete root move
	delete(eng.pvTableAB.PositionEntries, eng.Position.Zobrist())
	/*moveCount := 0
	for _, entry := range eng.pvTableAB.PositionEntries {
		moveCount += len(entry)
	}
	fmt.Println("info string position entries", len(eng.pvTableAB.PositionEntries), "moves", moveCount)	*/

	ht = make(map[uint64]HashEntry)

	return eng.alphaBetaRec(AlphaBetaInfo{
		Alpha:        -1e5,
		Beta:         1e5,
		MaxDepth:     uint8(depth),
		CurrentDepth: 0,
	})
}

/////////////////////////////////////////////////////////////////////
