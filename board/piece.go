package board

/////////////////////////////////////////////////////////////////////
// imports

import "strings"

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (p *Piece) ToString() string {
	letter, _ := PIECE_KIND_TO_PIECE_LETTER[p.Kind]

	if p.Color {
		letter = strings.ToUpper(letter)
	}

	return letter + p.Direction.ToString()
}

func (p *Piece) ToStringUpper() string {
	str := p.ToString()
	return strings.ToUpper(str[0:1]) + str[1:]
}

func (p *Piece) IsEmpty() bool {
	return p.Kind == NO_PIECE
}

/////////////////////////////////////////////////////////////////////
