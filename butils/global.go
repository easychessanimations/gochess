package butils

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// global functions

// ColorFigure returns a piece with col and fig
func ColorFigure(col Color, fig Figure) Piece {
	return Piece(fig<<1) + Piece(col>>1)
}

// RankFile returns a square with rank r and file f
// r and f should be between 0 and 7
func RankFile(r, f int) Square {
	return Square(r*8 + f)
}

// MakeMove constructs a move
func MakeMove(moveType MoveType, from, to Square, capture, target Piece) Move {
	piece := target
	if moveType == Promotion {
		piece = ColorFigure(target.Color(), Pawn)
	}

	return Move(from)<<0 +
		Move(to)<<8 +
		Move(moveType)<<16 +
		Move(target)<<20 +
		Move(capture)<<24 +
		Move(piece)<<28
}

// SquareFromString parses a square from a string
// the string has standard chess format [a-h][1-8]
func SquareFromString(s string) (Square, error) {
	if len(s) != 2 {
		return SquareA1, fmt.Errorf("invalid square %s", s)
	}

	f, r := -1, -1
	if 'a' <= s[0] && s[0] <= 'h' {
		f = int(s[0] - 'a')
	}
	if 'A' <= s[0] && s[0] <= 'H' {
		f = int(s[0] - 'A')
	}
	if '1' <= s[1] && s[1] <= '8' {
		r = int(s[1] - '1')
	}
	if f == -1 || r == -1 {
		return SquareA1, fmt.Errorf("invalid square %s", s)
	}

	return RankFile(r, f), nil
}

// CastlingRook returns the rook moved during castling
// together with starting and stopping squares
func CastlingRook(kingEnd Square) (Piece, Square, Square) {
	// Explanation how rookStart works for king on E1.
	// if kingEnd == C1 == b010, then rookStart == A1 == b000
	// if kingEnd == G1 == b110, then rookStart == H1 == b111
	// So bit 3 will set bit 2 and bit 1.
	//
	// Explanation how rookEnd works for king on E1.
	// if kingEnd == C1 == b010, then rookEnd == D1 == b011
	// if kingEnd == G1 == b110, then rookEnd == F1 == b101
	// So bit 3 will invert bit 2. bit 1 is always set.
	piece := Piece(Rook<<1) + (1 - Piece(kingEnd>>5))
	rookStart := kingEnd&^3 | (kingEnd & 4 >> 1) | (kingEnd & 4 >> 2)
	rookEnd := kingEnd ^ (kingEnd & 4 >> 1) | 1
	return piece, rookStart, rookEnd
}

/////////////////////////////////////////////////////////////////////
