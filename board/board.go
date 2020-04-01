package board

import (
	"fmt"
	"strings"
)

func AboutStr() string {
	return "go chess board"
}

func About() {
	fmt.Println(AboutStr())
}

const STANDARD_START_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type Square struct {
	file int8
	rank int8
}

type PieceKind uint8

const (
	NO_PIECE PieceKind = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
	Hawk
	Elephant
	Sentry
	Jailer
	Lancer
)

type PieceColor bool

const WHITE PieceColor = true
const BLACK PieceColor = false

type PieceDirection Square

type Piece struct {
	Kind      PieceKind
	Color     PieceColor
	Direction PieceDirection
}

var PIECE_KIND_TO_PIECE_LETTER map[PieceKind]string = map[PieceKind]string{
	NO_PIECE: "-",
	Pawn:     "p",
	Knight:   "n",
	Bishop:   "b",
	Rook:     "r",
	Queen:    "q",
	King:     "k",
	Hawk:     "h",
	Elephant: "e",
	Sentry:   "s",
	Jailer:   "j",
	Lancer:   "l",
}

func (p *Piece) ToString() string {
	letter, _ := PIECE_KIND_TO_PIECE_LETTER[p.Kind]
	if p.Color {
		letter = strings.ToUpper(letter)
	}
	return letter
}

type BoardRep struct {
	NumFiles int8
	NumRanks int8
	Rep      map[Square]Piece
}

func (br *BoardRep) Init(numFiles int8, numRanks int8) {
	br.NumFiles = numFiles
	br.NumRanks = numRanks

	br.Rep = make(map[Square]Piece)
}

func (br *BoardRep) SquareOk(file int8, rank int8) bool {
	return ((file >= 0) &&
		(file < br.NumFiles) &&
		(rank >= 0) &&
		(rank < br.NumRanks))
}

var PIECE_LETTER_TO_PIECE_KIND map[string]PieceKind = map[string]PieceKind{
	"-": NO_PIECE,
	"p": Pawn,
	"n": Knight,
	"b": Bishop,
	"r": Rook,
	"q": Queen,
	"k": King,
	"h": Hawk,
	"e": Elephant,
	"s": Sentry,
	"j": Jailer,
	"l": Lancer,
}

func PieceLetterToPiece(pieceLetter string) Piece {
	pieceKind, _ := PIECE_LETTER_TO_PIECE_KIND[strings.ToLower((pieceLetter))]
	color := WHITE
	if pieceLetter >= "a" {
		color = BLACK
	}
	return Piece{pieceKind, color, PieceDirection{}}
}

func (br *BoardRep) SetFromFen(fen string) {
	var file int8 = 0
	var rank int8 = 0
	for index := 0; index < len(fen); {
		chr := fen[index : index+1]
		if (chr >= "0") && (chr <= "9") {
			for cumul := chr[0] - "0"[0]; cumul > 0; cumul-- {
				br.Rep[Square{file, rank}] = Piece{}
				file++
			}
		} else if chr == "/" {
			rank++
			file = 0
		} else {
			piece := PieceLetterToPiece(chr)
			br.Rep[Square{file, rank}] = piece
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
			buff += piece.ToString()
		}
		buff += "\n"
	}
	return buff
}

type VariantKey uint8

const (
	VARIANT_STANDARD VariantKey = iota
	VARIANT_ATOMIC
	VARIANT_SEIRAWAN
	VARIANT_EIGHTPIECE
)

type Board struct {
	Variant VariantKey
	Rep     BoardRep
}

func NumFiles(variantKey VariantKey) int8 {
	return 8
}

func NumRanks(variantKey VariantKey) int8 {
	return 8
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

func (b *Board) Init(variantKey VariantKey) {
	// set variant
	b.Variant = variantKey

	// initialize rep to size required by variant
	b.Rep.Init(NumFiles(b.Variant), NumRanks(b.Variant))

	b.SetFromFen(STANDARD_START_FEN)
}
