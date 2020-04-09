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

func init() {
	fmt.Println("bboard init")
	fmt.Println("square array size", butils.SquareArraySize)
	fmt.Println("square e1", butils.SquareE1)
}

func (b *Board) Init(variant utils.VariantKey) {

}

func (b *Board) Reset() {

}

func (b *Board) Go(depth int) {

}

func (b *Board) Stop() {

}

func (b *Board) Print() {

}

func (b *Board) ExecCommand(command string) {

}

func (b *Board) SetFromVariantUciOptionAndFen(fen string) {

}

func (b *Board) MakeAlgebMove(algeb string, addSan bool) {

}

func (b *Board) ResetVariantFromUciOption() {

}

/////////////////////////////////////////////////////////////////////
