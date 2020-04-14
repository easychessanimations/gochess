package butils

import "fmt"

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

// POV returns the square from col's point of view
// that is for Black the rank is flipped, file stays the same
// useful in evaluation based on king's or pawns' positions
func (sq Square) POV(col Color) Square {
	return sq ^ (Square(col-2) & 0x38)
}

// Rank returns a number from 0 to 7 representing the rank of the square
func (sq Square) Rank() int {
	return int(sq / 8)
}

// File returns a number from 0 to 7 representing the file of the square
func (sq Square) File() int {
	return int(sq % 8)
}

// Bitboard returns a bitboard that has sq set
func (sq Square) Bitboard() Bitboard {
	return 1 << uint(sq)
}

// string representation of a Square
func (sq Square) String() string {
	squareToString := "a1b1c1d1e1f1g1h1a2b2c2d2e2f2g2h2a3b3c3d3e3f3g3h3a4b4c4d4e4f4g4h4a5b5c5d5e5f5g5h5a6b6c6d6e6f6g6h6a7b7c7d7e7f7g7h7a8b8c8d8e8f8g8h8"
	return squareToString[sq*2 : sq*2+2]
}

// AddDelta adds a delta to the square
func (sq Square) AddDelta(delta [2]int) (Square, error) {
	newRank := sq.Rank() + delta[0]
	newFile := sq.File() + delta[1]
	if newRank >= 0 && newRank <= 7 && newFile >= 0 && newFile <= 7 {
		return RankFile(newRank, newFile), nil
	}

	return SquareA1, fmt.Errorf("adding delta resulted in invalid square, square %v delta %v", sq, delta)
}

/////////////////////////////////////////////////////////////////////
