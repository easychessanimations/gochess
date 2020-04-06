package uciengine

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// constants

const ENGINE_NAME = "gochess"
const ENGINE_DESCRIPTION = "multi variant multi platform uci engine"
const ENGINE_AUTHOR = "easychessanimations"

const SEARCH_MAX_DEPTH = 100
const DEFAULT_QUIESCENCE_DEPTH = SEARCH_MAX_DEPTH

var UCI_OPTIONS = []UciOption{
	{
		Kind:       "spin",
		Name:       "Quiescence Depth",
		ValueKind:  "int",
		MinInt:     0,
		MaxInt:     SEARCH_MAX_DEPTH,
		DefaultInt: DEFAULT_QUIESCENCE_DEPTH,
	},
	{
		Kind:        "check",
		Name:        "Use AlphaBeta",
		ValueKind:   "bool",
		DefaultBool: true,
	},
}

/////////////////////////////////////////////////////////////////////
