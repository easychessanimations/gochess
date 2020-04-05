package uciengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	"github.com/easychessanimations/gochess/board"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// types

type UciEngine struct {
	Name        string
	Description string
	Author      string
	Board       board.Board
	Interactive bool
}

type UciOption struct {
	Kind               string
	Name               string
	ValueKind          string
	Default            string
	DefaultInt         int
	DefaultBool        bool
	DefaultStringArray []string
	Value              string
	ValueInt           int
	ValueBool          bool
	MinInt             int
	MaxInt             int
}

/////////////////////////////////////////////////////////////////////
