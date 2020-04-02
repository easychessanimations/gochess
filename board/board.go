package board

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
	"sort"
	"strings"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// about

func AboutStr() string {
	return "\n--------------\ngo chess board\n--------------\n"
}

func About() {
	fmt.Println(AboutStr())
}

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// global functions

func PieceLetterToPiece(pieceLetter string) Piece {
	pieceKind, _ := PIECE_LETTER_TO_PIECE_KIND[strings.ToLower((pieceLetter[0:1]))]

	color := WHITE
	if pieceLetter >= "a" {
		color = BLACK
	}

	dirStr := ""
	if pieceKind == Lancer {
		dirStr = pieceLetter[1:]
	}

	return Piece{pieceKind, color, DirectionStringToPieceDirection(dirStr)}
}

func DirectionStringToPieceDirection(dirStr string) PieceDirection {
	dir, ok := DIRECTION_STRING_TO_PIECE_DIRECTION[dirStr]

	if ok {
		return dir
	}

	return PieceDirection{}
}

func NumFiles(variantKey VariantKey) int8 {
	return 8
}

func NumRanks(variantKey VariantKey) int8 {
	return 8
}

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (pd PieceDirection) ToString() string {
	dirStr, ok := PIECE_DIRECTION_TO_DIRECTION_STRING[pd]

	if ok {
		return dirStr
	}

	return ""
}

func (p *Piece) ToString() string {
	letter, _ := PIECE_KIND_TO_PIECE_LETTER[p.Kind]

	if p.Color {
		letter = strings.ToUpper(letter)
	}

	return letter + p.Direction.ToString()
}

func (p *Piece) ToStringUpper() string {
	str := p.ToString()
	return strings.ToUpper(str[0:1]) + str[1:]
}

func (br *BoardRep) Init(variant VariantKey) {
	br.NumFiles = NumFiles(variant)
	br.LastFile = br.NumFiles - 1
	br.NumRanks = NumRanks(variant)
	br.LastRank = br.NumRanks - 1

	br.Rep = make(map[Square]Piece)

	var rank int8
	var file int8
	for rank = 0; rank < br.NumRanks; rank++ {
		for file = 0; file < br.NumFiles; file++ {
			br.Rep[Square{file, rank}] = Piece{}
		}
	}
}

func (br *BoardRep) HasSquare(sq Square) bool {
	_, ok := br.Rep[sq]

	return ok
}

func (br *BoardRep) SetPieceAtSquare(sq Square, p Piece) bool {
	if br.HasSquare(sq) {
		br.Rep[sq] = p

		return true
	}

	return false
}

func (b *Board) SetPieceAtSquare(sq Square, p Piece) bool {
	return b.Rep.SetPieceAtSquare(sq, p)
}

func (br *BoardRep) PieceAtSquare(sq Square) Piece {
	p, ok := br.Rep[sq]

	if ok {
		return p
	}

	return Piece{}
}

func (b *Board) PieceAtSquare(sq Square) Piece {
	return b.Rep.PieceAtSquare(sq)
}

func (p *Piece) IsEmpty() bool {
	return p.Kind == NO_PIECE
}

func (br *BoardRep) SetFromFen(fen string) {
	var file int8 = 0
	var rank int8 = 0
	for index := 0; index < len(fen); {
		chr := fen[index : index+1]
		if (chr >= "0") && (chr <= "9") {
			for cumul := chr[0] - "0"[0]; cumul > 0; cumul-- {
				br.SetPieceAtSquare(Square{file, rank}, Piece{})
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
			br.SetPieceAtSquare(Square{file, rank}, PieceLetterToPiece(pieceLetter))
			file++
		}
		index++
	}
}

func (br *BoardRep) ToString() string {
	buff := ""

	var rank int8
	var file int8
	for rank = 0; rank < br.NumRanks; rank++ {
		for file = 0; file < br.NumFiles; file++ {
			piece, _ := br.Rep[Square{file, rank}]
			buff += fmt.Sprintf("%-4s", piece.ToString())
		}
		buff += "\n"
	}

	return buff
}

func (b *Board) SetFromFen(fen string) {
	fenParts := strings.Split(fen, " ")

	b.Rep.SetFromFen((fenParts[0]))
}

func (b *Board) ToString() string {
	return b.Rep.ToString()
}

func (b *Board) Print() {
	fmt.Println(b.ToString())
}

func (b *Board) Init(variant VariantKey) {
	// set variant
	b.Variant = variant

	// initialize rep to size required by variant
	b.Rep.Init(b.Variant)
}

func (sq *Square) Add(delta Square) Square {
	return Square{sq.File + delta.File, sq.Rank + delta.Rank}
}

func (b *Board) HasSquare(sq Square) bool {
	return b.Rep.HasSquare(sq)
}

func (b *Board) SquareToFileLetter(sq Square) string {
	return string([]byte{"a"[0] + byte(sq.File)})
}

func (b *Board) SquareToRankLetter(sq Square) string {
	return string([]byte{"1"[0] + byte(b.Rep.LastRank-sq.Rank)})
}

func (b *Board) SquareToAlgeb(sq Square) string {
	return b.SquareToFileLetter(sq) + b.SquareToRankLetter(sq)
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
	if move.Capture {
		if fromPiece.Kind == Pawn {
			buff = b.SquareToFileLetter(move.FromSq)
		}
		buff += "x"
	}
	buff += toAlgeb

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

func (b *BoardRep) IsSquareEmpty(sq Square) bool {
	return b.PieceAtSquare(sq).Kind == NO_PIECE
}

func (b *Board) IsSquareEmpty(sq Square) bool {
	return b.Rep.IsSquareEmpty(sq)
}

func (b *Board) PslmsForPawnAtSquare(p Piece, sq Square) []Move {
	pslms := make([]Move, 0)

	// black pawn goes down
	var rankDir int8 = 1
	if p.Color {
		// white pawn goes up
		rankDir = -1
	}

	pushOneSq := sq.Add(Square{0, rankDir})

	if b.HasSquare(pushOneSq) {
		if b.IsSquareEmpty(pushOneSq) {
			move := Move{FromSq: sq, ToSq: pushOneSq}

			pslms = append(pslms, move)

			pushTwoSq := pushOneSq.Add(Square{0, rankDir})

			if b.HasSquare(pushTwoSq) {
				if b.IsSquareEmpty(pushTwoSq) {
					plm := Move{FromSq: sq, ToSq: pushTwoSq}

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

			if top.Color != p.Color {
				plm := Move{FromSq: sq, ToSq: captureSquare, Capture: true}

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

/////////////////////////////////////////////////////////////////////
