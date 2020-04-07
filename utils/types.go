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

type CastlingSide uint8

type Move struct {
	FromSq          Square
	ToSq            Square
	Capture         bool
	PawnCapture     bool
	PawnPushByOne   bool
	PawnPushByTwo   bool
	EpSquare        Square
	EpCapture       bool
	EpClearSquare   Square
	PromotionPiece  Piece
	PromotionSquare Square
	Castling        bool
	CastlingSide    CastlingSide
	RookOrigPiece   Piece
}

type MoveList []Move

type MoveEvalBuffItem struct {
	Move Move
	Eval int
}

type MoveEvalBuff []MoveEvalBuffItem

func (meb MoveEvalBuff) Len() int {
	return len(meb)
}
func (meb MoveEvalBuff) Swap(i, j int) {
	meb[i], meb[j] = meb[j], meb[i]
}
func (meb MoveEvalBuff) Less(i, j int) bool {
	return meb[i].Eval > meb[j].Eval
}

type MoveBuffItem struct {
	Move  Move
	San   string
	Algeb string
}

type MoveBuff []MoveBuffItem

func (mb MoveBuff) Len() int {
	return len(mb)
}
func (mb MoveBuff) Swap(i, j int) {
	mb[i], mb[j] = mb[j], mb[i]
}
func (mb MoveBuff) Less(i, j int) bool {
	return mb[i].San < mb[j].San
}

/////////////////////////////////////////////////////////////////////
