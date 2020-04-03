package board

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
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

func (b *Board) ToString() string {
	buff := b.Rep.ToString()

	buff += "\n" + b.ReportFen() + "\n"

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
	//fromAlgeb := b.SquareToAlgeb(move.FromSq)
	toAlgeb := b.SquareToAlgeb(move.ToSq)
	fromPiece := b.PieceAtSquare(move.FromSq)
	pieceLetter := fromPiece.ToStringUpper()
	buff := pieceLetter //+ fromAlgeb
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

	return buff
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

		currentSq = sq.Add(Square(dir))

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

			currentSq = currentSq.Add(Square(dir))
		}
	}

	return pslms
}

func (b *Board) IsSquareEmpty(sq Square) bool {
	return b.Rep.IsSquareEmpty(sq)
}

func (b *Board) PslmsForPawnAtSquare(p Piece, sq Square) []Move {
	pslms := make([]Move, 0)

	// black pawn goes down
	var rankDir int8 = 1
	if p.Color == WHITE {
		// white pawn goes up
		rankDir = -1
	}

	pushOneSq := sq.Add(Square{0, rankDir})

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

			pushTwoSq := pushOneSq.Add(Square{0, rankDir})

			if b.HasSquare(pushTwoSq) {
				if b.IsSquareEmpty(pushTwoSq) {
					epsq := NO_SQUARE

					var df int8
					for df = -1; df <= 1; df += 2 {
						testsq := pushTwoSq.Add(Square{df, 0})
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
		captureSquare := sq.Add(Square{fileDir, rankDir})

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
					EpClearSquare: captureSquare.Add(Square{0, -rankDir}),
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

	for sq, p := range b.Rep.Rep {
		if (p.Color == color) && (p.Kind != NO_PIECE) {
			pslms = append(pslms, b.PslmsForPieceAtSquare(p, sq)...)
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
		mb = append(mb, MoveBuffItem{move, b.MoveToSan(move)})
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

func (b *Board) Push(move Move) {
	restoreRep := make([]SetPiece, 0)

	fromp := b.PieceAtSquare(move.FromSq)

	restoreRep = append(restoreRep, SetPiece{move.FromSq, fromp})

	top := b.PieceAtSquare(move.ToSq)

	restoreRep = append(restoreRep, SetPiece{move.ToSq, top})

	if fromp.Kind == King {
		b.Pos.CastlingRights[b.Pos.Turn].ClearAll()
	}

	b.SetPieceAtSquare(move.FromSq, Piece{})

	b.SetPieceAtSquare(move.ToSq, fromp)

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
		b.Pos.Clone(),
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

/////////////////////////////////////////////////////////////////////
