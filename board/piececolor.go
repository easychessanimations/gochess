package board

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (p *PieceColor) ToString() string {
	if *p {
		return "w"
	}

	return "b"
}

func (p *PieceColor) SetFromFen(fen string) {
	if fen == "w" {
		*p = WHITE
		return
	}

	*p = BLACK
}

/////////////////////////////////////////////////////////////////////
