package butils

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

// Color returns piece's color
// 21844	=       101010101010100b

/*func (pi Piece) Color() Color {
	return Color(21844 >> pi & 3)
}*/

// Color returns piece's color
func (pi Piece) Color() Color {
	if (pi == NoPiece) || (pi == DummyPiece) {
		return NoColor
	}
	return Color((pi & 1) + 1)
}

// Figure returns piece's figure
func (pi Piece) Figure() Figure {
	return Figure(pi) >> 1
}

/////////////////////////////////////////////////////////////////////
