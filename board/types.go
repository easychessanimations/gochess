package board

import "time"

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

const MAX_FILES = 8
const MAX_RANKS = 8

type BoardRep [MAX_RANKS][MAX_FILES]Piece

type VariantKey uint8

type CastlingRight struct {
	Color          PieceColor
	Side           CastlingSide
	CanCastle      bool
	RookOrigSquare Square
	RookOrigPiece  Piece
}

type ColorCastlingRights [2]CastlingRight

type CastlingRights [2]ColorCastlingRights

type Pos struct {
	Turn           PieceColor
	CastlingRights CastlingRights
	EpSquare       Square
	HalfmoveClock  int
	FullmoveNumber int
}

type MoveStackItem struct {
	Rep  BoardRep
	Pos  Pos
	Move Move
	San  string
}

type Board struct {
	Variant           VariantKey
	NumFiles          int8
	LastFile          int8
	NumRanks          int8
	LastRank          int8
	Rep               BoardRep
	Pos               Pos
	MoveStack         []MoveStackItem
	Nodes             int
	Start             time.Time
	LogFunc           func(string)
	SortedSanMoveBuff MoveBuff
}

type PieceDescriptor struct {
	Directions          []PieceDirection
	Sliding             bool
	CanJumpOverOwnPiece bool
	CanCapture          bool
}

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

type AlphaBetaInfo struct {
	Alpha           int
	Beta            int
	Depth           int
	QuiescenceDepth int
	CurrentDepth    int
	Line            []Move
}
