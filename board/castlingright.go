package board

/////////////////////////////////////////////////////////////////////
// imports

import "strings"

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (cr *CastlingRight) Clear() {
	cr.CanCastle = false
}

// https://en.wikipedia.org/wiki/X-FEN

func (cr *CastlingRight) Init(color PieceColor, side CastlingSide, b *Board) {
	cr.Color = color
	cr.Side = side
	cr.Clear()
}

func (cr *CastlingRight) SetFromFen(fen string, b *Board) {
	cr.CanCastle = false

	wk := b.WhereIsKing(cr.Color)

	// sanity check, no king, no castling
	if wk == NO_SQUARE {
		return
	}

	// sanity check, king should be on castling rank
	if wk.Rank != b.CastlingRank(cr.Color) {
		return
	}

	sqs := b.SquaresInDirection(wk, PieceDirection{int8(cr.Side*2) - 1, 0})

	lastSq := NO_SQUARE
	lastKind := NO_PIECE

	for _, sq := range sqs {
		p := b.PieceAtSquare(sq)

		if ((p.Kind == Rook) || (p.Kind == Jailer)) && (p.Color == cr.Color) {
			letter := b.SquareToFileLetter(sq)

			if cr.Color == WHITE {
				letter = strings.ToUpper(letter)
			}

			if strings.Contains(fen, letter) {
				// x-fen
				cr.RookOrigSquare = sq
				cr.RookOrigPiece = p
				cr.CanCastle = true
				return
			} else {
				lastSq = sq
				lastKind = p.Kind
			}
		}
	}

	if lastKind == NO_PIECE {
		// no rook, no castling
		return
	}

	// fall back to conventional fen with outermost rook
	if strings.Contains(fen, CastlingLetter(cr.Color, cr.Side)) {
		cr.RookOrigSquare = lastSq
		cr.RookOrigPiece = Piece{Kind: lastKind, Color: cr.Color}
		cr.CanCastle = true
	}
}

func (cr *CastlingRight) ToString(b *Board) string {
	if !cr.CanCastle {
		return ""
	}

	wk := b.WhereIsKing(cr.Color)

	// sanity check, no king, no castling
	if wk == NO_SQUARE {
		return ""
	}

	// sanity check, king should be on castling rank
	if wk.Rank != b.CastlingRank(cr.Color) {
		return ""
	}

	sqs := b.SquaresInDirection(wk, PieceDirection{int8(cr.Side*2) - 1, 0})

	rcnt := 0

	// coount rooks
	for _, sq := range sqs {
		p := b.PieceAtSquare(sq)

		if ((p.Kind == Rook) || (p.Kind == Jailer)) && (p.Color == cr.Color) {
			rcnt++

			if rcnt > 1 {
				break
			}
		}
	}

	if rcnt == 0 {
		// no rook, no castling
		return ""
	} else if rcnt > 1 {
		// more than one rook, needs x-fen
		letter := b.SquareToFileLetter(cr.RookOrigSquare)
		if cr.Color == WHITE {
			letter = strings.ToUpper(letter)
		}
		return letter
	} else {
		// conventional fen enough
		return CastlingLetter(cr.Color, cr.Side)
	}
}

/////////////////////////////////////////////////////////////////////
