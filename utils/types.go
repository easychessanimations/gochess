package utils

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// types

type VariantKey uint8

const (
	VARIANT_STANDARD VariantKey = iota
	VARIANT_ATOMIC
	VARIANT_SEIRAWAN
	VARIANT_EIGHTPIECE
)

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

type Square struct {
	File int8
	Rank int8
}

type PieceKind uint8

type PieceColor uint8

type PieceDirection Square

type Piece struct {
	Kind      PieceKind
	Color     PieceColor
	Direction PieceDirection
}

type PieceDescriptor struct {
	Directions          []PieceDirection
	Sliding             bool
	CanJumpOverOwnPiece bool
	CanCapture          bool
}

/////////////////////////////////////////////////////////////////////
