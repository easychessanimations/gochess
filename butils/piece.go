package butils

import (
	"fmt"
	"strings"
)

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

// Figure.Symbol returns the symbol of the figure ( lower case )
func (f Figure) Symbol() string {
	return FigureToSymbol[f]
}

// Figure.SanSymbol returns the san symbol of the figure ( letter upper case )
func (f Figure) SanSymbol() string {
	symbol := f.Symbol()

	return fmt.Sprintf("%s%s", strings.ToUpper(symbol[0:1]), symbol[1:])
}

// Figure.IsLancer determines whether figure is lancer
func (f Figure) IsLancer() bool {
	return (f & LANCER_MASK) == Lancer
}

// Figure.LancerDirection determines the lancer direction of the figure
// figure has to be a lancer
func (lancer Figure) LancerDirection() int {
	return int(lancer & LANCER_DIRECTION_MASK)
}

// Figure.BaseFigure returns the base figure of the figure
// apart from lancer it is the same as figure
func (f Figure) BaseFigure() Figure {
	if f.IsLancer() {
		return Lancer
	}
	return f
}

// Figure.AlgebLetter returns the algeb letter of the figure ( lower case )
// for lancer direction is omitted
func (f Figure) AlgebLetter() string {
	return f.BaseFigure().Symbol()[0:1]
}

// Figure.SanLetter returns the san letter of the figure ( upper case )
// for lancer direction is omitted
func (f Figure) SanLetter() string {
	return fmt.Sprintf("%s", strings.ToUpper(f.AlgebLetter()))
}

// Color returns piece's color
func (pi Piece) Color() Color {
	if (pi == NoPiece) || (pi == DummyPiece) {
		return NoColor
	}
	return Color((pi & 1) + 1)
}

// Figure returns the figure of the piece
func (pi Piece) Figure() Figure {
	return Figure(pi) >> 1
}

// Piece.BaseFigure returns the base figure of the piece
// convenience function
func (pi Piece) BaseFigure() Figure {
	return pi.Figure().BaseFigure()
}

// AlgebSymbol returns the algebraic symbol of the piece
func (pi Piece) AlgebSymbol() string {
	return pi.Figure().Symbol()
}

// FenSymbol returns the fen symbol of the piece ( letter case depends on color )
func (pi Piece) FenSymbol() string {
	symbol := pi.AlgebSymbol()

	if pi.Color() == Black {
		return symbol
	}

	return fmt.Sprintf("%s%s", strings.ToUpper(symbol[0:1]), symbol[1:])
}

// PrettySymbol returns the pretty print symbol of the piece
func (pi Piece) PrettySymbol() string {
	symbol := pi.FenSymbol()

	if len(symbol) == 1 {
		return fmt.Sprintf(" %s ", symbol)
	}

	return fmt.Sprintf("%-3s", symbol)
}

// Piece.AlgebLetter returns the algeb letter of the piece ( upper case )
// convenience function
func (pi Piece) AlgebLetter() string {
	return pi.Figure().AlgebLetter()
}

// Piece.SanLetter returns the san letter of the piece ( lower case )
// convenience function
func (pi Piece) SanLetter() string {
	return pi.Figure().SanLetter()
}

// Piece.LancerDirection determines the lancer direction of the piece
// convenience function
func (lancer Piece) LancerDirection() int {
	return lancer.Figure().LancerDirection()
}

// Piece.SanSymbol returns the san symbol of the piece ( letter upper case )
// convenience function
func (pi Piece) SanSymbol() string {
	return pi.Figure().SanSymbol()
}

/////////////////////////////////////////////////////////////////////
