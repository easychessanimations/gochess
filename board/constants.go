package board

const STOP_AT_FIRST = true
const ALL_ATTACKS = false
const ADD_SAN = true

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

const WHITE PieceColor = 1
const BLACK PieceColor = 0

type CastlingSide uint8

const KING_SIDE CastlingSide = 1
const QUEEN_SIDE CastlingSide = 0

const (
	NO_PIECE_KIND PieceKind = iota
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

var NO_PIECE = Piece{Kind: NO_PIECE_KIND}
var NO_SQUARE = Square{-1, -1}

var PIECE_KIND_TO_PIECE_LETTER = map[PieceKind]string{
	NO_PIECE_KIND: "-",
	Pawn:          "p",
	Knight:        "n",
	Bishop:        "b",
	Rook:          "r",
	Queen:         "q",
	King:          "k",
	Hawk:          "h",
	Elephant:      "e",
	Sentry:        "s",
	Jailer:        "j",
	Lancer:        "l",
}

var PIECE_LETTER_TO_PIECE_KIND = map[string]PieceKind{
	"-": NO_PIECE_KIND,
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

var PIECE_KIND_TO_PIECE_DESCRIPTOR = map[PieceKind]PieceDescriptor{
	Knight: PieceDescriptor{
		Directions: []PieceDirection{
			PieceDirection{1, 2},
			PieceDirection{-1, 2},
			PieceDirection{1, -2},
			PieceDirection{-1, -2},
			PieceDirection{2, 1},
			PieceDirection{-2, 1},
			PieceDirection{2, -1},
			PieceDirection{-2, -1},
		},
		Sliding:             false,
		CanJumpOverOwnPiece: true,
		CanCapture:          true,
	},
	Bishop: PieceDescriptor{
		Directions: []PieceDirection{
			PieceDirection{1, 1},
			PieceDirection{1, -1},
			PieceDirection{-1, 1},
			PieceDirection{-1, -1},
		},
		Sliding:             true,
		CanJumpOverOwnPiece: false,
		CanCapture:          true,
	},
	Rook: PieceDescriptor{
		Directions: []PieceDirection{
			PieceDirection{1, 0},
			PieceDirection{-1, 0},
			PieceDirection{0, 1},
			PieceDirection{0, -1},
		},
		Sliding:             true,
		CanJumpOverOwnPiece: false,
		CanCapture:          true,
	},
	Queen: PieceDescriptor{
		Directions: []PieceDirection{
			PieceDirection{1, 0},
			PieceDirection{-1, 0},
			PieceDirection{0, 1},
			PieceDirection{0, -1},
			PieceDirection{1, 1},
			PieceDirection{1, -1},
			PieceDirection{-1, 1},
			PieceDirection{-1, -1},
		},
		Sliding:             true,
		CanJumpOverOwnPiece: false,
		CanCapture:          true,
	},
	King: PieceDescriptor{
		Directions: []PieceDirection{
			PieceDirection{1, 0},
			PieceDirection{-1, 0},
			PieceDirection{0, 1},
			PieceDirection{0, -1},
			PieceDirection{1, 1},
			PieceDirection{1, -1},
			PieceDirection{-1, 1},
			PieceDirection{-1, -1},
		},
		Sliding:             false,
		CanJumpOverOwnPiece: false,
		CanCapture:          true,
	},
	Sentry: PieceDescriptor{
		Directions: []PieceDirection{
			PieceDirection{1, 1},
			PieceDirection{1, -1},
			PieceDirection{-1, 1},
			PieceDirection{-1, -1},
		},
		Sliding:             true,
		CanJumpOverOwnPiece: false,
		CanCapture:          true,
	},
	Jailer: PieceDescriptor{
		Directions: []PieceDirection{
			PieceDirection{1, 0},
			PieceDirection{-1, 0},
			PieceDirection{0, 1},
			PieceDirection{0, -1},
		},
		Sliding:             true,
		CanJumpOverOwnPiece: false,
		CanCapture:          false,
	},
	Lancer: PieceDescriptor{
		Directions: []PieceDirection{
			PieceDirection{1, 0},
			PieceDirection{-1, 0},
			PieceDirection{0, 1},
			PieceDirection{0, -1},
			PieceDirection{1, 1},
			PieceDirection{1, -1},
			PieceDirection{-1, 1},
			PieceDirection{-1, -1},
		},
		Sliding:             true,
		CanJumpOverOwnPiece: false,
		CanCapture:          true,
	},
}

var PROMOTION_PIECES = map[VariantKey][]Piece{
	VARIANT_STANDARD: []Piece{
		Piece{Kind: Queen},
		Piece{Kind: Rook},
		Piece{Kind: Bishop},
		Piece{Kind: Knight},
	},
	VARIANT_ATOMIC: []Piece{
		Piece{Kind: Queen},
		Piece{Kind: Rook},
		Piece{Kind: Bishop},
		Piece{Kind: Knight},
	},
	VARIANT_SEIRAWAN: []Piece{
		Piece{Kind: Queen},
		Piece{Kind: Rook},
		Piece{Kind: Bishop},
		Piece{Kind: Knight},
		Piece{Kind: Elephant},
		Piece{Kind: Hawk},
	},
	VARIANT_EIGHTPIECE: []Piece{
		Piece{Kind: Queen},
		Piece{Kind: Rook},
		Piece{Kind: Bishop},
		Piece{Kind: Knight},
		Piece{Kind: Jailer},
		Piece{Kind: Sentry},
		Piece{Kind: Lancer, Direction: PieceDirection{1, 0}},
		Piece{Kind: Lancer, Direction: PieceDirection{1, 1}},
		Piece{Kind: Lancer, Direction: PieceDirection{0, 1}},
		Piece{Kind: Lancer, Direction: PieceDirection{0, -1}},
		Piece{Kind: Lancer, Direction: PieceDirection{-1, 0}},
		Piece{Kind: Lancer, Direction: PieceDirection{-1, -1}},
		Piece{Kind: Lancer, Direction: PieceDirection{-1, 0}},
		Piece{Kind: Lancer, Direction: PieceDirection{-1, 1}},
	},
}

const MAX_PIECE_KINDS = 20

var PIECE_VALUES [MAX_PIECE_KINDS]int

func init() {
	PIECE_VALUES[Pawn] = 100
	PIECE_VALUES[Knight] = 300
	PIECE_VALUES[Bishop] = 300
	PIECE_VALUES[Rook] = 500
	PIECE_VALUES[Queen] = 900
	PIECE_VALUES[King] = 0
	PIECE_VALUES[Hawk] = 600
	PIECE_VALUES[Elephant] = 800
	PIECE_VALUES[Sentry] = 300
	PIECE_VALUES[Jailer] = 400
	PIECE_VALUES[Lancer] = 700
}

const INFINITE_SCORE = 100000

const MATE_SCORE = 10000
const DRAW_SCORE = 0

const KNIGHT_ON_EDGE_DEDUCTION = 10
const KNIGHT_CLOSE_TO_EDGE_DEDUCTION = 10
const CENTER_PAWN_BONUS = 30
const MOBILITY_BONUS = 2
const RANDOM_BONUS = 10
const CAPTURE_BONUS = 2000
const NON_PAWN_MOVE_BONUS = 1000
