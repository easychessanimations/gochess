package board

type Square struct {
	File int8
	Rank int8
}

type PieceKind uint8

type PieceColor bool

type PieceDirection Square

type Piece struct {
	Kind      PieceKind
	Color     PieceColor
	Direction PieceDirection
}

type BoardRep struct {
	NumFiles int8
	LastFile int8
	NumRanks int8
	LastRank int8
	Rep      map[Square]Piece
}

type VariantKey uint8

type Pos struct {
	Turn PieceColor
}

type MoveStackItem struct {
	Rep BoardRep
	Pos Pos
}

type Board struct {
	Variant   VariantKey
	Rep       BoardRep
	Pos       Pos
	MoveStack []MoveStackItem
}

type PieceDescriptor struct {
	Directions          []PieceDirection
	Sliding             bool
	CanJumpOverOwnPiece bool
	CanCapture          bool
}

type Move struct {
	FromSq  Square
	ToSq    Square
	Capture bool
}

type MoveBuffItem struct {
	Move Move
	Str  string
}

type MoveBuff []MoveBuffItem

func (mb MoveBuff) Len() int {
	return len(mb)
}
func (mb MoveBuff) Swap(i, j int) {
	mb[i], mb[j] = mb[j], mb[i]
}
func (mb MoveBuff) Less(i, j int) bool {
	return mb[i].Str < mb[j].Str
}
