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

func (ph *PositionEntry) GetMoveEntry(move Move) MoveEntry {
	me, ok := ph.MoveEntries[move]

	if ok {
		return me
	}

	me = MoveEntry{}

	ph.MoveEntries[move] = me

	return me
}

/////////////////////////////////////////////////////////////////////
