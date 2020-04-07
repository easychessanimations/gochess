package board

/////////////////////////////////////////////////////////////////////
// imports

import (
	"time"

	"github.com/easychessanimations/gochess/utils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// types

const MAX_FILES = 8
const MAX_RANKS = 8

type BoardRep [MAX_RANKS][MAX_FILES]utils.Piece

type CastlingRight struct {
	Color          utils.PieceColor
	Side           CastlingSide
	CanCastle      bool
	RookOrigSquare utils.Square
	RookOrigPiece  utils.Piece
}

type ColorCastlingRights [2]CastlingRight

type CastlingRights [2]ColorCastlingRights

type Pos struct {
	Rep            BoardRep
	Turn           utils.PieceColor
	CastlingRights CastlingRights
	EpSquare       utils.Square
	HalfmoveClock  int
	FullmoveNumber int
}

type MoveStackItem struct {
	Pos  Pos
	Move Move
	San  string
}

type Board struct {
	Variant                           utils.VariantKey
	NumFiles                          int8
	LastFile                          int8
	NumRanks                          int8
	LastRank                          int8
	Pos                               Pos
	MoveStack                         []MoveStackItem
	Nodes                             int
	Start                             time.Time
	LogFunc                           func(string)
	LogAnalysisInfoFunc               func(string)
	SortedSanMoveBuff                 MoveBuff
	SelDepth                          int
	Alphas                            int
	Betas                             int
	Searching                         bool
	PositionHash                      PositionHash
	TestBoard                         *Board
	GetUciOptionByNameWithDefaultFunc func(string, utils.UciOption) utils.UciOption
	Multipvs                          []MultipvInfo
}

type Move struct {
	FromSq          utils.Square
	ToSq            utils.Square
	Capture         bool
	PawnCapture     bool
	PawnPushByOne   bool
	PawnPushByTwo   bool
	EpSquare        utils.Square
	EpCapture       bool
	EpClearSquare   utils.Square
	PromotionPiece  utils.Piece
	PromotionSquare utils.Square
	Castling        bool
	CastlingSide    CastlingSide
	RookOrigPiece   utils.Piece
}

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

type AlphaBetaInfo struct {
	Alpha           int
	Beta            int
	Depth           int
	QuiescenceDepth int
	CurrentDepth    int
	Line            []string
}

type MoveEntry struct {
	Eval int
}

type PositionEntry struct {
	MoveEntries map[Move]MoveEntry
}

type PositionHash struct {
	PositionEntries map[Pos]PositionEntry
}

type MultipvInfo struct {
	Depth    int
	SelDepth int
	Nodes    int
	Time     float32
	Nps      float32
	Alphas   int
	Betas    int
	Score    int
	Pv       string
	PvMoves  []Move
}

/////////////////////////////////////////////////////////////////////
