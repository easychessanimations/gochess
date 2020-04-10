package bengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	"sync"
	"time"
	"unsafe"

	. "github.com/easychessanimations/gochess/butils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// global variables

// mvvlva values based on one pawn = 10
var mvvlvaBonus = [...]int16{0, 10, 40, 45, 68, 145, 256}

// piece bonuses when calulating the see
// the values are fixed to approximatively the figure bonus in mid game
var seeBonus = [FigureArraySize]int32{0, 100, 357, 377, 712, 12534, 20000}

// gaps from Best Increments for the Average Case of Shellsort, Marcin Ciura
var shellSortGaps = [...]int{132, 57, 23, 10, 4, 1}

// distance stores the number of king steps required
// to reach from one square to another on an empty board
var distance [SquareArraySize][SquareArraySize]int32

var (
	// Placeholder for the weights array when running in coach mode
	Weights []Score

	FeaturesMap     = make(map[featureType]*FeatureInfo)
	featuresMapLock sync.Mutex
)

// engine
var (
	initialized = false
)

// hash table
var (
	// DefaultHashTableSizeMB is the default size in MB
	DefaultHashTableSizeMB = 64
	// GlobalHashTable is the global transposition table
	GlobalHashTable *HashTable
)

// material
var (
	// Evaluation caches
	pawnsAndShelterCache pawnsTable

	// Figure bonuses to use when computing the futility margin
	futilityFigureBonus [FigureArraySize]int32
)

// murmuir seed
var murmurSeed = [ColorArraySize]uint64{
	0x77a166129ab66e91,
	0x4f4863d5038ea3a3,
	0xe14ec7e648a4068b,
}

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// global functions

// NewTimeControl returns a new time control with no time limit
// no depth limit, zero time increment and zero moves to go
func NewTimeControl(pos *Position, predicted bool) *TimeControl {
	return &TimeControl{
		WTime:      infinite,
		WInc:       0,
		BTime:      infinite,
		BInc:       0,
		Depth:      64,
		MovesToGo:  defaultMovesToGo,
		sideToMove: pos.Us(),
		predicted:  predicted,
		branch:     32,
	}
}

// NewFixedDepthTimeControl returns a TimeControl which limits the search depth
func NewFixedDepthTimeControl(pos *Position, depth int32) *TimeControl {
	tc := NewTimeControl(pos, false)
	tc.Depth = depth
	tc.MovesToGo = 1
	return tc
}

// NewDeadlineTimeControl returns a TimeControl corresponding to a single move before deadline
func NewDeadlineTimeControl(pos *Position, deadline time.Duration) *TimeControl {
	tc := NewTimeControl(pos, false)
	tc.WTime = deadline
	tc.BTime = deadline
	tc.MovesToGo = 1
	return tc
}

// mvvlva computes Most Valuable Victim / Least Valuable Aggressor
// https://chessprogramming.wikispaces.com/MVV-LVA
func mvvlva(m Move) int16 {
	a := m.Target().Figure()
	v := m.Capture().Figure()
	return mvvlvaBonus[v]*64 - mvvlvaBonus[a]
}

// isInBounds returns true if score matches range defined by α, β and flags
func isInBounds(flags hashFlags, α, β, score int32) bool {
	if flags&exact != 0 {
		// simply return if the score is exact
		return true
	}
	if flags&failedLow != 0 && score <= α {
		// previously the move failed low so the actual score is at most
		// entry.score; if that is lower than α this will also fail low
		return true
	}
	if flags&failedHigh != 0 && score >= β {
		// previously the move failed high so the actual score is at least
		// entry.score; if that's higher than β this will also fail high
		return true
	}
	return false
}

// getBound returns the bound for score relative to α and β
func getBound(α, β, score int32) hashFlags {
	if score <= α {
		return failedLow
	}
	if score >= β {
		return failedHigh
	}
	return exact
}

// max returns maximum of a and b
func max(a, b int32) int32 {
	if a >= b {
		return a
	}
	return b
}

// min returns minimum of a and b
func min(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}

func getFeatureStart(feature featureType, num int) int {
	featuresMapLock.Lock()
	defer featuresMapLock.Unlock()

	if info, has := FeaturesMap[feature]; has {
		return info.Start
	}
	FeaturesMap[feature] = &FeatureInfo{
		Name:  feature,
		Start: len(Weights),
		Num:   num,
	}
	for i := 0; i < num; i++ {
		Weights = append(Weights, Score{M: 0, E: 0, I: len(Weights)})
	}
	return FeaturesMap[feature].Start
}

// group
func groupByCount(feature featureType, n int32, accum *Accum) {
	start := getFeatureStart(feature, 1)
	accum.addN(Weights[start], n)
}

func groupByBucket(feature featureType, n int, limit int, accum *Accum) {
	if n >= limit {
		n = limit - 1
	}
	start := getFeatureStart(feature, limit)
	accum.add(Weights[start+n])
}

func groupByBoard(feature featureType, bb Bitboard, accum *Accum) {
	groupByCount(feature, bb.Count(), accum)
}

func groupBySquare(feature featureType, us Color, bb Bitboard, accum *Accum) {
	start := getFeatureStart(feature, 64)
	for bb != BbEmpty {
		sq := bb.Pop().POV(us)
		accum.add(Weights[start+int(sq)])
	}
}

func groupByBool(feature featureType, b bool, accum *Accum) {
	start := getFeatureStart(feature, 1)
	if b {
		accum.add(Weights[start])
	}
}

func groupByFileSq(feature featureType, us Color, sq Square, accum *Accum) {
	groupByBucket(feature, sq.POV(us).File(), 8, accum)
}

func groupByRankSq(feature featureType, us Color, sq Square, accum *Accum) {
	groupByBucket(feature, sq.POV(us).Rank(), 8, accum)
}

func groupByRank(feature featureType, us Color, bb Bitboard, accum *Accum) {
	for bb != BbEmpty {
		sq := bb.Pop()
		groupByRankSq(feature, us, sq, accum)
	}
}

// murmuxMix function mixes two integers k&h
//
// murmurMix is based on MurmurHash2 https://sites.google.com/site/murmurhash/ which is on public domain
func murmurMix(k, h uint64) uint64 {
	h ^= k
	h *= uint64(0xc6a4a7935bd1e995)
	return h ^ (h >> uint(51))
}

// pawnsHash returns a hash of the pawns and king in position
func pawnsHash(pos *Position) uint64 {
	h := murmurSeed[pos.Us()]
	h = murmurMix(h, uint64(pos.ByPiece2(White, Pawn, King)))
	h = murmurMix(h, uint64(pos.ByPiece2(Black, Pawn, King)))
	h = murmurMix(h, uint64(pos.ByFigure(Pawn)))
	return h
}

// stub
func prefetch(e *hashEntry) {}

// split splits lock into a lock and two hash table indexes
// expects mask to be at least 3 bits
func split(lock uint64, mask uint32) (uint32, uint32, uint32) {
	hi := uint32(lock >> 32)
	lo := uint32(lock)
	h0 := lo & mask
	h1 := h0 ^ 1
	return hi, h0, h1
}

// NewHashTable builds transposition table that takes up to hashSizeMB megabytes
func NewHashTable(hashSizeMB int) *HashTable {
	// Choose hashSize such that it is a power of two.
	hashEntrySize := uint64(unsafe.Sizeof(hashEntry{}))
	hashSize := uint64(hashSizeMB) << 20 / hashEntrySize

	for hashSize&(hashSize-1) != 0 {
		hashSize &= hashSize - 1
	}
	return &HashTable{
		table: make([]hashEntry, hashSize),
		mask:  uint32(hashSize - 1),
	}
}

// newPvTable returns a new pvTable
func newPvTable() pvTable {
	return pvTable(make([]pvEntry, pvTableSize))
}

// NewEngine creates a new engine to search for pos
// if pos is nil then the start position is used
func NewEngine(pos *Position, log Logger, options Options) *Engine {
	if options.MultiPV == 0 {
		options.MultiPV = 1
	}

	if log == nil {
		log = &NulLogger{}
	}
	history := &historyTable{}
	eng := &Engine{
		Options: options,
		Log:     log,
		pvTable: newPvTable(),
		history: history,
		stack:   stack{history: history},
	}
	eng.SetPosition(pos)
	return eng
}

// CacheHitRatio returns the ratio of transposition table hits over total number of lookups
func (s *Stats) CacheHitRatio() float32 {
	return float32(s.CacheHit) / float32(s.CacheHit+s.CacheMiss)
}

// historyHash hashes the move and returns an index into the history table
func historyHash(m Move) uint32 {
	// this is a murmur inspired hash so upper bits are better
	// mixed than the lower bits; the hash multiplier was chosen
	// to minimize the number of misses
	h := uint32(m) * 438650727
	return (h + (h << 17)) >> 20
}

// isFutile return true if m cannot raise the current static
// evaluation above α; this is just an heuristic and mistakes
// can happen
func isFutile(pos *Position, static, α, margin int32, m Move) bool {
	if m.MoveType() == Promotion || m.Piece().Figure() == Pawn && BbPawnStartRank.Has(m.To()) {
		// promotion and passed pawns can increase the static evaluation
		// by more than futilityMargin
		return false
	}
	δ := futilityFigureBonus[m.Capture().Figure()]
	return static+δ+margin < α
}

func seeScore(m Move) int32 {
	score := seeBonus[m.Capture().Figure()]
	if m.MoveType() == Promotion {
		score -= seeBonus[Pawn]
		score += seeBonus[m.Target().Figure()]
	}
	return score
}

// see returns the static exchange evaluation for m, where is
// valid for current position (not yet executed)
//
// https://chessprogramming.wikispaces.com/Static+Exchange+Evaluation
// https://chessprogramming.wikispaces.com/SEE+-+The+Swap+Algorithm
//
// the implementation here is optimized for the common case when there
// isn't any capture following the move; the score returned is based
// on some fixed values for figures, different from the ones
// defined in material.go
func see(pos *Position, m Move) int32 {
	us := pos.Us()
	sq := m.To()
	bb := sq.Bitboard()
	target := m.Target() // piece in position
	bb27 := bb &^ (BbRank1 | BbRank8)
	bb18 := bb & (BbRank1 | BbRank8)

	var occ [ColorArraySize]Bitboard
	occ[White] = pos.ByColor(White)
	occ[Black] = pos.ByColor(Black)

	// occupancy tables as if moves are executed
	occ[us] &^= m.From().Bitboard()
	occ[us] |= m.To().Bitboard()
	occ[us.Opposite()] &^= m.CaptureSquare().Bitboard()
	us = us.Opposite()

	all := occ[White] | occ[Black]

	// adjust score for move
	score := seeScore(m)
	tmp := [16]int32{score}
	gain := tmp[:1]

	for score >= 0 {
		// try every figure in order of value
		var fig Figure                  // attacking figure
		var att Bitboard                // attackers
		var pawn, bishop, rook Bitboard // mobilies for our figures

		ours := occ[us]
		mt := Normal

		// pawn attacks
		pawn = Backward(us, West(bb27)|East(bb27))
		if att = pawn & ours & pos.ByFigure(Pawn); att != 0 {
			fig = Pawn
			goto makeMove
		}

		if att = KnightMobility(sq) & ours & pos.ByFigure(Knight); att != 0 {
			fig = Knight
			goto makeMove
		}

		if SuperQueenMobility(sq)&ours == 0 {
			// no other figure can attack sq so we give up early
			break
		}

		bishop = BishopMobility(sq, all)
		if att = bishop & ours & pos.ByFigure(Bishop); att != 0 {
			fig = Bishop
			goto makeMove
		}

		rook = RookMobility(sq, all)
		if att = rook & ours & pos.ByFigure(Rook); att != 0 {
			fig = Rook
			goto makeMove
		}

		// pawn promotions are considered queens minus the pawn
		pawn = Backward(us, West(bb18)|East(bb18))
		if att = pawn & ours & pos.ByFigure(Pawn); att != 0 {
			fig, mt = Queen, Promotion
			goto makeMove
		}

		if att = (rook | bishop) & ours & pos.ByFigure(Queen); att != 0 {
			fig = Queen
			goto makeMove
		}

		if att = KingMobility(sq) & ours & pos.ByFigure(King); att != 0 {
			fig = King
			goto makeMove
		}

		// no attack found
		break

	makeMove:
		// make a new pseudo-legal move of the smallest attacker
		from := att.LSB()
		attacker := ColorFigure(us, fig)
		m := MakeMove(mt, from.AsSquare(), sq, target, attacker)
		target = attacker // attacker becomes the new target

		// update score
		score = seeScore(m) - score
		gain = append(gain, score)

		// update occupancy tables for executing the move
		occ[us] = occ[us] &^ from
		all = all &^ from

		// switch sides
		us = us.Opposite()
	}

	for i := len(gain) - 2; i >= 0; i-- {
		if -gain[i+1] < gain[i] {
			gain[i] = -gain[i+1]
		}
	}
	return gain[0]
}

// seeSign return true if see(m) < 0
func seeSign(pos *Position, m Move) bool {
	if m.Piece().Figure() <= m.Capture().Figure() {
		// Even if m.Piece() is captured, we are still positive.
		return false
	}
	return see(pos, m) < 0
}

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// init

func init() {
	GlobalHashTable = NewHashTable(DefaultHashTableSizeMB)
}

/////////////////////////////////////////////////////////////////////
