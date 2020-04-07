package board

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/easychessanimations/gochess/utils"
)

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

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

func (b *Board) Reset() {
	b.SetFromFen(utils.StartFenForVariant(b.Variant))
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

	if b.IS_EIGHTPIECE() {
		b.Pos.DisabledMove = b.AlgebToMoveRaw(fenParts[6])
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

	if b.IS_EIGHTPIECE() {
		buff += " " + b.MoveToAlgeb(b.Pos.DisabledMove)
	}

	return buff
}

/////////////////////////////////////////////////////////////////////
