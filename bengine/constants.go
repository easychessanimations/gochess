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

/*const (
	// value of each figure
	fNoFigure featureType = "NoFigure"
	fPawn     featureType = "Pawn"
	fKnight   featureType = "Knight"
	fBishop   featureType = "Bishop"
	fRook     featureType = "Rook"
	fQueen    featureType = "Queen"
	fKing     featureType = "King"

	// PSqT for each figure from white's POV
	// for pawns evaluate each square, but other figures
	// we only evaluate the coordinates
	fPawnSquare featureType = "PawnSquare"
	fKnightFile featureType = "KnightFile"
	fKnightRank featureType = "KnightRank"
	fBishopFile featureType = "BishopFile"
	fBishopRank featureType = "BishopRank"
	fRookFile   featureType = "RookFile"
	fRookRank   featureType = "RookRank"
	fQueenFile  featureType = "QueenFile"
	fQueenRank  featureType = "QueenRank"
	fKingFile   featureType = "KingFile"
	fKingRank   featureType = "KingRank"

	// mobility of each figure
	fKnightAttack featureType = "KnightAttack"
	fBishopAttack featureType = "BishopAttack"
	fRookAttack   featureType = "RookAttack"
	fQueenAttack  featureType = "QueenAttack"
	fKingAttack   featureType = "KingAttack"

	// pawn structure
	fBackwardPawns  featureType = "BackwardPawns"
	fConnectedPawns featureType = "ConnectedPawns"
	fDoubledPawns   featureType = "DoubledPawns"
	fIsolatedPawns  featureType = "IsolatedPawns"
	fRammedPawns    featureType = "RammedPawns"
	fPassedPawnRank featureType = "PassedPawnRank"
	fPawnMobility   featureType = "PawnMobility"
	// minors and majors attacked by pawns
	fMinorsPawnsAttack featureType = "MinorsPawnsAttack"
	fMajorsPawnsAttack featureType = "MajorsPawnsAttack"
	// minors and majors attacked after a pawn push
	fMinorsPawnsPotentialAttack featureType = "MinorsPawnsPotentialAttack"
	fMajorsPawnsPotentialAttack featureType = "MajorsPawnsPotentialAttack"
	// how close is the king from a friendly passed pawn
	fKingPassedPawnTropism featureType = "KingPassedPawnTropism"
	// how close is the king from an enemy passed pawn
	fKingEnemyPassedPawnTropism featureType = "KingEnemyPassedPawnTropism"

	// attacked minors
	fAttackedMinors featureType = "AttackedMinors"
	// bishop pair
	fBishopPair featureType = "BishopPair"
	// rook on open and semi-open files
	fRookOnOpenFile     featureType = "RookOnOpenFile"
	fRookOnSemiOpenFile featureType = "RookOnSemiOpenFile"
	fKingQueenTropism   featureType = "KingQueenTropism"

	// king safety
	fKingAttackers featureType = "KingAttackers"
	// pawn in king's area
	fKingShelterNear featureType = "KingShelterNear"
	// pawn in front of the king, on the same file
	fKingShelterFront featureType = "KingShelterFront"
	// pawn in front of the king, including adjacent files
	fKingShelterFar featureType = "KingShelterFar"
)*/

const (
	fNoFigure                   featureType = 0
	fPawn                       featureType = 1
	fKnight                     featureType = 2
	fBishop                     featureType = 3
	fRook                       featureType = 4
	fQueen                      featureType = 5
	fKing                       featureType = 6
	fPawnMobility               featureType = 7
	fMinorsPawnsAttack          featureType = 8
	fMajorsPawnsAttack          featureType = 9
	fMinorsPawnsPotentialAttack featureType = 10
	fMajorsPawnsPotentialAttack featureType = 11
	fKnightFile                 featureType = 12
	fKnightRank                 featureType = 20
	fKnightAttack               featureType = 28
	fBishopFile                 featureType = 29
	fBishopRank                 featureType = 37
	fBishopAttack               featureType = 45
	fRookFile                   featureType = 46
	fRookRank                   featureType = 54
	fRookAttack                 featureType = 62
	fRookOnOpenFile             featureType = 63
	fRookOnSemiOpenFile         featureType = 64
	fQueenFile                  featureType = 65
	fQueenRank                  featureType = 73
	fQueenAttack                featureType = 81
	fKingQueenTropism           featureType = 82
	fAttackedMinors             featureType = 83
	fBishopPair                 featureType = 84
	fKingAttackers              featureType = 85
	fPawnSquare                 featureType = 89
	fBackwardPawns              featureType = 153
	fConnectedPawns             featureType = 154
	fDoubledPawns               featureType = 155
	fIsolatedPawns              featureType = 156
	fRammedPawns                featureType = 157
	fKingFile                   featureType = 158
	fKingRank                   featureType = 166
	fKingAttack                 featureType = 174
	fKingShelterNear            featureType = 175
	fKingShelterFar             featureType = 176
	fKingShelterFront           featureType = 177
	fPassedPawnRank             featureType = 178
	fKingEnemyPassedPawnTropism featureType = 186
	fKingPassedPawnTropism      featureType = 194
)

/////////////////////////////////////////////////////////////////////
