package butils

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// types

// Color represents a side
type Color uint

// Figure represents a piece without a color
type Figure uint

// Piece is a figure owned by one side
type Piece uint

// Bitboard is a set representing the 8x8 chess board squares
type Bitboard uint64

// square identifies the location on the board
type Square uint

// MoveType defines the move type
type MoveType uint8

// Move stores a position dependent move
//
// Bit representation
//   00.00.00.3f - from
//   00.00.3f.00 - to
//   00.0f.00.00 - move type
//   00.f0.00.00 - target
//   0f.00.00.00 - capture
//   f0.00.00.00 - piece
type Move uint32

// TODO: extend move to 64 bit, to allow for more move information

// Castle represents the castling rights mask
type Castle uint

// state
type state struct {
	Zobrist         uint64                    // Zobrist key, can be zero
	Move            Move                      // last move played
	HalfmoveClock   int                       // last ply when a pawn was moved or a capture was made
	EnpassantSquare Square                    // en passant square; if no e.p, then SquareA1
	CastlingAbility Castle                    // remaining castling rights
	ByFigure        [FigureArraySize]Bitboard // bitboards of square occupancy by figure
	ByColor         [ColorArraySize]Bitboard  // bitboards of square occupancy by color

	IsCheckedKnown   bool // true if it's known whether the current player is in check or not
	IsChecked        bool // true if current player is in check; if true then IsCheckedKnown is also true
	GivesCheckMove   Move // last move checkd with GivesCheck
	GivesCheckResult bool // true if last move gives check
}

// Position represents the chess board and keeps track of the move history
type Position struct {
	sideToMove Color // which side is to move. sideToMove is updated by DoMove and UndoMove
	Ply        int   // current ply

	pieces          [SquareArraySize]Piece // tracks pieces at each square
	fullmoveCounter int                    // fullmove counter, incremented after black move
	states          []state                // a state for each Ply
	curr            *state                 // current state
}

type castleInfo struct {
	Castle Castle
	Piece  [2]Piece
	Square [2]Square
}

/////////////////////////////////////////////////////////////////////
