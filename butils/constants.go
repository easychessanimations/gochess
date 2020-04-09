package butils

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// constants

const (
	// Set of possible chess colors.

	NoColor Color = iota
	Black
	White
	_

	ColorArraySize = int(iota)
	ColorMinValue  = Black
	ColorMaxValue  = White
)

const (
	// Set of possible chess figures.

	NoFigure Figure = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King

	FigureArraySize = int(iota)
	FigureMinValue  = Pawn
	FigureMaxValue  = King
)

// piece constants must stay in sync with ColorFigure
// the order of pieces must match Polyglot format:
// http://hgm.nubati.net/book_format.html
const (
	NoPiece Piece = iota
	_
	BlackPawn
	WhitePawn
	BlackKnight
	WhiteKnight
	BlackBishop
	WhiteBishop
	BlackRook
	WhiteRook
	BlackQueen
	WhiteQueen
	BlackKing
	WhiteKing
	_
	_

	PieceArraySize = int(iota)
	PieceMinValue  = BlackPawn
	PieceMaxValue  = WhiteKing
)

// move types
const (
	NoMove    MoveType = iota // no move or null move
	Normal                    // regular move
	Promotion                 // pawn is promoted. Move.Promotion() gives the new piece
	Castling                  // king castles
	Enpassant                 // pawn takes enpassant
)

// null move is a move that does nothing, its value is 0
const (
	NullMove = Move(0)
)

// useful bitboards
const (
	BbEmpty          Bitboard = 0x0000000000000000
	BbFull           Bitboard = 0xffffffffffffffff
	BbBorder         Bitboard = 0xff818181818181ff
	BbPawnStartRank  Bitboard = 0x00ff00000000ff00
	BbPawnDoubleRank Bitboard = 0x000000ffff000000
	BbBlackSquares   Bitboard = 0xaa55aa55aa55aa55
	BbWhiteSquares   Bitboard = 0x55aa55aa55aa55aa
)

const (
	BbFileA Bitboard = 0x101010101010101 << iota
	BbFileB
	BbFileC
	BbFileD
	BbFileE
	BbFileF
	BbFileG
	BbFileH
)

const (
	BbRank1 Bitboard = 0x0000000000000FF << (8 * iota)
	BbRank2
	BbRank3
	BbRank4
	BbRank5
	BbRank6
	BbRank7
	BbRank8
)

const (
	// the set of all board squares

	SquareA1 = Square(iota)
	SquareB1
	SquareC1
	SquareD1
	SquareE1
	SquareF1
	SquareG1
	SquareH1
	SquareA2
	SquareB2
	SquareC2
	SquareD2
	SquareE2
	SquareF2
	SquareG2
	SquareH2
	SquareA3
	SquareB3
	SquareC3
	SquareD3
	SquareE3
	SquareF3
	SquareG3
	SquareH3
	SquareA4
	SquareB4
	SquareC4
	SquareD4
	SquareE4
	SquareF4
	SquareG4
	SquareH4
	SquareA5
	SquareB5
	SquareC5
	SquareD5
	SquareE5
	SquareF5
	SquareG5
	SquareH5
	SquareA6
	SquareB6
	SquareC6
	SquareD6
	SquareE6
	SquareF6
	SquareG6
	SquareH6
	SquareA7
	SquareB7
	SquareC7
	SquareD7
	SquareE7
	SquareF7
	SquareG7
	SquareH7
	SquareA8
	SquareB8
	SquareC8
	SquareD8
	SquareE8
	SquareF8
	SquareG8
	SquareH8

	SquareArraySize = int(iota)
	SquareMinValue  = SquareA1
	SquareMaxValue  = SquareH8
)

var (
	// maps figures to symbols for move notations
	lanFigureToSymbol = [...]string{"", "", "N", "B", "R", "Q", "K"}
	uciFigureToSymbol = [...]string{"", "", "n", "b", "r", "q", "k"}
)

const (
	// WhiteOO indicates that White can castle on King side.
	WhiteOO Castle = 1 << iota
	// WhiteOOO indicates that White can castle on Queen side.
	WhiteOOO
	// BlackOO indicates that Black can castle on King side.
	BlackOO
	// BlackOOO indicates that Black can castle on Queen side.
	BlackOOO

	// NoCastle indicates no castling rights.
	NoCastle Castle = 0
	// AnyCastle indicates all castling rights.
	AnyCastle Castle = WhiteOO | WhiteOOO | BlackOO | BlackOOO

	CastleArraySize = int(AnyCastle + 1)
	CastleMinValue  = NoCastle
	CastleMaxValue  = AnyCastle
)

var castleToString = [...]string{
	"-", "K", "Q", "KQ", "k", "Kk", "Qk", "KQk", "q", "Kq", "Qq", "KQq", "kq", "Kkq", "Qkq", "KQkq",
}

/////////////////////////////////////////////////////////////////////
