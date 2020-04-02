package board

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
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

func (br *BoardRep) Init(variant VariantKey) {
	br.NumFiles = NumFiles(variant)
	br.NumRanks = NumRanks(variant)

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

/////////////////////////////////////////////////////////////////////
