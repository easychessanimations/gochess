package minboard

/////////////////////////////////////////////////////////////////////
// imports

import (
	"math/rand"
	"strconv"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (b *Board) ExecCommand(command string) bool {
	if command == "d" {
		b.Pos.UndoMoveSafe()

		b.Print()

		return true
	}

	b.SortedSanMoveBuff = b.Pos.SortedLegalMoves()

	i, err := strconv.ParseInt(command, 10, 32)

	if err == nil {
		move := b.SortedSanMoveBuff[i-1].Move

		b.Pos.DoMove(move)

		b.Print()

		return true
	}

	move := b.SortedSanMoveBuff[rand.Intn(len(b.SortedSanMoveBuff))].Move

	b.Pos.DoMove(move)

	b.Print()

	return true
}

/////////////////////////////////////////////////////////////////////
