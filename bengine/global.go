package bengine

/////////////////////////////////////////////////////////////////////
// imports

import (
	"sync"
	"unsafe"

	. "github.com/easychessanimations/gochess/butils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// global variables

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

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// init

func init() {
	GlobalHashTable = NewHashTable(DefaultHashTableSizeMB)
}

/////////////////////////////////////////////////////////////////////
