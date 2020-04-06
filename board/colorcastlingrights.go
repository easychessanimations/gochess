package board

import "github.com/easychessanimations/gochess/utils"

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (ccr *ColorCastlingRights) Init(color utils.PieceColor, b *Board) {
	ccr[KING_SIDE].Init(color, KING_SIDE, b)
	ccr[QUEEN_SIDE].Init(color, QUEEN_SIDE, b)
}

func (ccr *ColorCastlingRights) SetFromFen(fen string, b *Board) {
	ccr[KING_SIDE].SetFromFen(fen, b)
	ccr[QUEEN_SIDE].SetFromFen(fen, b)
}

func (ccr *ColorCastlingRights) ToString(b *Board) string {
	return ccr[KING_SIDE].ToString(b) + ccr[QUEEN_SIDE].ToString(b)
}

func (ccr *ColorCastlingRights) ClearAll() {
	ccr[KING_SIDE].Clear()
	ccr[QUEEN_SIDE].Clear()
}

/////////////////////////////////////////////////////////////////////
