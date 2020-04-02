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

type Board struct {
	Variant VariantKey
	Rep     BoardRep
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
