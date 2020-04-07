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
	Side           utils.CastlingSide
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
	DisabledMove   utils.Move
}

type MoveStackItem struct {
	Pos  Pos
	Move utils.Move
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
	SortedSanMoveBuff                 utils.MoveBuff
	SelDepth                          int
	Alphas                            int
	Betas                             int
	Searching                         bool
	PositionHash                      PositionHash
	TestBoard                         *Board
	GetUciOptionByNameWithDefaultFunc func(string, utils.UciOption) utils.UciOption
	MultipvInfos                      MultipvInfos
	ExcludedMoves                     []utils.Move
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
	MoveEntries map[utils.Move]MoveEntry
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
	PvMoves  []utils.Move
}

type MultipvInfos []MultipvInfo

func (mpi MultipvInfos) Len() int {
	return len(mpi)
}
func (mpi MultipvInfos) Swap(i, j int) {
	mpi[i], mpi[j] = mpi[j], mpi[i]
}
func (mpi MultipvInfos) Less(i, j int) bool {
	if mpi[i].Depth > mpi[j].Depth {
		return true
	}
	return mpi[i].Score > mpi[j].Score
}

/////////////////////////////////////////////////////////////////////
