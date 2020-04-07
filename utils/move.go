package utils

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (m *Move) IsCapture() bool {
	return m.Capture || m.PawnCapture || m.EpCapture || m.SentryPush
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

func (m *Move) RoughlyEqualTo(testm Move) bool {
	return m.FromSq.EqualTo(testm.FromSq) && m.ToSq.EqualTo(testm.ToSq)
}

func (m *Move) EffectivePromotionSquare() Square {
	if m.PromotionSquare == NO_SQUARE {
		return m.ToSq
	}

	return m.PromotionSquare
}

func (m *Move) NormalizedDirection() PieceDirection {
	fileDiff := m.ToSq.File - m.FromSq.File
	rankDiff := m.ToSq.Rank - m.FromSq.Rank

	if fileDiff == 0 {
		if rankDiff > 0 {
			return PieceDirection{0, 1}
		} else if rankDiff == 0 {
			return PieceDirection{0, 0}
		} else {
			return PieceDirection{0, -1}
		}
	}

	if rankDiff == 0 {
		if fileDiff > 0 {
			return PieceDirection{1, 0}
		} else {
			return PieceDirection{-1, 0}
		}
	}

	// non diagonal direction cannot be normalized
	if (fileDiff * fileDiff) != (rankDiff * rankDiff) {
		return PieceDirection{0, 0}
	}

	if fileDiff > 0 {
		if rankDiff > 0 {
			return PieceDirection{1, 1}
		} else {
			return PieceDirection{1, -1}
		}
	} else {
		if rankDiff > 0 {
			return PieceDirection{-1, 1}
		} else {
			return PieceDirection{-1, -1}
		}
	}
}

/////////////////////////////////////////////////////////////////////
