package board

import (
	"fmt"
	"strings"

	"github.com/easychessanimations/gochess/utils"
)

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

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

func (b *Board) LineToString(line []utils.Move) string {
	buff := []string{}

	for _, move := range line {
		buff = append(buff, b.MoveToAlgeb(move))
	}

	return strings.Join(buff, " ")
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

func (b *Board) LegalMovesToString() string {
	lms := b.LegalMovesForAllPieces()

	b.SortedSanMoveBuff = b.MovesSortedBySan(lms)

	buff := fmt.Sprintf("Legal moves ( %d ) ", len(lms))

	for i, mbi := range b.SortedSanMoveBuff {
		buff += fmt.Sprintf("%d. %s [ %s ] ", i+1, mbi.San, mbi.Algeb)
	}

	return buff
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

func (b *Board) Print() {
	b.Log(b.ToString())
}

/////////////////////////////////////////////////////////////////////
