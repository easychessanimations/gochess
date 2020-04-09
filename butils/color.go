package butils

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

// Forward returns bb shifted one rank forward wrt color
func Forward(col Color, bb Bitboard) Bitboard {
	if col == White {
		return bb << 8
	}
	if col == Black {
		return bb >> 8
	}
	return bb
}

// Backward returns bb shifted one rank backward wrt color
func Backward(col Color, bb Bitboard) Bitboard {
	if col == White {
		return bb >> 8
	}
	if col == Black {
		return bb << 8
	}
	return bb
}

// Opposite returns the reversed color
// result is undefined if c is not White or Black
func (c Color) Opposite() Color {
	return White + Black - c
}

// Multiplier returns -1 for Black, 1 for White
// useful for computing the position score relative to current player
// result is undefined if c is not White or Black
func (c Color) Multiplier() int32 {
	return int32(int(c)*2 - 3)
}

/////////////////////////////////////////////////////////////////////
