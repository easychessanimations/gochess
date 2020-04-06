package board

import "github.com/easychessanimations/gochess/utils"

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (cr *CastlingRights) Init(b *Board) {
	cr[utils.WHITE].Init(utils.WHITE, b)
	cr[utils.BLACK].Init(utils.BLACK, b)
}

func (cr *CastlingRights) SetFromFen(fen string, b *Board) {
	cr[utils.WHITE].SetFromFen(fen, b)
	cr[utils.BLACK].SetFromFen(fen, b)
}

func (cr *CastlingRights) ToString(b *Board) string {
	fen := cr[utils.WHITE].ToString(b) + cr[utils.BLACK].ToString(b)

	if fen == "" {
		return "-"
	}

	return fen
}

/////////////////////////////////////////////////////////////////////
