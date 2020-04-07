package board

import (
	"math/rand"
	"strconv"

	"github.com/easychessanimations/gochess/utils"
)

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (b *Board) MakeAlgebMove(algeb string, addSan bool) {
	move := b.AlgebToMove(algeb)

	if move != NO_MOVE {
		b.Push(move, addSan)
	}
}

func (b *Board) Pop() {
	l := len(b.MoveStack)
	if l == 0 {
		return
	}

	msi := b.MoveStack[l-1]

	b.Pos = msi.Pos

	b.MoveStack = b.MoveStack[:l-1]
}

func (b *Board) Push(move utils.Move, addSan bool) {
	san := "?"

	if addSan {
		san = b.MoveToSan(move)
	}

	oldPos := b.Pos.Clone()

	//////////////////////////////////////////////

	b.Pos.DisabledMove = NO_MOVE

	fromp := b.PieceAtSquare(move.FromSq)

	ccr := &b.Pos.CastlingRights[b.Pos.Turn]

	if fromp.Kind == utils.King {
		ccr.ClearAll()
	}

	b.SetPieceAtSquare(move.FromSq, utils.NO_PIECE)

	if move.IsPromotion() {
		b.SetPieceAtSquare(move.EffectivePromotionSquare(), move.PromotionPiece)
	}

	if move.SentryPush {
		disabledMove := utils.Move{
			FromSq: move.PromotionSquare,
			ToSq:   move.ToSq,
		}

		b.Pos.DisabledMove = disabledMove
	}

	if move.EpCapture {
		b.SetPieceAtSquare(move.EpClearSquare, utils.NO_PIECE)
	}

	if move.Castling {
		b.SetPieceAtSquare(move.ToSq, utils.NO_PIECE)
		kctsq := b.KingCastlingTargetSq(b.Pos.Turn, move.CastlingSide)
		b.SetPieceAtSquare(kctsq, utils.Piece{Kind: utils.King, Color: b.Pos.Turn})
		rctsq := b.RookCastlingTargetSq(b.Pos.Turn, move.CastlingSide)
		b.SetPieceAtSquare(rctsq, move.RookOrigPiece)
	} else {
		b.SetPieceAtSquare(move.ToSq, fromp)
	}

	if b.IS_ATOMIC() {
		if move.IsCapture() {
			// atomic explosion
			b.SetPieceAtSquare(move.ToSq, utils.NO_PIECE)

			for _, sq := range b.AdjacentSquares(move.ToSq) {
				p := b.PieceAtSquare(sq)

				if p.Kind != utils.Pawn {
					b.SetPieceAtSquare(sq, utils.NO_PIECE)
				}
			}
		}
	}

	var side utils.CastlingSide
	for side = utils.QUEEN_SIDE; side <= utils.KING_SIDE; side++ {
		cs := &ccr[side]
		if cs.CanCastle {
			rp := b.PieceAtSquare(cs.RookOrigSquare)

			if !cs.RookOrigPiece.KindColorEqualTo(rp) {
				// rook changed, delete castling right
				cs.CanCastle = false
			}
		}
	}

	b.Pos.Turn = b.Pos.Turn.Inverse()

	b.Pos.EpSquare = utils.NO_SQUARE

	if move.PawnPushByTwo {
		b.Pos.EpSquare = move.EpSquare
	}

	if move.ShouldDeleteHalfmoveClock() {
		b.Pos.HalfmoveClock = 0
	} else {
		b.Pos.HalfmoveClock++
	}

	if b.Pos.Turn == utils.WHITE {
		b.Pos.FullmoveNumber++
	}

	b.MoveStack = append(b.MoveStack, MoveStackItem{
		oldPos,
		move,
		san,
	})
}

func (b *Board) ExecCommand(command string) bool {
	b.SortedSanMoveBuff = b.MovesSortedBySan(b.LegalMovesForAllPieces())

	i, err := strconv.ParseInt(command, 10, 32)

	if err == nil {
		move := b.SortedSanMoveBuff[i-1].Move

		b.Push(move, ADD_SAN)

		b.Print()

		return true
	} else {
		if command == "g" {
			bm, _ := b.Go(10)

			b.Push(bm, ADD_SAN)

			b.Print()

			return true
		} else if command == "perf" {
			b.Perf(3)

			return true
		} else if command == "d" {
			b.Pop()

			b.Print()

			return true
		} else if command == "" {
			randIndex := rand.Intn(len(b.SortedSanMoveBuff))

			move := b.SortedSanMoveBuff[randIndex-1].Move

			b.Push(move, ADD_SAN)

			b.Print()

			return true
		} else if command != "" {
			for _, mbi := range b.SortedSanMoveBuff {
				if (mbi.San == command) || (mbi.Algeb == command) {
					move := mbi.Move

					b.Push(move, ADD_SAN)

					b.Print()

					return true
				}
			}
		}
	}

	return false
}

/////////////////////////////////////////////////////////////////////
