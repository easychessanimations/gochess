package butils

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

// From returns the starting square
func (m Move) From() Square {
	return Square(m >> 0 & 0x3f)
}

// To returns the destination square
func (m Move) To() Square {
	return Square(m >> 8 & 0x3f)
}

// MoveType returns the move type
func (m Move) MoveType() MoveType {
	return MoveType(m >> 16 & 0xf)
}

// CaptureSquare returns the captured piece square
// if no piece is captured, the result is the destination square
func (m Move) CaptureSquare() Square {
	if m.MoveType() != Enpassant {
		return m.To()
	}
	return m.From()&0x38 + m.To()&0x7
}

// Capture returns the captured pieces
func (m Move) Capture() Piece {
	return Piece((m >> MOVE_CAPTURE_SHIFT) & Move(PIECE_MASK))
}

// Target returns the piece on the to square after the move is executed
func (m Move) Target() Piece {
	return Piece((m >> MOVE_TARGET_SHIFT) & Move(PIECE_MASK))
}

// Piece returns the piece moved
func (m Move) Piece() Piece {
	return Piece((m >> MOVE_PIECE_SHIFT) & Move(PIECE_MASK))
}

// Color returns which player is moving
func (m Move) Color() Color {
	return m.Piece().Color()
}

// Figure returns which figure is moved
func (m Move) Figure() Figure {
	return m.Piece().Figure()
}

// Promotion returns the promoted piece if any
func (m Move) Promotion() Piece {
	if m.MoveType() != Promotion {
		return NoPiece
	}
	return m.Target()
}

// IsViolent returns true if the move is a capture or a queen promotion
// castling and minor promotions (including captures) are not violent
// TODO: IsViolent should be in sync with GenerateViolentMoves
func (m Move) IsViolent() bool {
	if m.MoveType() != Promotion {
		return m.Capture() != NoPiece
	}
	return m.Target().Figure() == Queen
}

// IsQuiet returns true if the move is not violent
func (m Move) IsQuiet() bool {
	return !m.IsViolent()
}

// UCI converts a move to UCI format
// the protocol specification at http://wbec-ridderkerk.nl/html/UCIProtocol.html
// incorrectly states that this is the long algebraic notation (LAN)
func (m Move) UCI() string {
	buff := m.From().String() + m.To().String()
	promFigure := m.Promotion().Figure()
	if promFigure != NoFigure {
		buff += promFigure.Symbol()
	}
	return buff
}

// LAN converts a move to Long Algebraic Notation
// http://en.wikipedia.org/wiki/Algebraic_notation_%28chess%29#Long_algebraic_notation
// e.g. a2-a3, b7-b8Q, Nb1xc3, Ke1-c1 (white king queen side castling)
func (m Move) LAN() string {
	r := m.Piece().SanLetter() + m.From().String()
	if m.Capture() != NoPiece {
		r += "x"
	} else {
		r += "-"
	}
	r += m.To().String()
	promPiece := m.Promotion()
	if promPiece != NoPiece {
		r += promPiece.SanSymbol()
	}
	return r
}

// move as a string
func (m Move) String() string {
	return m.LAN()
}

/////////////////////////////////////////////////////////////////////
