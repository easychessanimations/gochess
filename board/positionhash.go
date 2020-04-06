package board

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// about

func (ph *PositionHash) Init() {
	ph.PositionEntries = make(map[Pos]PositionEntry)
}

func (pe *PositionEntry) Init() {
	pe.MoveEntries = make(map[Move]MoveEntry)
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

func (pe *PositionEntry) GetMoveEntry(move Move) MoveEntry {
	me, ok := pe.MoveEntries[move]

	if ok {
		return me
	}

	me = MoveEntry{}

	pe.MoveEntries[move] = me

	return me
}

func (pe *PositionEntry) SetMoveEntry(move Move, me MoveEntry) {
	pe.MoveEntries[move] = me
}

/////////////////////////////////////////////////////////////////////
