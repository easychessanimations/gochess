package board

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/easychessanimations/gochess/utils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

var DEBUG = false

func (b *Board) IS_ATOMIC() bool {
	return b.Variant == utils.VARIANT_ATOMIC
}

func (b *Board) IS_EIGHTPIECE() bool {
	return b.Variant == utils.VARIANT_EIGHTPIECE
}

func (b *Board) IsExploded(color utils.PieceColor) bool {
	wk := b.WhereIsKing(color)

	return wk == utils.NO_SQUARE
}

func (b *Board) AdjacentSquares(sq utils.Square) []utils.Square {
	asqs := []utils.Square{}

	var df int8
	var dr int8
	for df = -1; df <= 1; df++ {
		for dr = -1; dr <= 1; dr++ {
			if (df != 0) || (dr != 0) {
				testsq := sq.Add(utils.PieceDirection{df, dr})
				if b.HasSquare(testsq) {
					asqs = append(asqs, testsq)
				}
			}
		}
	}

	return asqs
}

func (b *Board) RookAdjacentSquares(sq utils.Square) []utils.Square {
	rasqs := []utils.Square{}

	var df int8
	var dr int8
	for df = -1; df <= 1; df++ {
		for dr = -1; dr <= 1; dr++ {
			if (df*df + dr*dr) == 1 {
				testsq := sq.Add(utils.PieceDirection{df, dr})
				if b.HasSquare(testsq) {
					rasqs = append(rasqs, testsq)
				}
			}
		}
	}

	return rasqs
}

func (b *Board) IsSquareJailedForColor(sq utils.Square, color utils.PieceColor) bool {
	rasqs := b.RookAdjacentSquares(sq)

	for _, rasq := range rasqs {
		p := b.PieceAtSquare(rasq)

		if (p.Kind == utils.Jailer) && (p.Color == color.Inverse()) {
			return true
		}
	}

	return false
}

func (b *Board) GetUciOptionByNameWithDefault(name string, uciOption utils.UciOption) utils.UciOption {
	if b.GetUciOptionByNameWithDefaultFunc != nil {
		return b.GetUciOptionByNameWithDefaultFunc(name, uciOption)
	}

	return uciOption
}

func (b *Board) SetPieceAtSquare(sq utils.Square, p utils.Piece) bool {
	if b.HasSquare(sq) {
		b.Pos.Rep[sq.Rank][sq.File] = p

		return true
	}

	return false
}

func (b *Board) PieceAtSquare(sq utils.Square) utils.Piece {
	if b.HasSquare(sq) {
		return b.Pos.Rep[sq.Rank][sq.File]
	}

	return utils.NO_PIECE
}

func (b *Board) Line() string {
	buff := "Line : "

	for i, msi := range b.MoveStack {
		if (i % 2) == 0 {
			buff += fmt.Sprintf("%d.", i/2+1)
		}

		buff += msi.San + " "
	}

	return buff
}

func (b *Board) ReportMaterial() string {
	materialWhite, materialBonusWhite, mobilityWhite := b.Material(utils.WHITE)
	materialBlack, materialBonusBlack, mobilityBlack := b.Material(utils.BLACK)

	totalWhite := materialWhite + materialBonusWhite + mobilityWhite
	totalBlack := materialBlack + materialBonusBlack + mobilityBlack

	buff := fmt.Sprintf("%-10s %10s %10s %10s %10s\n", "", "material", "mat. bonus", "mobility", "total")
	buff += fmt.Sprintf("%-10s %10d %10d %10d %10d\n", "white", materialWhite, materialBonusWhite, mobilityWhite, totalWhite)
	buff += fmt.Sprintf("%-10s %10d %10d %10d %10d\n", "black", materialBlack, materialBonusBlack, mobilityBlack, totalBlack)
	buff += fmt.Sprintf("%-10s %10d %10d %10d %10d", "balance", materialWhite-materialBlack, materialBonusWhite-materialBonusBlack, mobilityWhite-mobilityBlack, totalWhite-totalBlack)

	return buff
}

func (b *Board) AlgebToMoveRaw(algeb string) utils.Move {
	if algeb == "-" {
		return NO_MOVE
	}

	fromsq := b.SquareFromAlgeb(algeb[0:2])
	tosq := b.SquareFromAlgeb(algeb[2:4])

	move := utils.Move{
		FromSq: fromsq,
		ToSq:   tosq,
	}

	return move
}

func (b *Board) AlgebToMove(algeb string) utils.Move {
	lms := b.LegalMovesForAllPieces()

	for _, lm := range lms {
		if b.MoveToAlgeb(lm) == algeb {
			return lm
		}
	}

	return NO_MOVE
}

func (b *Board) MakeAlgebMove(algeb string, addSan bool) {
	move := b.AlgebToMove(algeb)

	if move != NO_MOVE {
		b.Push(move, addSan)
	}
}

func (b *Board) ExecCommand(command string) bool {
	b.SortedSanMoveBuff = b.MovesSortedBySan(b.LegalMovesForAllPieces())

	i, err := strconv.ParseInt(command, 10, 32)

	if err == nil {
		move := b.SortedSanMoveBuff[i-1].Move

		b.Push(move, ADD_SAN)

		b.Print()

		return true
	} else {
		if command == "g" {
			bm, _ := b.Go(10)

			b.Push(bm, ADD_SAN)

			b.Print()

			return true
		} else if command == "perf" {
			b.Perf(3)

			return true
		} else if command == "d" {
			b.Pop()

			b.Print()

			return true
		} else if command == "" {
			randIndex := rand.Intn(len(b.SortedSanMoveBuff))

			move := b.SortedSanMoveBuff[randIndex-1].Move

			b.Push(move, ADD_SAN)

			b.Print()

			return true
		} else if command != "" {
			for _, mbi := range b.SortedSanMoveBuff {
				if (mbi.San == command) || (mbi.Algeb == command) {
					move := mbi.Move

					b.Push(move, ADD_SAN)

					b.Print()

					return true
				}
			}
		}
	}

	return false
}

func (b *Board) ToString() string {
	buff := ""

	var rank int8
	var file int8
	for rank = 0; rank < b.NumRanks; rank++ {
		for file = 0; file < b.NumFiles; file++ {
			piece := b.Pos.Rep[rank][file]
			buff += fmt.Sprintf("%-4s", piece.ToString())
		}
		buff += "\n"
	}

	buff += "\n" + utils.VariantKeyToVariantKeyString(b.Variant)
	buff += "\n" + b.ReportFen() + "\n"

	buff += "\n" + b.ReportMaterial() + "\n"

	buff += fmt.Sprintf("\nEval for turn : %d\n", b.EvalForTurn())

	buff += "\n" + b.Line() + "\n"

	buff += "\n" + b.LegalMovesToString()

	return buff
}

func (b *Board) Log(content string) {
	if b.LogFunc != nil {
		b.LogFunc(content)
	} else {
		fmt.Println(content)
	}
}

func (b *Board) LogAnalysisInfo(content string) {
	if b.LogAnalysisInfoFunc != nil {
		b.LogAnalysisInfoFunc(content)
	} else {
		fmt.Println(content)
	}
}

func (b *Board) LegalMovesToString() string {
	lms := b.LegalMovesForAllPieces()

	b.SortedSanMoveBuff = b.MovesSortedBySan(lms)

	buff := fmt.Sprintf("Legal moves ( %d ) ", len(lms))

	for i, mbi := range b.SortedSanMoveBuff {
		buff += fmt.Sprintf("%d. %s [ %s ] ", i+1, mbi.San, mbi.Algeb)
	}

	return buff
}

func (b *Board) Print() {
	b.Log(b.ToString())

	/*DEBUG = true
	var rank int8
	var file int8
	for rank = 0; rank < b.NumRanks; rank++ {
		for file = 0; file < b.NumFiles; file++ {
			sq := utils.Square{file, rank}
			fmt.Println(sq, len(b.AttacksOnSquareBySentry(sq, utils.WHITE, ALL_ATTACKS)))
		}
	}
	DEBUG = false*/
}

func (b *Board) HasSquare(sq utils.Square) bool {
	return (sq.File >= 0) && (sq.File < b.NumFiles) && (sq.Rank >= 0) && (sq.Rank < b.NumRanks)
}

func (b *Board) SquareToFileLetter(sq utils.Square) string {
	return string([]byte{"a"[0] + byte(sq.File)})
}

func (b *Board) SquareToRankLetter(sq utils.Square) string {
	return string([]byte{"1"[0] + byte(b.LastRank-sq.Rank)})
}

func (b *Board) SquareToAlgeb(sq utils.Square) string {
	if sq.File < 0 {
		return "-"
	}
	return b.SquareToFileLetter(sq) + b.SquareToRankLetter(sq)
}

func (b *Board) SquareFromAlgeb(algeb string) utils.Square {
	if algeb == "-" {
		return utils.NO_SQUARE
	}

	return utils.Square{int8(algeb[0] - "a"[0]), int8(byte(b.LastRank) - algeb[1] - "1"[0])}
}

func (b *Board) MoveToAlgeb(move utils.Move) string {
	if move == NO_MOVE {
		return "-"
	}

	return b.SquareToAlgeb(move.FromSq) + b.SquareToAlgeb(move.ToSq)
}

func (b *Board) MoveToSan(move utils.Move) string {
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
		if move.CastlingSide == utils.QUEEN_SIDE {
			return "O-O-O" + checkStr
		}

		return "O-O" + checkStr
	}

	fromAlgeb := b.SquareToAlgeb(move.FromSq)
	toAlgeb := b.SquareToAlgeb(move.ToSq)
	fromPiece := b.PieceAtSquare(move.FromSq)
	pieceLetter := fromPiece.ToStringUpper()

	qualifier := ""

	if fromPiece.Kind != utils.Pawn {
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

	if fromPiece.Kind == utils.Pawn {
		buff = ""
	}

	if move.IsCapture() {
		if fromPiece.Kind == utils.Pawn {
			buff = b.SquareToFileLetter(move.FromSq)
		}
		buff += "x"
	}

	buff += toAlgeb

	if move.IsPromotion() {
		buff += "=" + move.PromotionPiece.ToStringUpper()

		if move.PromotionSquare != utils.NO_SQUARE {
			buff += "@" + b.SquareToAlgeb(move.PromotionSquare)
		}
	}

	return buff + checkStr
}

func (b *Board) IsSquareEmpty(sq utils.Square) bool {
	return b.PieceAtSquare(sq) == utils.NO_PIECE
}

func (b *Board) PawnRankDir(color utils.PieceColor) int8 {
	// black pawn goes downward in rank
	var rankDir int8 = 1

	if color == utils.WHITE {
		// white pawn goes upward in rank
		rankDir = -1
	}

	return rankDir
}

func (b *Board) SquaresForPiece(p utils.Piece) []utils.Square {
	sqs := []utils.Square{}

	var rank int8
	var file int8
	for rank = 0; rank < b.NumRanks; rank++ {
		for file = 0; file < b.NumFiles; file++ {
			sq := utils.Square{file, rank}
			testp := b.PieceAtSquare(sq)
			if testp.EqualTo(p) {
				sqs = append(sqs, sq)
			}
		}
	}

	return sqs
}

func (b *Board) AttackingPieceKinds() []utils.PieceKind {
	apks := []utils.PieceKind{
		utils.Pawn,
		utils.King,
		utils.Queen,
		utils.Rook,
		utils.Bishop,
		utils.Knight,
	}

	if b.Variant == utils.VARIANT_SEIRAWAN {
		apks = append(apks, []utils.PieceKind{
			utils.Elephant,
			utils.Hawk,
		}...)
	}

	if b.Variant == utils.VARIANT_EIGHTPIECE {
		apks = append(apks, []utils.PieceKind{
			utils.Sentry,
			// TODO: lancer attacks
		}...)
	}

	return apks
}

func (b *Board) KingsAdjacent() bool {
	wk := b.WhereIsKing(utils.WHITE)

	if wk == utils.NO_SQUARE {
		return false
	}

	testk := utils.Piece{Kind: utils.King, Color: utils.BLACK}

	for _, sq := range b.AdjacentSquares(wk) {
		testp := b.PieceAtSquare(sq)
		if testp.KindColorEqualTo(testk) {
			return true
		}
	}

	return false
}

func (b *Board) MovesSortedBySan(moves []utils.Move) utils.MoveBuff {
	mb := make(utils.MoveBuff, 0)

	for _, move := range moves {
		san := b.MoveToSan(move)
		algeb := b.MoveToAlgeb(move)

		mb = append(mb, utils.MoveBuffItem{move, san, algeb})
	}

	sort.Sort(utils.MoveBuff(mb))

	return mb
}

func (b *Board) PromotionRank(color utils.PieceColor) int8 {
	if color == utils.WHITE {
		return 0
	}

	return 7
}

func (b *Board) CastlingRank(color utils.PieceColor) int8 {
	if color == utils.WHITE {
		return 7
	}

	return 0
}

func (b *Board) RookCastlingTargetSq(color utils.PieceColor, side utils.CastlingSide) utils.Square {
	rank := b.CastlingRank(color)

	var file int8 = 2

	if side == utils.KING_SIDE {
		file = 5
	}

	return utils.Square{file, rank}
}

func (b *Board) KingCastlingTargetSq(color utils.PieceColor, side utils.CastlingSide) utils.Square {
	rank := b.CastlingRank(color)

	var file int8 = 3

	if side == utils.KING_SIDE {
		file = 6
	}

	return utils.Square{file, rank}
}

func (b *Board) SquaresInDirection(origSq utils.Square, dir utils.PieceDirection) []utils.Square {
	sqs := make([]utils.Square, 0)

	currentSq := origSq.Add(dir)

	for b.HasSquare(currentSq) {
		sqs = append(sqs, currentSq)

		currentSq = currentSq.Add(dir)
	}

	return sqs
}

func (b *Board) WhereIsKing(color utils.PieceColor) utils.Square {
	var rank int8
	var file int8
	for rank = 0; rank < b.NumRanks; rank++ {
		for file = 0; file < b.NumFiles; file++ {
			p := b.Pos.Rep[rank][file]
			if (p.Kind == utils.King) && (p.Color == color) {
				return utils.Square{file, rank}
			}
		}
	}

	return utils.NO_SQUARE
}

func (b *Board) LineToString(line []utils.Move) string {
	buff := []string{}

	for _, move := range line {
		buff = append(buff, b.MoveToAlgeb(move))
	}

	return strings.Join(buff, " ")
}

func (b *Board) Pop() {
	l := len(b.MoveStack)
	if l == 0 {
		return
	}

	msi := b.MoveStack[l-1]

	b.Pos = msi.Pos

	b.MoveStack = b.MoveStack[:l-1]
}

func (b *Board) Push(move utils.Move, addSan bool) {
	san := "?"

	if addSan {
		san = b.MoveToSan(move)
	}

	oldPos := b.Pos.Clone()

	fromp := b.PieceAtSquare(move.FromSq)

	ccr := &b.Pos.CastlingRights[b.Pos.Turn]

	if fromp.Kind == utils.King {
		ccr.ClearAll()
	}

	b.SetPieceAtSquare(move.FromSq, utils.NO_PIECE)

	if move.IsPromotion() {
		b.SetPieceAtSquare(move.EffectivePromotionSquare(), move.PromotionPiece)
	}

	if move.EpCapture {
		b.SetPieceAtSquare(move.EpClearSquare, utils.NO_PIECE)
	}

	if move.Castling {
		b.SetPieceAtSquare(move.ToSq, utils.NO_PIECE)
		kctsq := b.KingCastlingTargetSq(b.Pos.Turn, move.CastlingSide)
		b.SetPieceAtSquare(kctsq, utils.Piece{Kind: utils.King, Color: b.Pos.Turn})
		rctsq := b.RookCastlingTargetSq(b.Pos.Turn, move.CastlingSide)
		b.SetPieceAtSquare(rctsq, move.RookOrigPiece)
	} else {
		b.SetPieceAtSquare(move.ToSq, fromp)
	}

	if b.IS_ATOMIC() {
		if move.IsCapture() {
			// atomic explosion
			b.SetPieceAtSquare(move.ToSq, utils.NO_PIECE)

			for _, sq := range b.AdjacentSquares(move.ToSq) {
				p := b.PieceAtSquare(sq)

				if p.Kind != utils.Pawn {
					b.SetPieceAtSquare(sq, utils.NO_PIECE)
				}
			}
		}
	}

	var side utils.CastlingSide
	for side = utils.QUEEN_SIDE; side <= utils.KING_SIDE; side++ {
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

	b.Pos.EpSquare = utils.NO_SQUARE

	if move.PawnPushByTwo {
		b.Pos.EpSquare = move.EpSquare
	}

	if move.ShouldDeleteHalfmoveClock() {
		b.Pos.HalfmoveClock = 0
	} else {
		b.Pos.HalfmoveClock++
	}

	if b.Pos.Turn == utils.WHITE {
		b.Pos.FullmoveNumber++
	}

	b.MoveStack = append(b.MoveStack, MoveStackItem{
		oldPos,
		move,
		san,
	})
}

/////////////////////////////////////////////////////////////////////
