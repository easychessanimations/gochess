package board

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (pos *Pos) Init() {
	pos.CastlingRights = CastlingRights{WHITE: SideCastlingRights{}, BLACK: SideCastlingRights{}}
}

func (pos *Pos) Clone() Pos {
	newPos := *pos

	crw, _ := pos.CastlingRights[WHITE]
	crb, _ := pos.CastlingRights[BLACK]

	newPos.CastlingRights = CastlingRights{
		WHITE: crw,
		BLACK: crb,
	}

	return newPos
}

/////////////////////////////////////////////////////////////////////
