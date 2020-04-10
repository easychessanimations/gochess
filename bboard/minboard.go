package minboard

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"

	"github.com/easychessanimations/gochess/butils"
	"github.com/easychessanimations/gochess/utils"
)

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

func (b *Board) Init(variant utils.VariantKey) {
	b.Reset()
}

func (b *Board) Reset() {
	b.Pos, _ = butils.PositionFromFEN(butils.FENStartPos)
}

func (b *Board) Go(depth int) {

}

func (b *Board) Stop() {

}

func (b *Board) Print() {
	b.Log(b.Pos.PrettyPrintString())
}

func (b *Board) SetFromVariantUciOptionAndFen(fen string) {

}

func (b *Board) MakeAlgebMove(algeb string, addSan bool) {

}

func (b *Board) ResetVariantFromUciOption() {
	variantUciOption := b.GetUciOptionByNameWithDefault("UCI_Variant", utils.UciOption{
		Value: "standard",
	})

	b.Variant = utils.VariantKeyStringToVariantKey(variantUciOption.Value)

	b.Reset()
}

func (b *Board) GetUciOptionByNameWithDefault(name string, uciOption utils.UciOption) utils.UciOption {
	if b.GetUciOptionByNameWithDefaultFunc != nil {
		return b.GetUciOptionByNameWithDefaultFunc(name, uciOption)
	}

	return uciOption
}

/////////////////////////////////////////////////////////////////////
