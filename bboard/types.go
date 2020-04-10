package minboard

/////////////////////////////////////////////////////////////////////
// imports

import (
	"github.com/easychessanimations/gochess/butils"
	"github.com/easychessanimations/gochess/utils"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// types

type Board struct {
	Variant                           utils.VariantKey
	Pos                               *butils.Position
	SortedSanMoveBuff                 butils.MoveBuff
	LogFunc                           func(string)
	LogAnalysisInfoFunc               func(string)
	GetUciOptionByNameWithDefaultFunc func(string, utils.UciOption) utils.UciOption
}

/////////////////////////////////////////////////////////////////////
