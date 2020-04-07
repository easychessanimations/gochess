package board

import "github.com/easychessanimations/gochess/utils"

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (ccr *ColorCastlingRights) Init(color utils.PieceColor, b *Board) {
	ccr[utils.KING_SIDE].Init(color, utils.KING_SIDE, b)
	ccr[utils.QUEEN_SIDE].Init(color, utils.QUEEN_SIDE, b)
}

func (ccr *ColorCastlingRights) SetFromFen(fen string, b *Board) {
	ccr[utils.KING_SIDE].SetFromFen(fen, b)
	ccr[utils.QUEEN_SIDE].SetFromFen(fen, b)
}

func (ccr *ColorCastlingRights) ToString(b *Board) string {
	return ccr[utils.KING_SIDE].ToString(b) + ccr[utils.QUEEN_SIDE].ToString(b)
}

func (ccr *ColorCastlingRights) ClearAll() {
	ccr[utils.KING_SIDE].Clear()
	ccr[utils.QUEEN_SIDE].Clear()
}

/////////////////////////////////////////////////////////////////////
