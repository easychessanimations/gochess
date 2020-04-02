package board

type Square struct {
	file int8
	rank int8
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
	NumRanks int8
	Rep      map[Square]Piece
}

type VariantKey uint8

type Board struct {
	Variant VariantKey
	Rep     BoardRep
}
