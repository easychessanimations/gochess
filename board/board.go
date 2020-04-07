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

	"github.com/easychessanimations/gochess/utils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (b *Board) IS_ATOMIC() bool {
	return b.Variant == utils.VARIANT_ATOMIC
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

func (b *Board) SetFromRawFen(fen string) {
	var file int8 = 0
	var rank int8 = 0
	for index := 0; index < len(fen); {
		chr := fen[index : index+1]
		if (chr >= "0") && (chr <= "9") {
			for cumul := chr[0] - "0"[0]; cumul > 0; cumul-- {
				b.SetPieceAtSquare(utils.Square{file, rank}, utils.NO_PIECE)
				file++
			}
		} else if chr == "/" {
			rank++
			file = 0
		} else {
			pieceLetter := chr
			if (chr == "l") || (chr == "L") {
				index++
				dirFirst := fen[index : index+1]
				dirSecond := ""
				if (dirFirst == "n") || (dirFirst == "s") {
					index++
					dirSecond = fen[index : index+1]
					if (dirSecond != "w") && (dirSecond != "e") {
						dirSecond = ""
					}
				}
				pieceLetter = chr + dirFirst + dirSecond
			}
			b.SetPieceAtSquare(utils.Square{file, rank}, utils.PieceLetterToPiece(pieceLetter))
			file++
		}
		index++
	}
}

func (b *Board) ResetVariantFromUciOption() {
	variantUciOption := b.GetUciOptionByNameWithDefault("UCI_Variant", utils.UciOption{
		Value: "standard",
	})

	b.Variant = utils.VariantKeyStringToVariantKey(variantUciOption.Value)

	b.Reset()
}

func (b *Board) SetFromVariantUciOptionAndFen(fen string) {
	b.ResetVariantFromUciOption()

	b.SetFromFen(fen)
}

func (b *Board) SetFromFen(fen string) {
	fenParts := strings.Split(fen, " ")

	b.SetFromRawFen(fenParts[0])

	b.Pos.Turn.SetFromFen(fenParts[1])

	b.Pos.CastlingRights.SetFromFen(fenParts[2], b)

	b.Pos.EpSquare = b.SquareFromAlgeb(fenParts[3])

	hmc, _ := strconv.ParseInt(fenParts[4], 10, 32)

	b.Pos.HalfmoveClock = int(hmc)

	fmn, _ := strconv.ParseInt(fenParts[5], 10, 32)

	b.Pos.FullmoveNumber = int(fmn)
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

func (b *Board) AlgebToMove(algeb string) Move {
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
}

func (b *Board) Init(variant utils.VariantKey) {
	// set variant
	b.Variant = variant

	// initialize rep to size required by variant
	b.NumFiles = utils.NumFiles(variant)
	b.LastFile = b.NumFiles - 1
	b.NumRanks = utils.NumRanks(variant)
	b.LastRank = b.NumRanks - 1

	var rank int8
	var file int8
	for rank = 0; rank < b.NumRanks; rank++ {
		for file = 0; file < b.NumFiles; file++ {
			b.Pos.Rep[rank][file] = utils.NO_PIECE
		}
	}

	// init move stack
	b.MoveStack = make([]MoveStackItem, 0)

	// init position
	b.Pos.Init(b)
}

func (b *Board) HasSquare(sq utils.Square) bool {
	return (sq.File >= 0) && (sq.File < b.NumFiles) && (sq.Rank >= 0) && (sq.Rank < b.NumRanks)
}

func (b *Board) ReportRawFen() string {
	buff := ""
	cumul := 0

	var file int8
	var rank int8
	for rank = 0; rank < b.NumRanks; rank++ {
		for file = 0; file < b.NumFiles; file++ {
			p := b.PieceAtSquare(utils.Square{file, rank})

			if p == utils.NO_PIECE {
				cumul++
			} else {
				if cumul > 0 {
					buff += string([]byte{"0"[0] + byte(cumul)})
					cumul = 0
				}

				buff += p.ToString()
			}
		}

		if cumul > 0 {
			buff += string([]byte{"0"[0] + byte(cumul)})
			cumul = 0
		}

		if rank < (b.NumRanks - 1) {
			buff += "/"
		}
	}

	return buff
}

func (b *Board) ReportFen() string {
	buff := b.ReportRawFen()

	buff += " " + b.Pos.Turn.ToString()

	buff += " " + b.Pos.CastlingRights.ToString(b)

	buff += " " + b.SquareToAlgeb(b.Pos.EpSquare)

	buff += " " + fmt.Sprintf("%d", b.Pos.HalfmoveClock)

	buff += " " + fmt.Sprintf("%d", b.Pos.FullmoveNumber)

	return buff
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
	}

	return buff + checkStr
}

func (b *Board) PslmsForVectorPieceAtSquare(p utils.Piece, sq utils.Square) []Move {
	pslms := make([]Move, 0)

	pdesc, ok := utils.PIECE_KIND_TO_PIECE_DESCRIPTOR[p.Kind]

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

				if top != utils.NO_PIECE {
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

func (b *Board) AttacksOnSquareByPawn(sq utils.Square, color utils.PieceColor, stopAtFirst bool) []Move {
	attacks := make([]Move, 0)

	rdir := -b.PawnRankDir(color)

	var df int8
	for df = -1; df <= 1; df += 2 {
		testsq := sq.Add(utils.PieceDirection{df, rdir})

		if b.HasSquare(testsq) {
			testp := b.PieceAtSquare(testsq)

			if (testp.Kind == utils.Pawn) && (testp.Color == color) {
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

func (b *Board) AttacksOnSquareByVectorPiece(sq utils.Square, p utils.Piece, stopAtFirst bool) []Move {
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

func (b *Board) AttacksOnSquareByPiece(sq utils.Square, p utils.Piece, stopAtFirst bool) []Move {
	if p.Kind == utils.Pawn {
		return b.AttacksOnSquareByPawn(sq, p.Color, stopAtFirst)
	}

	return b.AttacksOnSquareByVectorPiece(sq, p, stopAtFirst)
}

func (b *Board) IsSquareAttackedByPiece(sq utils.Square, p utils.Piece) bool {
	attacks := b.AttacksOnSquareByPiece(sq, p, STOP_AT_FIRST)

	return len(attacks) > 0
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

	return apks
}

func (b *Board) IsSquareAttackedByColor(sq utils.Square, color utils.PieceColor) bool {
	apks := b.AttackingPieceKinds()

	for _, apk := range apks {
		if b.IsSquareAttackedByPiece(sq, utils.Piece{Kind: apk, Color: color}) {
			return true
		}
	}

	return false
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

func (b *Board) IsInCheck(color utils.PieceColor) bool {
	wk := b.WhereIsKing(color)

	if wk == utils.NO_SQUARE {
		// missing king is considered check
		return true
	}

	if b.IS_ATOMIC() {
		if b.IsExploded(color.Inverse()) {
			// no check if opponent king exploded but our king not
			return false
		}

		if b.KingsAdjacent() {
			// no check when kings adjacent
			return false
		}
	}

	return b.IsSquareAttackedByColor(wk, color.Inverse())
}

func (b *Board) PslmsForPawnAtSquare(p utils.Piece, sq utils.Square) []Move {
	pslms := make([]Move, 0)

	rankDir := b.PawnRankDir(p.Color)

	pushOneSq := sq.Add(utils.PieceDirection{0, rankDir})

	if b.HasSquare(pushOneSq) {
		if b.IsSquareEmpty(pushOneSq) {
			if pushOneSq.Rank == b.PromotionRank(p.Color) {
				promotionMoves := b.CreatePromotionMoves(
					sq,        // from
					pushOneSq, // to
					false,     // pawn capture
					true,      // push by one
					p.Color,   // color
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

			pushTwoSq := pushOneSq.Add(utils.PieceDirection{0, rankDir})

			if b.HasSquare(pushTwoSq) {
				if b.IsSquareEmpty(pushTwoSq) {
					epsq := utils.NO_SQUARE

					var df int8
					for df = -1; df <= 1; df += 2 {
						testsq := pushTwoSq.Add(utils.PieceDirection{df, 0})
						if b.HasSquare(testsq) {
							tp := b.PieceAtSquare(testsq)

							if (tp.Kind == utils.Pawn) && (tp.Color != p.Color) {
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
		captureSquare := sq.Add(utils.PieceDirection{fileDir, rankDir})

		if b.HasSquare(captureSquare) {
			top := b.PieceAtSquare(captureSquare)

			if (top != utils.NO_PIECE) && (top.Color != p.Color) {
				if pushOneSq.Rank == b.PromotionRank(p.Color) {
					promotionMoves := b.CreatePromotionMoves(
						sq,            // from
						captureSquare, // to
						true,          // pawn capture
						false,         // push by one
						p.Color,       // color
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
					EpClearSquare: captureSquare.Add(utils.PieceDirection{0, -rankDir}),
				}

				pslms = append(pslms, plm)
			}
		}
	}

	return pslms
}

func (b *Board) PslmsForPieceAtSquare(p utils.Piece, sq utils.Square) []Move {
	if p.Kind == utils.Pawn {
		return b.PslmsForPawnAtSquare(p, sq)
	}

	return b.PslmsForVectorPieceAtSquare(p, sq)
}

func (b *Board) PslmsForAllPiecesOfColor(color utils.PieceColor) []Move {
	pslms := make([]Move, 0)

	var rank int8
	var file int8
	for rank = 0; rank < b.NumRanks; rank++ {
		for file = 0; file < b.NumFiles; file++ {
			sq := utils.Square{file, rank}
			p := b.PieceAtSquare(sq)
			if (p.Color == color) && (p != utils.NO_PIECE) {
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
	b.SetFromFen(utils.StartFenForVariant(b.Variant))
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
	fromsq utils.Square,
	tosq utils.Square,
	pawnCapture bool,
	pawnPushByOne bool,
	color utils.PieceColor,
) []Move {
	promotionMoves := make([]Move, 0)

	promotionPieces, _ := utils.PROMOTION_PIECES[b.Variant]

	for _, pp := range promotionPieces {
		ppc := pp

		ppc.Color = color

		promotionMove := Move{
			FromSq:         fromsq,
			ToSq:           tosq,
			PawnCapture:    pawnCapture,
			PawnPushByOne:  pawnPushByOne,
			PromotionPiece: ppc,
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

	oldPos := b.Pos.Clone()

	fromp := b.PieceAtSquare(move.FromSq)

	ccr := &b.Pos.CastlingRights[b.Pos.Turn]

	if fromp.Kind == utils.King {
		ccr.ClearAll()
	}

	b.SetPieceAtSquare(move.FromSq, utils.NO_PIECE)

	if move.IsPromotion() {
		b.SetPieceAtSquare(move.ToSq, move.PromotionPiece)
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

func (b *Board) Pop() {
	l := len(b.MoveStack)
	if l == 0 {
		return
	}

	msi := b.MoveStack[l-1]

	b.Pos = msi.Pos

	b.MoveStack = b.MoveStack[:l-1]
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

func (b *Board) RookCastlingTargetSq(color utils.PieceColor, side CastlingSide) utils.Square {
	rank := b.CastlingRank(color)

	var file int8 = 2

	if side == KING_SIDE {
		file = 5
	}

	return utils.Square{file, rank}
}

func (b *Board) KingCastlingTargetSq(color utils.PieceColor, side CastlingSide) utils.Square {
	rank := b.CastlingRank(color)

	var file int8 = 3

	if side == KING_SIDE {
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

func (b *Board) PickLegalMovesFrom(pslms []Move, color utils.PieceColor) []Move {
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

func (b *Board) LegalMovesForAllPiecesOfColor(color utils.PieceColor) []Move {
	pslms := b.PslmsForAllPiecesOfColor(color)

	lms := b.PickLegalMovesFrom(pslms, color)

	return lms
}

func (b *Board) HasLegalMoveColor(color utils.PieceColor) bool {
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
	b.Alphas = 0
	b.Betas = 0

	b.Searching = true

	b.Start = time.Now()
}

func (b *Board) GetNps() (float32, float32) {
	elapsed := float32(time.Now().Sub(b.Start)) / float32(1e9)

	nps := float32(b.Nodes) / float32(elapsed)

	return nps, elapsed
}

func (b *Board) StopPerf() {
	nps, elapsed := b.GetNps()

	b.Log(fmt.Sprintf("perf elapsed %.2f nodes %d nps %.0f", elapsed, b.Nodes, nps))
}

func (b *Board) Perf(maxDepth int) {
	b.StartPerf()

	b.Log(fmt.Sprintf("perf up to depth %d", maxDepth))

	b.PerfRecursive(0, maxDepth)

	b.StopPerf()
}

func (b *Board) Material(color utils.PieceColor) (int, int, int) {
	material := 0
	materialBonus := 0
	mobility := 0

	var rank int8
	var file int8
	for rank = 0; rank < b.NumRanks; rank++ {
		for file = 0; file < b.NumFiles; file++ {
			p := b.PieceAtSquare(utils.Square{file, rank})
			if (p.Color == color) && (p != utils.NO_PIECE) {
				material += PIECE_VALUES[p.Kind]

				if p.Kind == utils.Pawn {
					if (rank >= 3) && (rank <= 4) {
						if (file >= 3) && (file <= 4) {
							materialBonus += CENTER_PAWN_BONUS
						}
					}
				}

				if p.Kind == utils.Knight {
					if (rank <= 1) || (rank >= 6) {
						materialBonus -= KNIGHT_ON_EDGE_DEDUCTION
					}
					if (file <= 1) || (file >= 6) {
						materialBonus -= KNIGHT_ON_EDGE_DEDUCTION
					}
					if (rank <= 0) || (rank >= 7) {
						materialBonus -= KNIGHT_CLOSE_TO_EDGE_DEDUCTION
					}
					if (file <= 0) || (file >= 7) {
						materialBonus -= KNIGHT_CLOSE_TO_EDGE_DEDUCTION
					}
				}
			}
		}
	}

	pslms := b.PslmsForAllPiecesOfColor(color)

	mobility += MOBILITY_BONUS * len(pslms)

	if b.IsExploded(color) {
		material -= (MATE_SCORE / 2)
	}

	return material, materialBonus, mobility
}

func (b *Board) MaterialBalance() int {
	materialWhite, materialBonusWhite, mobilityWhite := b.Material(utils.WHITE)
	materialBlack, materialBonusBlack, mobilityBlack := b.Material(utils.BLACK)
	return materialWhite + materialBonusWhite + mobilityWhite - (materialBlack + materialBonusBlack + mobilityBlack)
}

func (b *Board) Eval() int {
	return b.MaterialBalance() + rand.Intn(RANDOM_BONUS)
}

func (b *Board) EvalForColor(color utils.PieceColor) int {
	eval := b.Eval()

	if color == utils.WHITE {
		return eval
	}

	return -eval
}

func (b *Board) EvalForTurn() int {
	return b.EvalForColor(b.Pos.Turn)
	//return rand.Intn(RANDOM_BONUS)
}

func (b *Board) LineToString(line []Move) string {
	buff := []string{}

	for _, move := range line {
		buff = append(buff, b.MoveToAlgeb(move))
	}

	return strings.Join(buff, " ")
}

/////////////////////////////////////////////////////////////////////
