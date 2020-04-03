package board

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (cr *CastlingRights) Init(b *Board) {
	cr[WHITE].Init(WHITE, b)
	cr[BLACK].Init(BLACK, b)
}

func (cr *CastlingRights) SetFromFen(fen string, b *Board) {
	cr[WHITE].SetFromFen(fen, b)
	cr[BLACK].SetFromFen(fen, b)
}

func (cr *CastlingRights) ToString(b *Board) string {
	fen := cr[WHITE].ToString(b) + cr[BLACK].ToString(b)

	if fen == "" {
		return "-"
	}

	return fen
}

/////////////////////////////////////////////////////////////////////
