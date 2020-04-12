package butils

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

// lancer attacks provides masks for all possible lancer directions
// for all possible squares

/////////////////////////////////////////////////////////////////////
// constants

const NUM_LANCER_DIRECTIONS = 8

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// types

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// global variables

// LANCER_DIRECTION_TO_DELTA maps lancer directions to corresponding deltas
var LANCER_DIRECTION_TO_DELTA = [NUM_LANCER_DIRECTIONS][2]int{
	{1, 0},   // n
	{1, 1},   // ne
	{0, 1},   // e
	{-1, 1},  // se
	{-1, 0},  // s
	{-1, -1}, // sw
	{0, -1},  // w
	{1, -1},  // nw
}

// LancerDirectionMasksForSquares maps squares and lancer directions to lancer attack masks
var LancerDirectionMasksForSquares [SquareArraySize][NUM_LANCER_DIRECTIONS]Bitboard

const BaseLancer = LancerMinValue

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// global functions

// MakeLancer creates a lancer piece from color and direction
func MakeLancer(color Color, direction int) Piece {
	figure := BaseLancer | Figure(direction)

	return ColorFigure(color, figure)
}

// MakeLancerMove constructs a lancer move
func MakeLancerMove(from, to Square, piece, capture, target Piece) Move {
	return Move(from)<<MOVE_FROM_SHIFT +
		Move(to)<<MOVE_TO_SHIFT +
		Move(Promotion)<<MOVE_TYPE_SHIFT +
		Move(target)<<MOVE_TARGET_SHIFT +
		Move(capture)<<MOVE_CAPTURE_SHIFT +
		Move(piece)<<MOVE_PIECE_SHIFT
}

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// init

func init() {
	for sq := SquareMinValue; sq <= SquareMaxValue; sq++ {
		for ldi := 0; ldi < NUM_LANCER_DIRECTIONS; ldi++ {
			delta := LANCER_DIRECTION_TO_DELTA[ldi]

			mask := slidingAttack(sq, [][2]int{delta}, BbEmpty)

			LancerDirectionMasksForSquares[sq][ldi] = mask
		}
	}
}

/////////////////////////////////////////////////////////////////////
