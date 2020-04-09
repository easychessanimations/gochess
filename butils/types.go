package butils

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// types

// Color represents a side.
type Color uint

// Figure represents a piece without a color
type Figure uint

// Piece is a figure owned by one side.
type Piece uint

// Bitboard is a set representing the 8x8 chess board squares
type Bitboard uint64

// square identifies the location on the board
type Square uint

// MoveType defines the move type.
type MoveType uint8

// Move stores a position dependent move.
//
// Bit representation
//   00.00.00.3f - from
//   00.00.3f.00 - to
//   00.0f.00.00 - move type
//   00.f0.00.00 - target
//   0f.00.00.00 - capture
//   f0.00.00.00 - piece
type Move uint32

// TODO: extend move to 64 bit, to allow for more move information (like promotion square)?

// Castle represents the castling rights mask
type Castle uint

/////////////////////////////////////////////////////////////////////
