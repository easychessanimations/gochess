package board

import "github.com/easychessanimations/gochess/utils"

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (m *Move) IsCapture() bool {
	return m.Capture || m.PawnCapture || m.EpCapture
}

func (m *Move) IsPawnMove() bool {
	return m.PawnPushByOne || m.PawnPushByTwo || m.PawnCapture || m.EpCapture
}

func (m *Move) ShouldDeleteHalfmoveClock() bool {
	return m.IsCapture() || m.IsPawnMove()
}

func (m *Move) IsPromotion() bool {
	return m.PromotionPiece != utils.NO_PIECE
}

func (m *Move) RoughlyEqualTo(testm Move) bool {
	return m.FromSq.EqualTo(testm.FromSq) && m.ToSq.EqualTo(testm.ToSq)
}

/////////////////////////////////////////////////////////////////////
