package board

/////////////////////////////////////////////////////////////////////
// imports

import (
	"sort"

	"github.com/easychessanimations/gochess/utils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

var DEBUG = false

/////////////////////////////////////////////////////////////////////

func (b *Board) IS_ATOMIC() bool {
	return b.Variant == utils.VARIANT_ATOMIC
}

func (b *Board) IS_EIGHTPIECE() bool {
	return b.Variant == utils.VARIANT_EIGHTPIECE
}

/////////////////////////////////////////////////////////////////////

func (b *Board) HasSquare(sq utils.Square) bool {
	return (sq.File >= 0) && (sq.File < b.NumFiles) && (sq.Rank >= 0) && (sq.Rank < b.NumRanks)
}

func (b *Board) IsSquareEmpty(sq utils.Square) bool {
	return b.PieceAtSquare(sq) == utils.NO_PIECE
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

/////////////////////////////////////////////////////////////////////

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

func (b *Board) EmptyAdjacentSquares(sq utils.Square) []utils.Square {
	asqs := b.AdjacentSquares(sq)

	easqs := []utils.Square{}

	for _, testsq := range asqs {
		if b.IsSquareEmpty(testsq) {
			easqs = append(easqs, testsq)
		}
	}

	return easqs
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

func (b *Board) IsExploded(color utils.PieceColor) bool {
	wk := b.WhereIsKing(color)

	return wk == utils.NO_SQUARE
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

func (b *Board) SquaresInDirection(origSq utils.Square, dir utils.PieceDirection) []utils.Square {
	sqs := make([]utils.Square, 0)

	currentSq := origSq.Add(dir)

	for b.HasSquare(currentSq) {
		sqs = append(sqs, currentSq)

		currentSq = currentSq.Add(dir)
	}

	return sqs
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

/////////////////////////////////////////////////////////////////////

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

func (b *Board) PawnRankDir(color utils.PieceColor) int8 {
	// black pawn goes downward in rank
	var rankDir int8 = 1

	if color == utils.WHITE {
		// white pawn goes upward in rank
		rankDir = -1
	}

	return rankDir
}

func (b *Board) PawnBaseRank(color utils.PieceColor) int8 {
	if color == utils.WHITE {
		return 6
	}

	return 1
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

/////////////////////////////////////////////////////////////////////

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

func (b *Board) GetUciOptionByNameWithDefault(name string, uciOption utils.UciOption) utils.UciOption {
	if b.GetUciOptionByNameWithDefaultFunc != nil {
		return b.GetUciOptionByNameWithDefaultFunc(name, uciOption)
	}

	return uciOption
}

/////////////////////////////////////////////////////////////////////
