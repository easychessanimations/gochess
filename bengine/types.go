package bengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	"sync"
	"time"

	. "github.com/easychessanimations/gochess/butils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// types

/*
// feature type
type featureType string
*/

type featureType int

type FeatureInfo struct {
	Name  featureType // name of this feature
	Start int         // start position in the weights array
	Num   int         // number of weights used
}

// Score represents a pair of mid and end game scores
type Score struct {
	M, E int32 // mid game, end game
	//I    int   // index in Weights
}

// Accum accumulates scores
type Accum struct {
	M, E int32 // mid game, end game
	//Values []int8 // input values
}

// Options keeps engine's options
type Options struct {
	AnalyseMode   bool // true to display info strings
	MultiPV       int  // number of principal variation lines to compute
	HandicapLevel int
}

// Stats stores statistics about the search
type Stats struct {
	CacheHit  uint64 // number of times the position was found transposition table
	CacheMiss uint64 // number of times the position was not found in the transposition table
	Nodes     uint64 // number of nodes searched
	Depth     int32  // depth search
	SelDepth  int32  // maximum depth reached on PV (doesn't include the hash moves)
}

// Logger logs search progress
type Logger interface {
	// BeginSearch signals a new search is started
	BeginSearch()
	// EndSearch signals end of search
	EndSearch()
	// PrintPV logs the principal variation after iterative deepening completed one depth
	PrintPV(stats Stats, multiPV int, score int32, pv []Move)
	// CurrMove logs the current move. Current move index is 1-based
	CurrMove(depth int, move Move, num int)
}

// NulLogger is a logger that does nothing
type NulLogger struct{}

func (nl *NulLogger) BeginSearch()                                             {}
func (nl *NulLogger) EndSearch()                                               {}
func (nl *NulLogger) PrintPV(stats Stats, multiPV int, score int32, pv []Move) {}
func (nl *NulLogger) CurrMove(depth int, move Move, num int)                   {}

// historyEntry keeps counts of how well move performed in the past
type historyEntry struct {
	stat int32
	move Move
}

// historyTable is a hash table that contains history of moves
//
// old moves are automatically evicted when new moves are inserted
// so this cache is approx. LRU.
type historyTable [1 << 12]historyEntry

type orderedMove struct {
	move Move  // move
	key  int16 // sort key
}

// movesStack is a stack of moves
type moveStack struct {
	moves []orderedMove // list of moves with an order key
	buf   []Move        // a buffer of moves

	kind   int     // violent or all
	state  int     // current generation state
	hash   Move    // hash move
	killer [3]Move // two killer moves and one counter move
}

// stack is a stack of plies (movesStack)
type stack struct {
	position *Position
	moves    []moveStack
	history  *historyTable
	counter  [1 << 11]Move // counter moves table
}

type pvEntry struct {
	// lock is used to handled hash conflicts
	// normally set to position's Zobrist key
	lock uint64
	// move on pricipal variation for this position
	move Move
}

// pvTable is like hash table, but only to keep principal variation
//
// the additional table to store the PV was suggested by Robert Hyatt, see
//
// * http://www.talkchess.com/forum/viewtopic.php?topic_view=threads&p=369163&t=35982
// * http://www.talkchess.com/forum/viewtopic.php?t=36099
//
// during alpha-beta search entries that are on principal variation,
// are exact nodes, i.e. their score lies exactly between alpha and beta
type pvTable []pvEntry

// atomicFlag is an atomic bool that can only be set
type atomicFlag struct {
	lock sync.Mutex
	flag bool
}

func (af *atomicFlag) set() {
	af.lock.Lock()
	af.flag = true
	af.lock.Unlock()
}

func (af *atomicFlag) get() bool {
	af.lock.Lock()
	tmp := af.flag
	af.lock.Unlock()
	return tmp
}

// TimeControl is a time control that tries to split the
// remaining time over MovesToGo
type TimeControl struct {
	WTime, WInc time.Duration // time and increment for white
	BTime, BInc time.Duration // time and increment for black
	Depth       int32         // maximum depth search (including)
	MovesToGo   int32         // number of remaining moves, defaults to defaultMovesToGo

	sideToMove Color
	time, inc  time.Duration // time and increment for us
	limit      time.Duration

	predicted bool       // true if this move was predicted
	branch    int        // branching factor, multiplied by 16
	currDepth int32      // current depth searched
	stopped   atomicFlag // true to stop the search
	ponderhit atomicFlag // true if ponder was successful

	searchTime     time.Duration // alocated time for this move
	searchDeadline time.Time     // don't go to the next depth after this deadline
	stopDeadline   time.Time     // abort search after this deadline
}

// Engine implements the logic to search for the best move for a position
type Engine struct {
	Options  Options   // engine options
	Log      Logger    // logger
	Stats    Stats     // search statistics
	Position *Position // current Position

	rootPly         int           // position's ply at the start of the search
	stack           stack         // stack of moves
	pvTable         pvTable       // principal variation table
	pvTableAB       pvTableAB     // principal variation table for naive alphabeta
	UseAB           bool          // whether to use naive alphabeta
	history         *historyTable // keeps history of moves
	ignoreRootMoves []Move        // moves to ignore at root
	onlyRootMoves   []Move        // search only these root moves

	timeControl *TimeControl
	stopped     bool   // true if timeControl stopped the clock
	checkpoint  uint64 // when to check the time
}

type hashFlags uint8

// hashEntry is a value in the transposition table
type hashEntry struct {
	lock   uint32    // lock is used to handle hashing conflicts
	move   Move      // best move
	score  int16     // score of the position. if mate, score is relative to current position
	static int16     // static score of the position (not yet used)
	depth  int8      // remaining search depth
	kind   hashFlags // type of hash
}

// HashTable is a transposition table
// Engine uses this table to cache position scores so
// it doesn't have to research them again
type HashTable struct {
	table []hashEntry // len(table) is a power of two and equals mask+1
	mask  uint32      // mask is used to determine the index in the table
}

// Eval contains necessary information for evaluation
type Eval struct {
	// scores
	// - Accum[NoColor] is the combined score
	// - Accum[White] is White's score
	// - Accum[Black] is Black's score
	Accum [ColorArraySize]Accum
	// Position evaluated
	position *Position
}

// pawnsTable is a cache entry
type pawnsEntry struct {
	lock  uint64
	white Accum
	black Accum
}

// pawnsTable implements a fixed size cache
type pawnsTable [1 << 13]pawnsEntry

/////////////////////////////////////////////////////////////////////
