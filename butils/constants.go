package butils

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// constants

const (
	// set of possible chess colors

	NoColor Color = iota
	Black
	White
	_

	ColorArraySize = int(iota)
	ColorMinValue  = Black
	ColorMaxValue  = White
)

const (
	// set of possible chess figures

	NoFigure Figure = iota // 0
	Pawn                   // 1
	Knight                 // 2
	Bishop                 // 3
	Rook                   // 4
	Queen                  // 5
	King                   // 6
	Jailer                 // 7
	LancerN                // 8
	LancerNE               // 9
	LancerE                // 10
	LancerSE               // 11
	LancerS                // 12
	LancerSW               // 13
	LancerW                // 14
	LancerNW               // 15
	Sentry                 // 16
	_                      // 17
	_                      // 18
	_                      // 19
	_                      // 20
	_                      // 21
	_                      // 22
	_                      // 23
	_                      // 24
	_                      // 25
	_                      // 26
	_                      // 27
	_                      // 28
	_                      // 29
	_                      // 30
	_                      // 31

	FigureArraySize = int(iota) // 32
	FigureMinValue  = Pawn      // 1
	FigureMaxValue  = King      // 6
)

// piece constants must stay in sync with ColorFigure
// the order of pieces must match Polyglot format:
// http://hgm.nubati.net/book_format.html
const (
	NoPiece       Piece = iota // 0
	_                          // 1
	BlackPawn                  // 2
	WhitePawn                  // 3
	BlackKnight                // 4
	WhiteKnight                // 5
	BlackBishop                // 6
	WhiteBishop                // 7
	BlackRook                  // 8
	WhiteRook                  // 9
	BlackQueen                 // 10
	WhiteQueen                 // 11
	BlackKing                  // 12
	WhiteKing                  // 13
	BlackJailer                // 14
	WhiteJailer                // 15
	BlackLancerN               // 16
	WhiteLancerN               // 17
	BlackLancerNE              // 18
	WhiteLancerNE              // 19
	BlackLancerE               // 20
	WhiteLancerE               // 21
	BlackLancerSE              // 22
	WhiteLancerSE              // 23
	BlackLancerS               // 24
	WhiteLancerS               // 25
	BlackLancerSW              // 26
	WhiteLancerSW              // 27
	BlackLancerW               // 28
	WhiteLancerW               // 29
	BlackLancerNW              // 30
	WhiteLancerNW              // 31
	BlackSentry                // 32
	WhiteSentry                // 33
	_                          // 34
	_                          // 35
	_                          // 36
	_                          // 37
	_                          // 38
	_                          // 39
	_                          // 40
	_                          // 41
	_                          // 42
	_                          // 43
	_                          // 44
	_                          // 45
	_                          // 46
	_                          // 47
	_                          // 48
	_                          // 49
	_                          // 50
	_                          // 51
	_                          // 52
	_                          // 53
	_                          // 54
	_                          // 55
	_                          // 56
	_                          // 57
	_                          // 58
	_                          // 59
	_                          // 60
	_                          // 61
	_                          // 62
	_                          // 63

	PieceArraySize = int(iota) // 64
	PieceMinValue  = BlackPawn // 2
	PieceMaxValue  = WhiteKing // 13
)

// PieceArraySize should be an exponent of 2
// PieceArraySize = 2 ^ PIECE_ARRAY_SIZE_IN_BITS
// TODO: PIECE_ARRAY_SIZE_IN_BITS should be used instead of
// hard coded constants everywhere in the code
const PIECE_ARRAY_SIZE_IN_BITS = 6

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
	// WhiteOO indicates that White can castle on King side
	WhiteOO Castle = 1 << iota
	// WhiteOOO indicates that White can castle on Queen side
	WhiteOOO
	// BlackOO indicates that Black can castle on King side
	BlackOO
	// BlackOOO indicates that Black can castle on Queen side
	BlackOOO

	// NoCastle indicates no castling rights
	NoCastle Castle = 0
	// AnyCastle indicates all castling rights
	AnyCastle Castle = WhiteOO | WhiteOOO | BlackOO | BlackOOO

	CastleArraySize = int(AnyCastle + 1)
	CastleMinValue  = NoCastle
	CastleMaxValue  = AnyCastle
)

var castleToString = [...]string{
	"-", "K", "Q", "KQ", "k", "Kk", "Qk", "KQk", "q", "Kq", "Qq", "KQq", "kq", "Kkq", "Qkq", "KQkq",
}

const (
	// Violent indicates captures (including en passant) and queen promotions
	Violent int = 1 << iota
	// Quiet are all other moves including minor promotions and castling
	Quiet
	// All moves
	All = Violent | Quiet
)

var (
	// FENStartPos is the FEN string of the starting position
	FENStartPos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	// which castle rights are lost when pieces are moved
	lostCastleRights = [64]Castle{
		SquareA1: WhiteOOO,
		SquareE1: WhiteOOO | WhiteOO,
		SquareH1: WhiteOO,
		SquareA8: BlackOOO,
		SquareE8: BlackOOO | BlackOO,
		SquareH8: BlackOO,
	}

	// the zobrist* arrays contain magic numbers used for Zobrist hashing
	// more information on Zobrist hashing can be found in the paper:
	// http://research.cs.wisc.edu/techreports/1970/TR88.pdf
	zobristPiece     [PieceArraySize][SquareArraySize]uint64
	zobristEnpassant [SquareArraySize]uint64
	zobristCastle    [CastleArraySize]uint64
	zobristColor     [ColorArraySize]uint64

	// maps runes to figures
	symbolToFigure = [256]Figure{
		'p': Pawn, 'n': Knight, 'b': Bishop, 'r': Rook, 'q': Queen, 'k': King,
		'P': Pawn, 'N': Knight, 'B': Bishop, 'R': Rook, 'Q': Queen, 'K': King,
	}
	// maps pieces to symbols
	//prettyPieceToSymbol = []string{".", "?", "♟", "♙", "♞", "♘", "♝", "♗", "♜", "♖", "♛", "♕", "♚", "♔"}
	prettyPieceToSymbol = []string{" . ", " ? ", " p ", " P ", " n ", " N ", " b ", " B ", " r ", " R ", " q ", " Q ", " k ", " K "}
)

// conversions
var (
	colorToSymbol      = "?bw"
	pieceToSymbol      = ".?pPnNbBrRqQkK"
	symbolToCastleInfo = map[rune]castleInfo{
		'K': castleInfo{
			Castle: WhiteOO,
			Piece:  [2]Piece{WhiteKing, WhiteRook},
			Square: [2]Square{SquareE1, SquareH1},
		},
		'k': castleInfo{
			Castle: BlackOO,
			Piece:  [2]Piece{BlackKing, BlackRook},
			Square: [2]Square{SquareE8, SquareH8},
		},
		'Q': castleInfo{
			Castle: WhiteOOO,
			Piece:  [2]Piece{WhiteKing, WhiteRook},
			Square: [2]Square{SquareE1, SquareA1},
		},
		'q': castleInfo{
			Castle: BlackOOO,
			Piece:  [2]Piece{BlackKing, BlackRook},
			Square: [2]Square{SquareE8, SquareA8},
		},
	}
	symbolToColor = map[string]Color{
		"w": White,
		"b": Black,
	}
	symbolToPiece = map[rune]Piece{
		'p': BlackPawn,
		'n': BlackKnight,
		'b': BlackBishop,
		'r': BlackRook,
		'q': BlackQueen,
		'k': BlackKing,

		'P': WhitePawn,
		'N': WhiteKnight,
		'B': WhiteBishop,
		'R': WhiteRook,
		'Q': WhiteQueen,
		'K': WhiteKing,
	}
)

var (
	// bbPawnAttack contains pawn's attack tables
	bbPawnAttack [64]Bitboard
	// bbKnightAttack contains knight's attack tables
	bbKnightAttack [64]Bitboard
	// bbKingAttack contains king's attack tables (excluding castling)
	bbKingAttack [64]Bitboard
	bbKingArea   [64]Bitboard
	// bbSuperAttack contains queen piece's attack tables. This queen can jump
	bbSuperAttack [64]Bitboard

	rookMagic    [64]magicInfo
	rookDeltas   = [][2]int{{-1, +0}, {+1, +0}, {+0, -1}, {+0, +1}}
	bishopMagic  [64]magicInfo
	bishopDeltas = [][2]int{{-1, +1}, {+1, +1}, {+1, -1}, {-1, -1}}
)

/////////////////////////////////////////////////////////////////////
