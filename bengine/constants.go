package bengine

import "time"

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// constants

// engine
const (
	checkDepthExtension int32 = 1 // how much to extend search in case of checks
	lmrDepthLimit       int32 = 4 // do not do LMR below and including this limit
	futilityDepthLimit  int32 = 4 // maximum depth to do futility pruning

	initialAspirationWindow = 13
	futilityMargin          = 75
	checkpointStep          = 10000
)

// hash table
const (
	exact      hashFlags = 1 << iota // exact score is known
	failedLow                        // Search failed low, upper bound.
	failedHigh                       // Search failed high, lower bound
	hasStatic                        // entry contains static evaluation
)

// material
const (
	// KnownWinScore is strictly greater than all evaluation scores (mate not included)
	KnownWinScore = 25000
	// KnownLossScore is strictly smaller than all evaluation scores (mated not included)
	KnownLossScore = -KnownWinScore
	// MateScore - N is mate in N plies
	MateScore = 30000
	// MatedScore + N is mated in N plies
	MatedScore = -MateScore
	// InfinityScore is possible score, -InfinityScore is the minimum possible score
	InfinityScore = 32000
)

// Move generation states
const (
	msHash          = iota // return hash move
	msGenViolent           // generate violent moves
	msReturnViolent        // return violent moves in order
	msGenKiller            // generate killer moves
	msReturnKiller         // return killer moves  in order
	msGenRest              // generate remaining moves
	msReturnRest           // return remaining moves in order
	msDone                 // all moves returned
)

// pv
const (
	pvTableSize = 1 << 14
	pvTableMask = pvTableSize - 1
)

// time control
const (
	defaultMovesToGo = 35 // default number of more moves expected to play
	infinite         = 1000000000 * time.Second
	overhead         = 20 * time.Millisecond
)

// main
const (
	maxMultiPV       = 16
	maxHandicapLevel = 20
)

/////////////////////////////////////////////////////////////////////
