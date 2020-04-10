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

// Weights stores the network parameters
// network has train error 0.05679009 and validation error 0.05702872
var Weights = []Score{
	{M: -2, E: 0}, {M: 14364, E: 13854}, {M: 64220, E: 49415}, {M: 66592, E: 52025}, {M: 84852, E: 98328}, {M: 222394, E: 169927}, {M: -160, E: 130}, {M: 1199, E: 2146},
	{M: -6958, E: 77}, {M: -9864, E: 2079}, {M: -2773, E: 114}, {M: -1505, E: 100}, {M: -4315, E: -2085}, {M: -708, E: -88}, {M: 7, E: 1901}, {M: 1161, E: 3481},
	{M: 1977, E: 2787}, {M: 2058, E: 1354}, {M: 776, E: -44}, {M: -1118, E: -2678}, {M: -534, E: -2756}, {M: 10, E: 25}, {M: 2200, E: 1697}, {M: 4187, E: 3431},
	{M: 5955, E: 3342}, {M: 4977, E: 1044}, {M: -157, E: -390}, {M: -13491, E: -1993}, {M: 1590, E: 424}, {M: -3039, E: -94}, {M: 991, E: -43}, {M: 1121, E: 222},
	{M: -726, E: 1148}, {M: 0, E: 1220}, {M: 310, E: 698}, {M: 3442, E: -908}, {M: -1059, E: 13}, {M: -262, E: -338}, {M: 1134, E: -76}, {M: 2662, E: 745},
	{M: 645, E: 767}, {M: 1004, E: 374}, {M: -66, E: -84}, {M: -7848, E: 1101}, {M: -6984, E: -177}, {M: 1581, E: 667}, {M: -978, E: 386}, {M: -231, E: 594},
	{M: 326, E: 867}, {M: 2005, E: 19}, {M: 2387, E: -796}, {M: 3694, E: -1037}, {M: -969, E: 199}, {M: -6, E: -1217}, {M: 2252, E: -2485}, {M: -550, E: -1226},
	{M: 48, E: -1192}, {M: -572, E: 63}, {M: 52, E: 407}, {M: 101, E: 600}, {M: 1303, E: 1319}, {M: -787, E: 1963}, {M: 1339, E: 470}, {M: 8194, E: -1609},
	{M: 1878, E: 2066}, {M: 418, E: 1069}, {M: -555, E: -63}, {M: -709, E: -537}, {M: -258, E: -208}, {M: 5, E: 87}, {M: 539, E: -13}, {M: 2403, E: 67},
	{M: 2004, E: 1769}, {M: 13247, E: -6786}, {M: 9471, E: -6878}, {M: 5311, E: -2926}, {M: 414, E: 152}, {M: -3468, E: 2329}, {M: -601, E: -98}, {M: -6549, E: 3325},
	{M: 36, E: 921}, {M: 529, E: 1485}, {M: -2851, E: -2370}, {M: 2150, E: 4431}, {M: 5406, E: 8787}, {M: -3138, E: 783}, {M: -1358, E: 7}, {M: 16986, E: -4302},
	{M: 29780, E: 85}, {M: -86, E: -20}, {M: 71, E: -40}, {M: 1, E: -14}, {M: -76, E: 16}, {M: -126, E: 19}, {M: 38, E: 39}, {M: -59, E: 140},
	{M: 66, E: 107}, {M: -2370, E: 15}, {M: -2166, E: -644}, {M: -1993, E: 537}, {M: -1301, E: -139}, {M: -2122, E: 1284}, {M: 3850, E: -495}, {M: 2909, E: -1932},
	{M: -1674, E: -2480}, {M: -2015, E: -1135}, {M: -4037, E: -451}, {M: -893, E: -1554}, {M: -2342, E: -239}, {M: -923, E: -112}, {M: -451, E: -425}, {M: 1343, E: -2536},
	{M: -1497, E: -1896}, {M: -1661, E: 1581}, {M: -3958, E: 841}, {M: 30, E: -982}, {M: 675, E: -1662}, {M: 284, E: -711}, {M: 262, E: -1286}, {M: -2817, E: -301},
	{M: -2886, E: 2}, {M: 98, E: 2937}, {M: -9, E: 1467}, {M: 103, E: 414}, {M: 1561, E: -187}, {M: 998, E: 319}, {M: 90, E: 506}, {M: -10, E: 1856},
	{M: -1247, E: 2380}, {M: 4373, E: 5274}, {M: 954, E: 5331}, {M: 3116, E: 3809}, {M: 1007, E: 1454}, {M: 7494, E: 205}, {M: 12094, E: 927}, {M: 2921, E: 4545},
	{M: 2135, E: 5350}, {M: 123, E: 253}, {M: 126, E: 108}, {M: -97, E: -99}, {M: 44, E: -1229}, {M: 27, E: 1017}, {M: -13, E: -838}, {M: -206, E: -5},
	{M: -115, E: 1456}, {M: 226, E: -82}, {M: 38, E: 94}, {M: -124, E: -36}, {M: 157, E: -82}, {M: 87, E: -97}, {M: -69, E: 123}, {M: -105, E: 159},
	{M: 172, E: 40}, {M: -2908, E: -2023}, {M: 2162, E: 958}, {M: -1525, E: 45}, {M: -1515, E: -1395}, {M: 2855, E: 757}, {M: -1183, E: -4363}, {M: 3153, E: -1831},
	{M: 566, E: 293}, {M: -9013, E: 2241}, {M: -520, E: 1088}, {M: -4066, E: 2623}, {M: 3714, E: -20}, {M: 1372, E: -3061}, {M: -3598, E: -2663}, {M: -5, E: -938},
	{M: 11, E: 62}, {M: -3350, E: -10}, {M: -33, E: 352}, {M: 6186, E: -18}, {M: 22, E: 94}, {M: -37, E: -2029}, {M: -775, E: 596}, {M: 565, E: -568},
	{M: 3402, E: -1948}, {M: 2374, E: -66}, {M: -123, E: 142}, {M: 202, E: 2363}, {M: 260, E: 2843}, {M: -688, E: 6686}, {M: 49, E: 3610}, {M: 47, E: 17042},
	{M: 9455, E: 33952}, {M: -24, E: -176}, {M: 0, E: 33}, {M: -1108, E: 10140}, {M: -7291, E: 5046}, {M: -2197, E: -4209}, {M: 136, E: -12279}, {M: 3, E: -13911},
	{M: -28, E: -14428}, {M: -50, E: -12908}, {M: 147, E: -85}, {M: 4548, E: 12736}, {M: -43, E: 7638}, {M: 151, E: 1756}, {M: -2061, E: 273}, {M: -285, E: -218},
	{M: 1997, E: -64}, {M: 226, E: 60},
}

var (
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

// hash table
func init() {
	GlobalHashTable = NewHashTable(DefaultHashTableSizeMB)
}

// distance
func init() {
	for i := SquareMinValue; i <= SquareMaxValue; i++ {
		for j := SquareMinValue; j <= SquareMaxValue; j++ {
			f, r := int32(i.File()-j.File()), int32(i.Rank()-j.Rank())
			f, r = max(f, -f), max(r, -r) // absolute value
			distance[i][j] = max(f, r)
		}
	}
}

/////////////////////////////////////////////////////////////////////
