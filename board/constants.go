package board

const (
	VARIANT_STANDARD VariantKey = iota
	VARIANT_ATOMIC
	VARIANT_SEIRAWAN
	VARIANT_EIGHTPIECE
)

var START_FENS = map[VariantKey]string{
	VARIANT_STANDARD:   "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	VARIANT_ATOMIC:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	VARIANT_SEIRAWAN:   "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR[EHeh] w KQBCDFGkqbcdfg - 0 1",
	VARIANT_EIGHTPIECE: "jlsesqkbnr/pppppppp/8/8/8/8/PPPPPPPP/JLneSQKBNR w KQkq - 0 1",
}

const WHITE PieceColor = true
const BLACK PieceColor = false

const (
	NO_PIECE PieceKind = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
	Hawk
	Elephant
	Sentry
	Jailer
	Lancer
)

var PIECE_KIND_TO_PIECE_LETTER map[PieceKind]string = map[PieceKind]string{
	NO_PIECE: "-",
	Pawn:     "p",
	Knight:   "n",
	Bishop:   "b",
	Rook:     "r",
	Queen:    "q",
	King:     "k",
	Hawk:     "h",
	Elephant: "e",
	Sentry:   "s",
	Jailer:   "j",
	Lancer:   "l",
}

var PIECE_LETTER_TO_PIECE_KIND map[string]PieceKind = map[string]PieceKind{
	"-": NO_PIECE,
	"p": Pawn,
	"n": Knight,
	"b": Bishop,
	"r": Rook,
	"q": Queen,
	"k": King,
	"h": Hawk,
	"e": Elephant,
	"s": Sentry,
	"j": Jailer,
	"l": Lancer,
}

var DIRECTION_STRING_TO_PIECE_DIRECTION = map[string]PieceDirection{
	"n":  PieceDirection{0, -1},
	"ne": PieceDirection{1, -1},
	"e":  PieceDirection{1, 0},
	"se": PieceDirection{1, 1},
	"s":  PieceDirection{0, 1},
	"sw": PieceDirection{-1, 1},
	"w":  PieceDirection{-1, 0},
	"nw": PieceDirection{-1, -1},
}

var PIECE_DIRECTION_TO_DIRECTION_STRING = map[PieceDirection]string{
	PieceDirection{0, -1}:  "n",
	PieceDirection{1, -1}:  "ne",
	PieceDirection{1, 0}:   "e",
	PieceDirection{1, 1}:   "se",
	PieceDirection{0, 1}:   "s",
	PieceDirection{-1, 1}:  "sw",
	PieceDirection{-1, 0}:  "w",
	PieceDirection{-1, -1}: "nw",
}
