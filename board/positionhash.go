package board

import "github.com/easychessanimations/gochess/utils"

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// about

func (ph *PositionHash) Init() {
	ph.PositionEntries = make(map[Pos]PositionEntry)
}

func (pe *PositionEntry) Init() {
	pe.MoveEntries = make(map[utils.Move]MoveEntry)
}

func (ph *PositionHash) GetPositionEntry(pos Pos) PositionEntry {
	pe, ok := ph.PositionEntries[pos]

	if ok {
		return pe
	}

	pe = PositionEntry{}

	pe.Init()

	ph.PositionEntries[pos] = pe

	return pe
}

func (pe *PositionEntry) GetMoveEntry(move utils.Move) MoveEntry {
	me, ok := pe.MoveEntries[move]

	if ok {
		return me
	}

	me = MoveEntry{}

	pe.MoveEntries[move] = me

	return me
}

func (pe *PositionEntry) SetMoveEntry(move utils.Move, me MoveEntry) {
	pe.MoveEntries[move] = me
}

func (pe *PositionEntry) DecreaseEntries(by int) {
	for move, me := range pe.MoveEntries {
		me.Eval -= by

		pe.MoveEntries[move] = me
	}
}

/////////////////////////////////////////////////////////////////////
