package butils

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

// Rank returns a number from 0 to 7 representing the rank of the square
func (sq Square) Rank() int {
	return int(sq / 8)
}

// File returns a number from 0 to 7 representing the file of the square
func (sq Square) File() int {
	return int(sq % 8)
}

// string representation of a Square
func (sq Square) String() string {
	squareToString := "a1b1c1d1e1f1g1h1a2b2c2d2e2f2g2h2a3b3c3d3e3f3g3h3a4b4c4d4e4f4g4h4a5b5c5d5e5f5g5h5a6b6c6d6e6f6g6h6a7b7c7d7e7f7g7h7a8b8c8d8e8f8g8h8"
	return squareToString[sq*2 : sq*2+2]
}

/////////////////////////////////////////////////////////////////////
