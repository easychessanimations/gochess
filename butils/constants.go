package butils

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// constants

const (
	// set of possible chess colors

	NoColor    Color = iota // 0
	Black                   // 1
	White                   // 2
	DummyColor              // 3

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
	DummyPiece                 // 1
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
	lanFigureToSymbol = [...]string{
		"",  // 0
		"",  // 1
		"N", // 2
		"B", // 3
		"R", // 4
		"Q", // 5
		"K", // 6
	}
	uciFigureToSymbol = [...]string{
		"",  // 0
		"",  // 1
		"n", // 2
		"b", // 3
		"r", // 4
		"q", // 5
		"k", // 6
	}
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
	/*symbolToFigure = [256]Figure{
		'p': Pawn, 'n': Knight, 'b': Bishop, 'r': Rook, 'q': Queen, 'k': King,
		'P': Pawn, 'N': Knight, 'B': Bishop, 'R': Rook, 'Q': Queen, 'K': King,
	}*/
	symbolToFigure = map[string]Figure{
		".":   NoFigure, // 0
		"?":   NoFigure, // 0
		"p":   Pawn,     // 1
		"P":   Pawn,     // 1
		"n":   Knight,   // 2
		"N":   Knight,   // 2
		"b":   Bishop,   // 3
		"B":   Bishop,   // 3
		"r":   Rook,     // 4
		"R":   Rook,     // 4
		"q":   Queen,    // 5
		"Q":   Queen,    // 5
		"k":   King,     // 6
		"K":   King,     // 6
		"j":   Jailer,   // 7
		"J":   Jailer,   // 7
		"ln":  LancerN,  // 8
		"Ln":  LancerN,  // 8
		"lne": LancerNE, // 9
		"Lne": LancerNE, // 9
		"le":  LancerE,  // 10
		"Le":  LancerE,  // 10
		"lse": LancerSE, // 11
		"Lse": LancerSE, // 11
		"ls":  LancerS,  // 12
		"Ls":  LancerS,  // 12
		"lsw": LancerSW, // 13
		"Lsw": LancerSW, // 13
		"lw":  LancerW,  // 14
		"Lw":  LancerW,  // 14
		"lnw": LancerNW, // 15
		"Lnw": LancerNW, // 15
		"s":   Sentry,   // 16
		"S":   Sentry,   // 16
	}
	// maps Piece to fen piece letter
	pieceToSymbol = map[Piece]string{
		NoPiece:       ".",   // 0
		DummyPiece:    "?",   // 1
		BlackPawn:     "p",   // 2
		WhitePawn:     "P",   // 3
		BlackKnight:   "n",   // 4
		WhiteKnight:   "N",   // 5
		BlackBishop:   "b",   // 6
		WhiteBishop:   "B",   // 7
		BlackRook:     "r",   // 8
		WhiteRook:     "R",   // 9
		BlackQueen:    "q",   // 10
		WhiteQueen:    "Q",   // 11
		BlackKing:     "k",   // 12
		WhiteKing:     "K",   // 13
		BlackJailer:   "j",   // 14
		WhiteJailer:   "J",   // 15
		BlackLancerN:  "ln",  // 16
		WhiteLancerN:  "Ln",  // 17
		BlackLancerNE: "lne", // 18
		WhiteLancerNE: "Lne", // 19
		BlackLancerE:  "le",  // 21
		WhiteLancerE:  "Le",  // 21
		BlackLancerSE: "lse", // 22
		WhiteLancerSE: "Lse", // 23
		BlackLancerS:  "ls",  // 24
		WhiteLancerS:  "Ls",  // 25
		BlackLancerSW: "lsw", // 26
		WhiteLancerSW: "Lsw", // 27
		BlackLancerW:  "lw",  // 28
		WhiteLancerW:  "Lw",  // 29
		BlackLancerNW: "lnw", // 31
		WhiteLancerNW: "Lnw", // 31
		BlackSentry:   "s",   // 32
		WhiteSentry:   "S",   // 33
	}
	// maps pieces to symbols
	//prettyPieceToSymbol = []string{".", "?", "♟", "♙", "♞", "♘", "♝", "♗", "♜", "♖", "♛", "♕", "♚", "♔"}
	prettyPieceToSymbol = []string{
		" . ", // 0
		" ? ", // 1
		" p ", // 2
		" P ", // 3
		" n ", // 4
		" N ", // 5
		" b ", // 6
		" B ", // 7
		" r ", // 8
		" R ", // 9
		" q ", // 10
		" Q ", // 11
		" k ", // 12
		" K ", // 13
		" j ", // 14
		" J ", // 15
		"ln ", // 16
		"Ln ", // 17
		"lne", // 18
		"Lne", // 19
		"le ", // 20
		"Le ", // 21
		"lse", // 22
		"Lse", // 23
		"ls ", // 24
		"Ls ", // 25
		"lsw", // 26
		"Lsw", // 27
		"lw ", // 28
		"Lw ", // 29
		"lnw", // 30
		"Lnw", // 31
		" s ", // 32
		" S ", // 33
	}
)

// conversions
var (
	colorToSymbol = "?bw"
	//pieceToSymbol      = ".?pPnNbBrRqQkK"
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
	/*symbolToPiece = map[rune]Piece{
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
	}*/
	symbolToPiece = map[string]Piece{
		".":   NoPiece,       // 0
		"?":   DummyPiece,    // 1
		"p":   BlackPawn,     // 2
		"P":   WhitePawn,     // 3
		"n":   BlackKnight,   // 4
		"N":   WhiteKnight,   // 5
		"b":   BlackBishop,   // 6
		"B":   WhiteBishop,   // 7
		"r":   BlackRook,     // 8
		"R":   WhiteRook,     // 9
		"q":   BlackQueen,    // 10
		"Q":   WhiteQueen,    // 11
		"k":   BlackKing,     // 12
		"K":   WhiteKing,     // 13
		"j":   BlackJailer,   // 14
		"J":   WhiteJailer,   // 15
		"ln":  BlackLancerN,  // 16
		"Ln":  WhiteLancerN,  // 17
		"lne": BlackLancerNE, // 18
		"Lne": WhiteLancerNE, // 19
		"le":  BlackLancerE,  // 20
		"Le":  WhiteLancerE,  // 21
		"lse": BlackLancerSE, // 22
		"Lse": WhiteLancerSE, // 23
		"ls":  BlackLancerS,  // 24
		"Ls":  WhiteLancerS,  // 25
		"lsw": BlackLancerSW, // 26
		"Lsw": WhiteLancerSW, // 27
		"lw":  BlackLancerW,  // 28
		"Lw":  WhiteLancerW,  // 29
		"lnw": BlackLancerNW, // 30
		"Lnw": WhiteLancerNW, // 31
		"s":   BlackSentry,   // 32
		"S":   WhiteSentry,   // 33
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
