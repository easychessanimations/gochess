package board

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
	return m.PromotionPiece != NO_PIECE
}

/////////////////////////////////////////////////////////////////////
