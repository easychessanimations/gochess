package board

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (cr *CastlingRights) SetFromFen(fen string) {
	K := false
	Q := false
	k := false
	q := false

	for i := 0; i < len(fen); i++ {
		c := fen[i : i+1]

		if c == "K" {
			K = true
		}

		if c == "Q" {
			Q = true
		}

		if c == "k" {
			k = true
		}

		if c == "q" {
			q = true
		}
	}

	(*cr)[WHITE] = SideCastlingRights{KingSide: K, QueenSide: Q}
	(*cr)[BLACK] = SideCastlingRights{KingSide: k, QueenSide: q}
}

func (cr *CastlingRights) ToString() string {
	buff := ""

	crw, _ := (*cr)[WHITE]
	crb, _ := (*cr)[BLACK]

	if crw.KingSide {
		buff += "K"
	}

	if crw.QueenSide {
		buff += "Q"
	}

	if crb.KingSide {
		buff += "k"
	}

	if crb.QueenSide {
		buff += "q"
	}

	if len(buff) > 0 {
		return buff
	}

	return "-"
}

/////////////////////////////////////////////////////////////////////
