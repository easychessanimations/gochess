package bengine

import (
	"testing"

	. "github.com/easychessanimations/gochess/butils"
)

func TestPieceSize(t *testing.T) {
	shouldBe := (1 << PIECE_ARRAY_SIZE_IN_BITS)
	if PieceArraySize != shouldBe {
		t.Errorf("( PieceArraySize = %d ) != ( 2 ^ ( PIECE_ARRAY_SIZE_IN_BITS = %d ) = %d )", PieceArraySize, PIECE_ARRAY_SIZE_IN_BITS, shouldBe)
	}
}

/*
func TestGame(t *testing.T) {
	pos, _ := PositionFromFEN(FENStartPos)
	eng := NewEngine(pos, nil, Options{})
	for i := 0; i < 1; i++ {
		tc := NewFixedDepthTimeControl(pos, 3)
		tc.Start(false)
		_, pv := eng.Play(tc)
		eng.DoMove(pv[0])
	}
}

func TestMateIn1(t *testing.T) {
	for i, d := range MateIn1 {
		pos, _ := PositionFromFEN(d.FEN)
		bm, err := pos.UCIToMove(d.BM)
		if err != nil {
			t.Errorf("#%d cannot parse move %s", i, d.BM)
			continue
		}

		tc := NewFixedDepthTimeControl(pos, 2)
		tc.Start(false)
		eng := NewEngine(pos, nil, Options{})
		_, pv := eng.Play(tc)

		if len(pv) != 1 {
			t.Errorf("#%d Expected at most one move, got %d", i, len(pv))
			t.Errorf("position is %v", pos)
			continue
		}

		if pv[0] != bm {
			t.Errorf("#%d expected move %v, got %v", i, bm, pv[0])
			t.Errorf("position is %v", pos)
			continue
		}
	}
}

// Test score is the same if we start with the position or move.
func TestScore(t *testing.T) {
	for _, game := range TestGames {
		pos, _ := PositionFromFEN(FENStartPos)
		dynamic := NewEngine(pos, nil, Options{})
		static := NewEngine(pos, nil, Options{})

		moves := strings.Fields(game)
		for _, move := range moves {
			m, _ := pos.UCIToMove(move)
			if !pos.IsPseudoLegal(m) {
				// t.Fatalf("bad bad bad")
			}

			dynamic.DoMove(m)
			static.SetPosition(pos)
			if dynamic.Score() != static.Score() {
				t.Fatalf("expected static score %v, got dynamic score %v", static.Score(), dynamic.Score())
			}
		}
	}
}

func TestEndGamePosition(t *testing.T) {
	pos, _ := PositionFromFEN("6k1/5p1p/4p1p1/3p4/5P1P/8/3r2q1/6K1 w - - 2 55")
	tc := NewFixedDepthTimeControl(pos, 3)
	tc.Start(false)
	eng := NewEngine(pos, nil, Options{})
	_, pv := eng.Play(tc)
	if pv != nil {
		t.Errorf("got %d moves (nonil, pv), expected nil pv", len(pv))
	}
}

// pvLogger logs the PV.
// It will panic if pvs are not in order.
type pvLog struct {
	depth   int32
	multiPV int
	score   int32
	moves   []Move
}

type pvLogger []pvLog

func (l *pvLogger) BeginSearch()                           {}
func (l *pvLogger) EndSearch()                             {}
func (l *pvLogger) CurrMove(depth int, move Move, num int) {}

func (l *pvLogger) PrintPV(stats Stats, multiPV int, score int32, moves []Move) {
	*l = append(*l, pvLog{
		depth:   stats.Depth,
		multiPV: multiPV,
		score:   score,
		moves:   moves,
	})
}

func TestMultiPV(t *testing.T) {
	for f, fen := range TestFENs {
		pos, _ := PositionFromFEN(fen)
		tc := NewFixedDepthTimeControl(pos, 4)
		tc.Start(false)
		pvl := pvLogger{}
		eng := NewEngine(pos, &pvl, Options{MultiPV: 3})
		eng.Play(tc)

		// Check the number of iterations.
		numIterations := 0
		for i := range pvl {
			if pvl[i].multiPV == 1 {
				numIterations++
			}
		}
		if numIterations != 4+1 {
			t.Errorf("#%d %s: expected 4+1 iterations, got %d", f, fen, numIterations)
		}

		// Check score and depth order.
		for i := 1; i < len(pvl); i++ {
			if pvl[i-1].depth > pvl[i].depth {
				// TODO: this is not really correct if we repeat the PVS lines
				t.Errorf("#%d %s: wrong depth order", f, fen)
			}
			if pvl[i-1].depth == pvl[i].depth && pvl[i-1].score < pvl[i].score {
				t.Errorf("#%d %s: wrong score order", f, fen)
			}
		}

		// Check different moves for the same iterations.
		for i := range pvl {
			for j := range pvl {
				if i <= j {
					continue
				}
				if pvl[i].depth != pvl[j].depth || pvl[i].multiPV == pvl[j].multiPV {
					continue
				}
				if len(pvl[i].moves) == 0 || len(pvl[j].moves) == 0 {
					continue
				}
				if pvl[i].moves[0] == pvl[j].moves[0] {
					t.Errorf("#%d %s: got identical moves", f, fen)
				}
			}
		}
	}
}

func BenchmarkGame(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pos, _ := PositionFromFEN(FENStartPos)
		eng := NewEngine(pos, nil, Options{})
		for j := 0; j < 20; j++ {
			tc := NewFixedDepthTimeControl(pos, 4)
			tc.Start(false)
			_, pv := eng.Play(tc)
			if len(pv) > 0 {
				eng.DoMove(pv[0])
			}
		}
	}
}

func TestScoreRange(t *testing.T) {
	for _, fen := range TestFENs {
		pos, _ := PositionFromFEN(fen)
		score := Evaluate(pos).GetCentipawnsScore()
		if KnownLossScore >= score || score >= KnownWinScore {
			t.Errorf("expected %d in interval (%d, %d) for %s",
				score, KnownLossScore, KnownWinScore, fen)
		}
	}
}

func BenchmarkScore(b *testing.B) {
	for _, fen := range TestFENs {
		pos, _ := PositionFromFEN(fen)
		for i := 0; i < b.N; i++ {
			Evaluate(pos)
		}
	}
}

func TestDistance(t *testing.T) {
	data := []struct {
		i, j Square
		d    int32
	}{
		{SquareA1, SquareA8, 7},
		{SquareA1, SquareH8, 7},
		{SquareB2, SquareB2, 0},
		{SquareB2, SquareC3, 1},
		{SquareE5, SquareD4, 1},
		{SquareE5, SquareD4, 1},
		{SquareE1, SquareG5, 4},
	}

	for i, d := range data {
		if got, want := distance[d.i][d.j], d.d; got != want {
			t.Errorf("#%d wanted distance[%v][%v] == %d, got %d", i, d.i, d.j, want, got)
		}
	}
}

func TestMurmurMixSwap(t *testing.T) {
	c1 := uint64(3080512559332270987)
	c2 := uint64(1670079002898303149)

	h1 := murmurSeed[NoFigure]
	h1 = murmurMix(h1, c1)
	h1 = murmurMix(h1, c2)

	h2 := murmurSeed[NoFigure]
	h2 = murmurMix(h2, c2)
	h2 = murmurMix(h2, c1)

	if h1 == h2 {
		t.Errorf("murmurMix(c1, c2) == murmurMix(c2, c1) (%d, %d), wanted different", h1, h2)
	}
}

func TestOrdersViolentMovesByMVVLVA(t *testing.T) {
	for _, fen := range TestFENs {
		pos, _ := PositionFromFEN(fen)
		st := &stack{history: new(historyTable)}
		st.Reset(pos)
		st.GenerateMoves(Violent, NullMove)

		limit := int16(0x7fff)
		for move := st.PopMove(); move != NullMove; move = st.PopMove() {
			if curr := mvvlva(move); curr > limit {
				t.Errorf("moves not sorted: %v", move)
			} else {
				limit = curr
			}
		}
	}
}

func TestReturnsHashMove(t *testing.T) {
	pos, _ := PositionFromFEN(FENKiwipete)

	for i, str := range []string{"f3f5", "e2b5", "a1b1"} {
		hash, _ := pos.UCIToMove(str)
		st := &stack{}
		st.Reset(pos)
		st.GenerateMoves(Violent|Quiet, hash)
		if move := st.PopMove(); hash != move {
			t.Errorf("#%d expected move %v, got %v", i, hash, move)
		}
	}
}

func TestReturnsMoves(t *testing.T) {
	for _, fen := range TestFENs {
		pos, _ := PositionFromFEN(fen)
		seen := make(map[Move]int)

		var moves []Move
		pos.GenerateMoves(Violent|Quiet, &moves)
		for _, m := range moves {
			seen[m] |= 1
		}

		st := &stack{history: new(historyTable)}
		st.Reset(pos)
		st.GenerateMoves(Violent|Quiet, moves[1234567891%len(moves)])
		for m := st.PopMove(); m != NullMove; m = st.PopMove() {
			if seen[m]&2 != 0 {
				t.Errorf("move %v is duplicate: %s", m, fen)
			}
			seen[m] |= 2
		}

		for m, v := range seen {
			if v == 1 {
				t.Errorf("move %v not generated", m)
			}
			if v == 2 {
				t.Errorf("move %v not expected", m)
			}
		}
	}
}

func TestCachePutGet(t *testing.T) {
	c1 := uint64(3080512559332270987)
	c2 := uint64(1670079002898303149)

	h1 := murmurSeed[NoFigure]
	h1 = murmurMix(h1, c1)
	h1 = murmurMix(h1, c2)

	ew, eb := Accum{1, 2}, Accum{3, 5}
	c := new(pawnsTable)
	c.put(h1, ew, eb)
	if gw, gb, ok := c.get(h1); !ok {
		t.Errorf("entry not in the cache, expecting a git")
	} else if ew != gw || eb != gb {
		t.Errorf("got get(%d) == %v, %v; wanted %v. %v", h1, gw, gb, ew, eb)
	}

	h2 := murmurSeed[NoFigure]
	h2 = murmurMix(h2, c2)
	h2 = murmurMix(h2, c1)
	if _, _, ok := c.get(h2); ok {
		t.Errorf("entry in the cache, expecting a miss")
	}
}

func TestPV(t *testing.T) {
	pos, _ := PositionFromFEN(FENStartPos)
	pvTable := newPvTable()
	for _, game := range TestGames {
		var moves []Move
		movesStr := strings.Fields(game)
		for _, moveStr := range movesStr {
			move, _ := pos.UCIToMove(moveStr)
			pos.DoMove(move)
			moves = append(moves, move)
		}

		for i := len(moves) - 1; i >= 0; i-- {
			pos.UndoMove()
			pvTable.Put(pos, moves[i])
		}

		pv := pvTable.Get(pos)
		if len(pv) == 0 {
			t.Errorf("expected at least on move on principal variation")
		}
		if len(pv) > len(moves) {
			// This can actually happen during the game.
			t.Errorf("got more moves on pv than in the game")
		}
		for i := range pv {
			if moves[i] != pv[i] {
				t.Errorf("#%d Expected move %v, got %v", i, pv[i], moves[i])
			}
		}
	}
}

func seeSlow(pos *Position, m Move, score int32) int32 {
	if m == NullMove || score > 0 {
		return score
	}

	// Compute the score change.
	score += seeScore(m)

	// Find the smallest attacker.
	var moves []Move
	pos.GenerateMoves(Violent, &moves)
	next := NullMove
	for _, n := range moves {
		if n.To() != m.To() {
			continue
		}

		// If the move is a promotion, consider the attacker to be a queen.
		fig, sq := n.Target().Figure(), n.From()
		if next == NullMove || fig < next.Target().Figure() || (fig == next.Piece().Figure() && sq < next.From()) {
			next = n
		}
	}

	// Recursively compute the see.
	pos.DoMove(next)
	see := -seeSlow(pos, next, -score)
	pos.UndoMove()

	if see > score {
		return score
	}
	return see
}

func TestSEE(t *testing.T) {
	good, bad := 0, 0
	for i, fen := range TestFENs {
		var moves []Move
		pos, _ := PositionFromFEN(fen)
		pos.GenerateMoves(Violent|Quiet, &moves)
		for _, m := range moves {
			actual := see(pos, m)
			pos.DoMove(m)
			expected := seeSlow(pos, m, 0)
			pos.UndoMove()

			if expected != actual {
				t.Errorf("#%d expected %d, got %d\nfor %v on %v", i, expected, actual, m, fen)
				bad++
			} else {
				good++
			}
		}
	}

	if bad != 0 {
		t.Errorf("Failed %d out of %d", bad, good+bad)
	}
}

// A benchmark position from http://www.stmintz.com/ccc/index.php?id=60880
var seeBench = "1rr3k1/4ppb1/2q1bnp1/1p2B1Q1/6P1/2p2P2/2P1B2R/2K4R w - - 0 1"

func BenchmarkSEESlow(b *testing.B) {
	var moves []Move
	pos, _ := PositionFromFEN(seeBench)
	pos.GenerateMoves(Violent|Quiet, &moves)
	for i := 0; i < b.N; i++ {
		for _, m := range moves {
			seeSlow(pos, m, 0)
		}
	}
}

func BenchmarkSEEFast(b *testing.B) {
	var moves []Move
	pos, _ := PositionFromFEN(seeBench)
	pos.GenerateMoves(Violent|Quiet, &moves)
	for i := 0; i < b.N; i++ {
		for _, m := range moves {
			see(pos, m)
		}
	}
}
*/
