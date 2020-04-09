package butils

/////////////////////////////////////////////////////////////////////
// imports

import (
	"math/bits"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

// RankBb returns a bitboard with all bits on rank set
func RankBb(rank int) Bitboard {
	return BbRank1 << uint(8*rank)
}

// FileBb returns a bitboard with all bits on file set
func FileBb(file int) Bitboard {
	return BbFileA << uint(file)
}

// North shifts all squares one rank up
func North(bb Bitboard) Bitboard {
	return bb << 8
}

// South shifts all squares one rank down
func South(bb Bitboard) Bitboard {
	return bb >> 8
}

// meaning of &^ operator
// https://stackoverflow.com/questions/34459450/what-is-the-operator-in-golang
// turn mask bits off

// East shifts all squares one file right
// delete h-file, then shift left
func East(bb Bitboard) Bitboard {
	return bb &^ BbFileH << 1
}

// West shifts all squares one file left
// delete a-file, then shift right
func West(bb Bitboard) Bitboard {
	return bb &^ BbFileA >> 1
}

// Fill returns a bitboard with all files with squares filled.
func Fill(bb Bitboard) Bitboard {
	return NorthFill(bb) | SouthFill(bb)
}

// ForwardSpan computes forward span wrt color.
func ForwardSpan(col Color, bb Bitboard) Bitboard {
	if col == White {
		return NorthSpan(bb)
	}
	if col == Black {
		return SouthSpan(bb)
	}
	return bb
}

// ForwardFill computes forward fill wrt color.
func ForwardFill(col Color, bb Bitboard) Bitboard {
	if col == White {
		return NorthFill(bb)
	}
	if col == Black {
		return SouthFill(bb)
	}
	return bb
}

// BackwardSpan computes backward span wrt color.
func BackwardSpan(col Color, bb Bitboard) Bitboard {
	if col == White {
		return SouthSpan(bb)
	}
	if col == Black {
		return NorthSpan(bb)
	}
	return bb
}

// BackwardFill computes forward fill wrt color.
func BackwardFill(col Color, bb Bitboard) Bitboard {
	if col == White {
		return SouthFill(bb)
	}
	if col == Black {
		return NorthFill(bb)
	}
	return bb
}

// NorthFill returns a bitboard with all north bits set.
func NorthFill(bb Bitboard) Bitboard {
	bb |= (bb << 8)
	bb |= (bb << 16)
	bb |= (bb << 32)
	return bb
}

// NorthSpan is like NorthFill shifted on up.
func NorthSpan(bb Bitboard) Bitboard {
	return NorthFill(North(bb))
}

// SouthFill returns a bitboard with all south bits set.
func SouthFill(bb Bitboard) Bitboard {
	bb |= (bb >> 8)
	bb |= (bb >> 16)
	bb |= (bb >> 32)
	return bb
}

// SouthSpan is like SouthFill shifted on up.
func SouthSpan(bb Bitboard) Bitboard {
	return SouthFill(South(bb))
}

// Has returns bb if sq is occupied in bitboard.
func (bb Bitboard) Has(sq Square) bool {
	return bb>>sq&1 != 0
}

// AsSquare returns the occupied square if the bitboard has a single piece.
// If the board has more then one piece the result is undefined.
// https://golang.org/pkg/math/bits/#TrailingZeros64
func (bb Bitboard) AsSquare() Square {
	return Square(bits.TrailingZeros64(uint64(bb)) & 0x3f)
}

// LSB picks a square in the board.
// Returns empty board for empty board.
func (bb Bitboard) LSB() Bitboard {
	return bb & (-bb)
}

// Count returns the number of squares set in bb.
// https://golang.org/pkg/math/bits/#OnesCount64
func (bb Bitboard) Count() int32 {
	return int32(bits.OnesCount64(uint64(bb)))
}

// Pop pops a set square from the bitboard.
func (bb *Bitboard) Pop() Square {
	sq := *bb & (-*bb)
	*bb -= sq
	return Square(bits.TrailingZeros64(uint64(sq)) & 0x3f)
}

/////////////////////////////////////////////////////////////////////
