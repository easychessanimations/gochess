package board

/////////////////////////////////////////////////////////////////////
// imports

import (
	"github.com/easychessanimations/gochess/utils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// constants

var NO_MOVE = Move{FromSq: utils.NO_SQUARE}

const STOP_AT_FIRST = true
const ALL_ATTACKS = false
const ADD_SAN = true

type CastlingSide uint8

const KING_SIDE CastlingSide = 1
const QUEEN_SIDE CastlingSide = 0

var PIECE_VALUES [utils.MAX_PIECE_KINDS]int

func init() {
	PIECE_VALUES[utils.Pawn] = 100
	PIECE_VALUES[utils.Knight] = 300
	PIECE_VALUES[utils.Bishop] = 300
	PIECE_VALUES[utils.Rook] = 500
	PIECE_VALUES[utils.Queen] = 900
	PIECE_VALUES[utils.King] = 0
	PIECE_VALUES[utils.Hawk] = 600
	PIECE_VALUES[utils.Elephant] = 800
	PIECE_VALUES[utils.Sentry] = 300
	PIECE_VALUES[utils.Jailer] = 400
	PIECE_VALUES[utils.Lancer] = 700
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

const SEARCH_MAX_DEPTH = 100
const DEFAULT_QUIESCENCE_DEPTH = 0
const DEFAULT_UCI_VARIANT_STRING = "standard"
const DEFAULT_SEARCH_DEPTH = 10
const MAX_MULTIPV = 500
const DEFAULT_MULTIPV = 1

var UCI_OPTIONS = []utils.UciOption{
	{
		Kind:      "combo",
		Name:      "UCI_Variant",
		ValueKind: "string",
		Default:   DEFAULT_UCI_VARIANT_STRING,
		Value:     DEFAULT_UCI_VARIANT_STRING,
		DefaultStringArray: []string{
			"standard",
			"atomic",
			"seirawan",
			"eightpiece",
		},
		DefaultBool: true,
	},
	{
		Kind:       "spin",
		Name:       "MultiPV",
		ValueKind:  "int",
		MinInt:     0,
		MaxInt:     MAX_MULTIPV,
		DefaultInt: DEFAULT_MULTIPV,
		ValueInt:   DEFAULT_MULTIPV,
	},
	{
		Kind:       "spin",
		Name:       "Quiescence Depth",
		ValueKind:  "int",
		MinInt:     0,
		MaxInt:     0,
		DefaultInt: DEFAULT_QUIESCENCE_DEPTH,
		ValueInt:   DEFAULT_QUIESCENCE_DEPTH,
	},
}

var UCI_COMMAND_ALIASES = map[string]string{
	"v":  "setoption name UCI_Variant value standard\ni",
	"va": "setoption name UCI_Variant value atomic\ni",
	"vs": "setoption name UCI_Variant value seirawan\ni",
	"ve": "setoption name UCI_Variant value eightpiece\ni",
	"q":  "setoption name Quiescence Depth value 100",
}

/////////////////////////////////////////////////////////////////////
