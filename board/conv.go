package board

import "github.com/easychessanimations/gochess/utils"

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// conv

func (b *Board) SquareToFileLetter(sq utils.Square) string {
	return string([]byte{"a"[0] + byte(sq.File)})
}

func (b *Board) SquareToRankLetter(sq utils.Square) string {
	return string([]byte{"1"[0] + byte(b.LastRank-sq.Rank)})
}

func (b *Board) SquareToAlgeb(sq utils.Square) string {
	if sq.File < 0 {
		return "-"
	}
	return b.SquareToFileLetter(sq) + b.SquareToRankLetter(sq)
}

func (b *Board) SquareFromAlgeb(algeb string) utils.Square {
	if algeb == "-" {
		return utils.NO_SQUARE
	}

	return utils.Square{int8(algeb[0] - "a"[0]), int8(byte(b.LastRank) - (algeb[1] - "1"[0]))}
}

func (b *Board) AlgebToMoveRaw(algeb string) utils.Move {
	if algeb == "-" {
		return NO_MOVE
	}

	fromsq := b.SquareFromAlgeb(algeb[0:2])
	tosq := b.SquareFromAlgeb(algeb[2:4])

	move := utils.Move{
		FromSq: fromsq,
		ToSq:   tosq,
	}

	return move
}

func (b *Board) AlgebToMove(algeb string) utils.Move {
	lms := b.LegalMovesForAllPieces()

	for _, lm := range lms {
		if b.MoveToAlgeb(lm) == algeb {
			return lm
		}
	}

	return NO_MOVE
}

func (b *Board) MoveToAlgeb(move utils.Move) string {
	if move == NO_MOVE {
		return "-"
	}

	buff := b.SquareToAlgeb(move.FromSq) + b.SquareToAlgeb(move.ToSq)

	if move.PromotionPiece != utils.NO_PIECE {
		buff += move.PromotionPiece.ToStringLower()

		if move.PromotionSquare != utils.NO_SQUARE {
			buff += "@" + b.SquareToAlgeb(move.PromotionSquare)
		}
	}

	return buff
}

func (b *Board) MoveToSan(move utils.Move) string {
	checkStr := ""

	b.Push(move, !ADD_SAN)
	check := b.IsInCheck(b.Pos.Turn)
	if check {
		checkStr = "+"
		if !b.HasLegalMoveColor(b.Pos.Turn) {
			checkStr = "#"
		}
	}
	b.Pop()

	if move.Castling {
		if move.CastlingSide == utils.QUEEN_SIDE {
			return "O-O-O" + checkStr
		}

		return "O-O" + checkStr
	}

	fromAlgeb := b.SquareToAlgeb(move.FromSq)
	toAlgeb := b.SquareToAlgeb(move.ToSq)
	fromPiece := b.PieceAtSquare(move.FromSq)
	pieceLetter := fromPiece.LetterUpper()

	qualifier := ""

	testPiece := fromPiece

	if fromPiece.Kind != utils.Pawn {
		oldPos := b.Pos

		if fromPiece.Kind == utils.Sentry {
			// for disambiguation need to replace all sentrys with bishop
			testPiece = utils.Piece{
				Kind:  utils.Bishop,
				Color: fromPiece.Color,
			}

			ssqs := b.SquaresForPiece(fromPiece)

			for _, sqs := range ssqs {
				b.SetPieceAtSquare(sqs, testPiece)
			}
		}

		pslAttacks := b.AttacksOnSquareByPiece(move.ToSq, testPiece, ALL_ATTACKS)

		if fromPiece.Kind == utils.Sentry {
			// put back sentrys
			b.Pos = oldPos
		}

		attacks := b.PickLegalMovesFrom(pslAttacks, b.Pos.Turn)

		files := make(map[int8]bool, 0)
		ranks := make(map[int8]bool, 0)
		samefiles := false
		sameranks := false

		if len(attacks) > 1 {
			for _, attack := range attacks {
				_, hasfile := files[attack.FromSq.File]
				if hasfile {
					samefiles = true
				} else {
					files[attack.FromSq.File] = true
				}

				_, hasrank := ranks[attack.FromSq.Rank]
				if hasrank {
					sameranks = true
				} else {
					ranks[attack.FromSq.Rank] = true
				}
			}

			if samefiles && sameranks {
				qualifier = fromAlgeb
			} else if samefiles {
				qualifier = fromAlgeb[1:2]
			} else {
				qualifier = fromAlgeb[0:1]
			}
		}
	}

	buff := pieceLetter + qualifier

	if fromPiece.Kind == utils.Pawn {
		buff = ""
	}

	if move.IsCapture() {
		if fromPiece.Kind == utils.Pawn {
			buff = b.SquareToFileLetter(move.FromSq)
		}
		buff += "x"
	}

	buff += toAlgeb

	if move.IsPromotion() {
		buff += "=" + move.PromotionPiece.ToStringUpper()

		if move.PromotionSquare != utils.NO_SQUARE {
			buff += "@" + b.SquareToAlgeb(move.PromotionSquare)
		}
	}

	return buff + checkStr
}

/////////////////////////////////////////////////////////////////////
