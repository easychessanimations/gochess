package uciengine

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// types

type UciEngine struct {
	Name        string
	Description string
	Author      string
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
