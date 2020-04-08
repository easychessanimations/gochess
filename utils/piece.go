package utils

/////////////////////////////////////////////////////////////////////
// imports

import "strings"

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (p *Piece) Letter() string {
	letter, _ := PIECE_KIND_TO_PIECE_LETTER[p.Kind]

	if p.Color == WHITE {
		letter = strings.ToUpper(letter)
	}

	return letter
}

func (p *Piece) LetterUpper() string {
	return strings.ToUpper(p.Letter())
}

func (p *Piece) LetterLower() string {
	return strings.ToLower(p.Letter())
}

func (p *Piece) ToString() string {
	return p.Letter() + p.Direction.ToString()
}

func (p *Piece) ToStringUpper() string {
	str := p.ToString()
	return strings.ToUpper(str[0:1]) + str[1:]
}

func (p *Piece) ToStringLower() string {
	str := p.ToString()
	return strings.ToLower(str)
}

func (p *Piece) KindColorEqualTo(testP Piece) bool {
	return (p.Kind == testP.Kind) && (p.Color == testP.Color)
}

func (p *Piece) EqualTo(testP Piece) bool {
	return p.KindColorEqualTo(testP) && (p.Direction == testP.Direction)
}

func (p *Piece) ColorInverse() Piece {
	return Piece{Kind: p.Kind, Color: p.Color.Inverse(), Direction: p.Direction}
}

/////////////////////////////////////////////////////////////////////
