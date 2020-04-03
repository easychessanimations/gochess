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

func (cr *CastlingRight) Init(color PieceColor, side CastlingSide, b *Board) {
	cr.Color = color
	cr.Side = side
	cr.Clear()
}

func (cr *CastlingRight) SetFromFen(fen string, b *Board) {
	cr.CanCastle = false

	for i := 0; i < len(fen); i++ {
		c := fen[i : i+1]

		if cr.Color == WHITE {
			if (c == "K") && (cr.Side == KING_SIDE) {
				cr.CanCastle = true
			}

			if (c == "Q") && (cr.Side == QUEEN_SIDE) {
				cr.CanCastle = true
			}
		} else {
			if (c == "k") && (cr.Side == KING_SIDE) {
				cr.CanCastle = true
			}

			if (c == "q") && (cr.Side == QUEEN_SIDE) {
				cr.CanCastle = true
			}
		}
	}
}

func (cr *CastlingRight) ToString(b *Board) string {
	if !cr.CanCastle {
		return ""
	}

	letter := "k"

	if cr.Side == QUEEN_SIDE {
		letter = "q"
	}

	if cr.Color == WHITE {
		letter = strings.ToUpper(letter)
	}

	return letter
}

/////////////////////////////////////////////////////////////////////
