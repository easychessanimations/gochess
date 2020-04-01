package board

import (
	"fmt"
)

func AboutStr() string {
	return "go chess board"
}

func About() {
	fmt.Println(AboutStr())
}

type Square struct {
	file int8
	rank int8
}

type PieceKind uint8

const (
	Pawn PieceKind = iota
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

type BoardRep map[Square]Piece

type Board struct {
	Rep BoardRep
}

func (b *Board) Init() {
	b.Rep = make(map[Square]Piece)

	b.Rep[Square{0, 0}] = Piece{Pawn, WHITE, PieceDirection{}}
}
