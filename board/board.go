package board

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (b *Board) SetPieceAtSquare(sq Square, p Piece) bool {
	return b.Rep.SetPieceAtSquare(sq, p)
}

func (b *Board) PieceAtSquare(sq Square) Piece {
	return b.Rep.PieceAtSquare(sq)
}

func (b *Board) SetFromFen(fen string) {
	fenParts := strings.Split(fen, " ")

	b.Rep.SetFromFen(fenParts[0])

	b.Pos.Turn.SetFromFen(fenParts[1])

	b.Pos.CastlingRights.SetFromFen(fenParts[2], b)

	b.Pos.EpSquare = b.SquareFromAlgeb(fenParts[3])

	hmc, _ := strconv.ParseInt(fenParts[4], 10, 32)

	b.Pos.HalfmoveClock = int(hmc)

	fmn, _ := strconv.ParseInt(fenParts[5], 10, 32)

	b.Pos.FullmoveNumber = int(fmn)
}

func (b *Board) Line() string {
	buff := ""

	for i, msi := range b.MoveStack {
		if (i % 2) == 0 {
			buff += fmt.Sprintf("%d.", i/2+1)
		}

		buff += msi.San + " "
	}

	return buff
}

func (b *Board) ToString() string {
	buff := b.Rep.ToString()

	buff += "\n" + b.ReportFen() + "\n"

	buff += "\n" + b.Line() + "\n"

	return buff
}

func (b *Board) Print() {
	fmt.Println(b.ToString())
}

func (b *Board) Init(variant VariantKey) {
	// set variant
	b.Variant = variant

	// initialize rep to size required by variant
	b.Rep.Init(b.Variant)

	// init move stack
	b.MoveStack = make([]MoveStackItem, 0)

	// init position
	b.Pos.Init(b)
}

func (b *Board) HasSquare(sq Square) bool {
	return b.Rep.HasSquare(sq)
}

func (b *Board) ReportFen() string {
	buff := b.Rep.ReportFen()

	buff += " " + b.Pos.Turn.ToString()

	buff += " " + b.Pos.CastlingRights.ToString(b)

	buff += " " + b.SquareToAlgeb(b.Pos.EpSquare)

	buff += " " + fmt.Sprintf("%d", b.Pos.HalfmoveClock)

	buff += " " + fmt.Sprintf("%d", b.Pos.FullmoveNumber)

	return buff
}

func (b *Board) SquareToFileLetter(sq Square) string {
	return string([]byte{"a"[0] + byte(sq.File)})
}

func (b *Board) SquareToRankLetter(sq Square) string {
	return string([]byte{"1"[0] + byte(b.Rep.LastRank-sq.Rank)})
}

func (b *Board) SquareToAlgeb(sq Square) string {
	if sq.File < 0 {
		return "-"
	}
	return b.SquareToFileLetter(sq) + b.SquareToRankLetter(sq)
}

func (b *Board) SquareFromAlgeb(algeb string) Square {
	if algeb == "-" {
		return NO_SQUARE
	}

	return Square{int8(algeb[0] - "a"[0]), int8(byte(b.Rep.LastRank) - algeb[1] - "1"[0])}
}

func (b *Board) MoveToAlgeb(move Move) string {
	return b.SquareToAlgeb(move.FromSq) + b.SquareToAlgeb(move.ToSq)
}

func (b *Board) MoveToSan(move Move) string {
	checkStr := ""

	b.Push(move, !ADD_SAN)
	check := b.IsInCheck(b.Pos.Turn)
	if check {
		checkStr = "+"
		if !b.HasLegalMoveColor(b.Pos.Turn) {
			checkStr = "#"
		}
	}
	b.Pop()

	if move.Castling {
		if move.CastlingSide == QUEEN_SIDE {
			return "O-O-O" + checkStr
		}

		return "O-O" + checkStr
	}

	fromAlgeb := b.SquareToAlgeb(move.FromSq)
	toAlgeb := b.SquareToAlgeb(move.ToSq)
	fromPiece := b.PieceAtSquare(move.FromSq)
	pieceLetter := fromPiece.ToStringUpper()

	qualifier := ""

	if fromPiece.Kind != Pawn {
		pslAttacks := b.AttacksOnSquareByPiece(move.ToSq, fromPiece, ALL_ATTACKS)

		attacks := b.PickLegalMovesFrom(pslAttacks, b.Pos.Turn)

		files := make(map[int8]bool, 0)
		ranks := make(map[int8]bool, 0)
		samefiles := false
		sameranks := false

		if len(attacks) > 1 {
			for _, attack := range attacks {
				_, hasfile := files[attack.FromSq.File]
				if hasfile {
					samefiles = true
				} else {
					files[attack.FromSq.File] = true
				}

				_, hasrank := ranks[attack.FromSq.Rank]
				if hasrank {
					sameranks = true
				} else {
					ranks[attack.FromSq.Rank] = true
				}
			}

			if samefiles && sameranks {
				qualifier = fromAlgeb
			} else if samefiles {
				qualifier = fromAlgeb[1:2]
			} else {
				qualifier = fromAlgeb[0:1]
			}
		}
	}

	buff := pieceLetter + qualifier

	if fromPiece.Kind == Pawn {
		buff = ""
	}

	if move.IsCapture() {
		if fromPiece.Kind == Pawn {
			buff = b.SquareToFileLetter(move.FromSq)
		}
		buff += "x"
	}

	buff += toAlgeb

	if move.IsPromotion() {
		buff += "=" + move.PromotionPiece.ToStringUpper()
	}

	return buff + checkStr
}

func (b *Board) PslmsForVectorPieceAtSquare(p Piece, sq Square) []Move {
	pslms := make([]Move, 0)

	pdesc, ok := PIECE_KIND_TO_PIECE_DESCRIPTOR[p.Kind]

	if !ok {
		return pslms
	}

	currentSq := sq

	for _, dir := range pdesc.Directions {
		ok := true

		currentSq = sq.Add(dir)

		for ok {
			if b.HasSquare(currentSq) {
				top := b.PieceAtSquare(currentSq)

				capture := false
				add := true

				if !top.IsEmpty() {
					// non empty target square is capture
					capture = true

					if top.Color == p.Color {
						// cannot capture own piece
						add = false
					}
				}

				pslm := Move{
					FromSq: sq,
					ToSq:   currentSq,
				}

				if !pdesc.Sliding {
					ok = false
				}

				if capture {
					ok = false

					pslm.Capture = capture

					if !pdesc.CanCapture {
						add = false
					}
				}

				if add {
					pslms = append(pslms, pslm)
				}
			} else {
				ok = false
			}

			currentSq = currentSq.Add(dir)
		}
	}

	return pslms
}

func (b *Board) IsSquareEmpty(sq Square) bool {
	return b.Rep.IsSquareEmpty(sq)
}

func (b *Board) PawnRankDir(color PieceColor) int8 {
	// black pawn goes downward in rank
	var rankDir int8 = 1

	if color == WHITE {
		// white pawn goes upward in rank
		rankDir = -1
	}

	return rankDir
}

func (b *Board) AttacksOnSquareByPawn(sq Square, color PieceColor, stopAtFirst bool) []Move {
	attacks := make([]Move, 0)

	rdir := -b.PawnRankDir(color)

	var df int8
	for df = -1; df <= 1; df += 2 {
		testsq := sq.Add(PieceDirection{df, rdir})

		if b.HasSquare(testsq) {
			testp := b.PieceAtSquare(testsq)

			if (testp.Kind == Pawn) && (testp.Color == color) {
				attacks = append(attacks, Move{
					FromSq: testsq,
					ToSq:   sq,
				})

				if stopAtFirst {
					return attacks
				}
			}
		}
	}

	return attacks
}

func (b *Board) AttacksOnSquareByVectorPiece(sq Square, p Piece, stopAtFirst bool) []Move {
	attacks := make([]Move, 0)

	testp := p.ColorInverse()

	pslms := b.PslmsForVectorPieceAtSquare(testp, sq)

	for _, pslm := range pslms {
		if pslm.IsCapture() {
			testp := b.PieceAtSquare(pslm.ToSq)
			if testp.KindColorEqualTo(p) {
				attack := Move{
					FromSq: pslm.ToSq,
					ToSq:   pslm.FromSq,
				}

				attacks = append(attacks, attack)

				if stopAtFirst {
					return attacks
				}
			}
		}
	}

	return attacks
}

func (b *Board) AttacksOnSquareByPiece(sq Square, p Piece, stopAtFirst bool) []Move {
	if p.Kind == Pawn {
		return b.AttacksOnSquareByPawn(sq, p.Color, stopAtFirst)
	}

	return b.AttacksOnSquareByVectorPiece(sq, p, stopAtFirst)
}

func (b *Board) IsSquareAttackedByPiece(sq Square, p Piece) bool {
	attacks := b.AttacksOnSquareByPiece(sq, p, STOP_AT_FIRST)

	return len(attacks) > 0
}

func (b *Board) AttackingPieceKinds() []PieceKind {
	apks := []PieceKind{
		Pawn,
		King,
		Queen,
		Rook,
		Bishop,
		Knight,
	}

	if b.Variant == VARIANT_SEIRAWAN {
		apks = append(apks, []PieceKind{
			Elephant,
			Hawk,
		}...)
	}

	return apks
}

func (b *Board) IsSquareAttackedByColor(sq Square, color PieceColor) bool {
	apks := b.AttackingPieceKinds()

	for _, apk := range apks {
		if b.IsSquareAttackedByPiece(sq, Piece{Kind: apk, Color: color}) {
			return true
		}
	}

	return false
}

func (b *Board) IsInCheck(color PieceColor) bool {
	wk := b.WhereIsKing(color)

	if wk == NO_SQUARE {
		// missing king is considered check
		return true
	}

	return b.IsSquareAttackedByColor(wk, color.Inverse())
}

func (b *Board) PslmsForPawnAtSquare(p Piece, sq Square) []Move {
	pslms := make([]Move, 0)

	rankDir := b.PawnRankDir(p.Color)

	pushOneSq := sq.Add(PieceDirection{0, rankDir})

	if b.HasSquare(pushOneSq) {
		if b.IsSquareEmpty(pushOneSq) {
			if pushOneSq.Rank == b.PromotionRank(p.Color) {
				promotionMoves := b.CreatePromotionMoves(
					sq,        // from
					pushOneSq, // to
					false,     // pawn capture
					true,      // push by one
				)

				pslms = append(pslms, promotionMoves...)
			} else {
				move := Move{
					FromSq:        sq,
					ToSq:          pushOneSq,
					PawnPushByOne: true,
				}

				pslms = append(pslms, move)
			}

			pushTwoSq := pushOneSq.Add(PieceDirection{0, rankDir})

			if b.HasSquare(pushTwoSq) {
				if b.IsSquareEmpty(pushTwoSq) {
					epsq := NO_SQUARE

					var df int8
					for df = -1; df <= 1; df += 2 {
						testsq := pushTwoSq.Add(PieceDirection{df, 0})
						if b.HasSquare(testsq) {
							tp := b.PieceAtSquare(testsq)

							if (tp.Kind == Pawn) && (tp.Color != p.Color) {
								epsq = pushOneSq
							}
						}
					}

					plm := Move{
						FromSq:        sq,
						ToSq:          pushTwoSq,
						PawnPushByTwo: true,
						EpSquare:      epsq,
					}

					pslms = append(pslms, plm)
				}
			}
		}
	}

	var fileDir int8
	for fileDir = -1; fileDir <= 1; fileDir += 2 {
		captureSquare := sq.Add(PieceDirection{fileDir, rankDir})

		if b.HasSquare(captureSquare) {
			top := b.PieceAtSquare(captureSquare)

			if (top.Kind != NO_PIECE) && (top.Color != p.Color) {
				if pushOneSq.Rank == b.PromotionRank(p.Color) {
					promotionMoves := b.CreatePromotionMoves(
						sq,            // from
						captureSquare, // to
						true,          // pawn capture
						false,         // push by one
					)

					pslms = append(pslms, promotionMoves...)
				} else {
					plm := Move{
						FromSq:      sq,
						ToSq:        captureSquare,
						PawnCapture: true,
					}

					pslms = append(pslms, plm)
				}
			}

			if b.Pos.EpSquare == captureSquare {
				plm := Move{
					FromSq:        sq,
					ToSq:          captureSquare,
					EpCapture:     true,
					EpClearSquare: captureSquare.Add(PieceDirection{0, -rankDir}),
				}

				pslms = append(pslms, plm)
			}
		}
	}

	return pslms
}

func (b *Board) PslmsForPieceAtSquare(p Piece, sq Square) []Move {
	if p.Kind == Pawn {
		return b.PslmsForPawnAtSquare(p, sq)
	}

	return b.PslmsForVectorPieceAtSquare(p, sq)
}

func (b *Board) PslmsForAllPiecesOfColor(color PieceColor) []Move {
	pslms := make([]Move, 0)

	var rank int8
	var file int8
	for rank = 0; rank < b.Rep.NumRanks; rank++ {
		for file = 0; file < b.Rep.NumFiles; file++ {
			sq := Square{file, rank}
			p := b.PieceAtSquare(sq)
			if (p.Color == color) && (p.Kind != NO_PIECE) {
				pslms = append(pslms, b.PslmsForPieceAtSquare(p, sq)...)
			}
		}
	}

	wk := b.WhereIsKing(color)

	for side := QUEEN_SIDE; side <= KING_SIDE; side++ {
		cs := b.Pos.CastlingRights[color][side]

		if cs.Free(b) {
			move := Move{
				FromSq:        wk,
				ToSq:          cs.RookOrigSquare,
				Castling:      true,
				CastlingSide:  side,
				RookOrigPiece: cs.RookOrigPiece,
			}

			pslms = append(pslms, move)
		}
	}

	return pslms
}

func (b *Board) Reset() {
	fen, _ := START_FENS[b.Variant]
	b.SetFromFen(fen)
}

func (b *Board) MovesSortedBySan(moves []Move) MoveBuff {
	mb := make(MoveBuff, 0)

	for _, move := range moves {
		san := b.MoveToSan(move)
		algeb := b.MoveToAlgeb(move)

		mb = append(mb, MoveBuffItem{move, san, algeb})
	}

	sort.Sort(MoveBuff(mb))

	return mb
}

func (b *Board) CreatePromotionMoves(
	fromsq Square,
	tosq Square,
	pawnCapture bool,
	pawnPushByOne bool,
) []Move {
	promotionMoves := make([]Move, 0)

	promotionPieces, _ := PROMOTION_PIECES[b.Variant]

	for _, pp := range promotionPieces {
		promotionMove := Move{
			FromSq:         fromsq,
			ToSq:           tosq,
			PawnCapture:    pawnCapture,
			PawnPushByOne:  pawnPushByOne,
			PromotionPiece: pp,
		}

		promotionMoves = append(promotionMoves, promotionMove)
	}

	return promotionMoves
}

func (b *Board) Push(move Move, addSan bool) {
	san := "?"

	if addSan {
		san = b.MoveToSan(move)
	}

	restoreRep := make([]SetPiece, 0)
	oldPos := b.Pos.Clone()

	fromp := b.PieceAtSquare(move.FromSq)

	restoreRep = append(restoreRep, SetPiece{move.FromSq, fromp})

	top := b.PieceAtSquare(move.ToSq)

	restoreRep = append(restoreRep, SetPiece{move.ToSq, top})

	ccr := &b.Pos.CastlingRights[b.Pos.Turn]

	if fromp.Kind == King {
		ccr.ClearAll()
	}

	b.SetPieceAtSquare(move.FromSq, Piece{})

	if move.Castling {
		b.SetPieceAtSquare(move.ToSq, Piece{})
		kctsq := b.KingCastlingTargetSq(b.Pos.Turn, move.CastlingSide)
		kctp := b.PieceAtSquare(kctsq)
		restoreRep = append(restoreRep, SetPiece{kctsq, kctp})
		b.SetPieceAtSquare(kctsq, Piece{Kind: King, Color: b.Pos.Turn})
		rctsq := b.RookCastlingTargetSq(b.Pos.Turn, move.CastlingSide)
		rctp := b.PieceAtSquare(rctsq)
		restoreRep = append(restoreRep, SetPiece{rctsq, rctp})
		b.SetPieceAtSquare(rctsq, move.RookOrigPiece)
	} else {
		b.SetPieceAtSquare(move.ToSq, fromp)
	}

	var side CastlingSide
	for side = QUEEN_SIDE; side <= KING_SIDE; side++ {
		cs := &ccr[side]
		if cs.CanCastle {
			rp := b.PieceAtSquare(cs.RookOrigSquare)

			if !cs.RookOrigPiece.KindColorEqualTo(rp) {
				// rook changed, delete castling right
				cs.CanCastle = false
			}
		}
	}

	b.Pos.Turn = b.Pos.Turn.Inverse()

	b.Pos.EpSquare = NO_SQUARE

	if move.PawnPushByTwo {
		b.Pos.EpSquare = move.EpSquare
	}

	if move.EpCapture {
		b.SetPieceAtSquare(move.EpClearSquare, Piece{})
	}

	if move.ShouldDeleteHalfmoveClock() {
		b.Pos.HalfmoveClock = 0
	} else {
		b.Pos.HalfmoveClock++
	}

	if b.Pos.Turn == WHITE {
		b.Pos.FullmoveNumber++
	}

	b.MoveStack = append(b.MoveStack, MoveStackItem{
		restoreRep,
		oldPos,
		move,
		san,
	})
}

func (b *Board) Pop() {
	l := len(b.MoveStack)
	if l == 0 {
		return
	}

	msi := b.MoveStack[l-1]

	b.MoveStack = b.MoveStack[:l-1]

	for _, sp := range msi.RestoreRep {
		b.SetPieceAtSquare(sp.Sq, sp.P)
	}

	b.Pos = msi.Pos
}

func (b *Board) PromotionRank(color PieceColor) int8 {
	if color == WHITE {
		return 0
	}

	return 7
}

func (b *Board) CastlingRank(color PieceColor) int8 {
	if color == WHITE {
		return 7
	}

	return 0
}

func (b *Board) RookCastlingTargetSq(color PieceColor, side CastlingSide) Square {
	rank := b.CastlingRank(color)

	var file int8 = 2

	if side == KING_SIDE {
		file = 5
	}

	return Square{file, rank}
}

func (b *Board) KingCastlingTargetSq(color PieceColor, side CastlingSide) Square {
	rank := b.CastlingRank(color)

	var file int8 = 3

	if side == KING_SIDE {
		file = 6
	}

	return Square{file, rank}
}

func (b *Board) SquaresInDirection(sq Square, dir PieceDirection) []Square {
	return b.Rep.SquaresInDirection(sq, dir)
}

func (b *Board) WhereIsKing(color PieceColor) Square {
	return b.Rep.WhereIsKing(color)
}

func (b *Board) PickLegalMovesFrom(pslms []Move, color PieceColor) []Move {
	lms := make([]Move, 0)

	for _, pslm := range pslms {
		b.Push(pslm, !ADD_SAN)
		check := b.IsInCheck(color)
		b.Pop()

		if !check {
			lms = append(lms, pslm)
		}
	}

	return lms
}

func (b *Board) LegalMovesForAllPiecesOfColor(color PieceColor) []Move {
	pslms := b.PslmsForAllPiecesOfColor(color)

	lms := b.PickLegalMovesFrom(pslms, color)

	return lms
}

func (b *Board) HasLegalMoveColor(color PieceColor) bool {
	return len(b.LegalMovesForAllPiecesOfColor(color)) > 0
}

func (b *Board) LegalMovesForAllPieces() []Move {
	return b.LegalMovesForAllPiecesOfColor(b.Pos.Turn)
}

func (b *Board) PerfRecursive(depth int, maxDepth int) {
	b.Nodes++

	if depth > maxDepth {
		return
	}

	lms := b.LegalMovesForAllPieces()

	for _, lm := range lms {
		b.Push(lm, !ADD_SAN)
		b.PerfRecursive(depth+1, maxDepth)
		b.Pop()
	}
}

func (b *Board) StartPerf() {
	b.Nodes = 0

	b.Start = time.Now()
}

func (b *Board) StopPerf() {
	elapsed := float32(time.Now().Sub(b.Start)) / float32(1e9)

	nps := float32(b.Nodes) / float32(elapsed)

	fmt.Printf(">> perf elapsed %.2f nodes %d nps %.0f\n", elapsed, b.Nodes, nps)
}

func (b *Board) Perf(maxDepth int) {
	b.StartPerf()

	fmt.Printf(">> perf up to depth %d\n", maxDepth)

	b.PerfRecursive(0, maxDepth)
}

func (b *Board) Material(color PieceColor) int {
	material := 0

	var rank int8
	var file int8
	for rank = 0; rank < b.Rep.NumRanks; rank++ {
		for file = 0; file < b.Rep.NumFiles; file++ {
			p := b.PieceAtSquare(Square{file, rank})
			if (p.Color == color) && (p.Kind != NO_PIECE) {
				material += PIECE_VALUES[p.Kind]

				if p.Kind == Pawn {
					if (rank >= 3) && (rank <= 4) {
						if (file >= 3) && (file <= 4) {
							material += CENTER_PAWN_BONUS
						}
					}
				}

				if p.Kind == Knight {
					if (rank <= 1) || (rank >= 6) {
						material -= KNIGHT_ON_EDGE_DEDUCTION
					}
					if (file <= 1) || (file >= 6) {
						material -= KNIGHT_ON_EDGE_DEDUCTION
					}
					if (rank <= 0) || (rank >= 7) {
						material -= KNIGHT_CLOSE_TO_EDGE_DEDUCTION
					}
					if (file <= 0) || (file >= 7) {
						material -= KNIGHT_CLOSE_TO_EDGE_DEDUCTION
					}
				}
			}
		}
	}

	pslms := b.PslmsForAllPiecesOfColor(color)

	material += MOBILITY_MULTIPLIER * len(pslms)

	material += rand.Intn(RANDOM_BONUS)

	return material
}

func (b *Board) MaterialBalance() int {
	return b.Material(WHITE) - b.Material(BLACK)
}

func (b *Board) Eval() int {
	return b.MaterialBalance()
}

func (b *Board) EvalForColor(color PieceColor) int {
	eval := b.Eval()

	if color == WHITE {
		return eval
	}

	return -eval
}

func (b *Board) EvalForTurn() int {
	return b.EvalForColor(b.Pos.Turn)
}

// https://www.chessprogramming.org/Alpha-Beta
func (b *Board) AlphaBeta(alpha int, beta int, depthLeft int, quescDepth int) (Move, int) {
	b.Nodes++

	bm := Move{}

	if depthLeft <= -quescDepth {
		return bm, b.EvalForTurn()
	}

	lms := b.LegalMovesForAllPieces()

	hadMove := false

	for _, lm := range lms {
		if (depthLeft >= 0) || lm.IsCapture() {
			b.Push(lm, !ADD_SAN)
			_, score := b.AlphaBeta(-beta, -alpha, depthLeft-1, quescDepth)
			b.Pop()

			hadMove = true

			score *= -1

			if score >= beta {
				return bm, beta
			}

			if score > alpha {
				bm = lm
				alpha = score
			}
		}
	}

	if !hadMove {
		return bm, b.EvalForTurn()
	}

	return bm, alpha
}

func (b *Board) Go(depth int, quisecenceDepth int) (Move, int) {
	b.StartPerf()

	fmt.Printf(">> go depth %d quiescence depth %d\n", depth, quisecenceDepth)

	bm, score := b.AlphaBeta(-10000, 10000, depth, quisecenceDepth)

	b.StopPerf()

	fmt.Printf(">> move %s score %d\n", b.MoveToSan(bm), score)

	return bm, score
}

/////////////////////////////////////////////////////////////////////
