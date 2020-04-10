package bengine

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (a *Accum) add(s Score) {
	a.M += s.M
	a.E += s.E
	a.Values = resize(a.Values)
	a.Values[s.I] += 1
}

func (a *Accum) addN(s Score, n int32) {
	a.M += s.M * n
	a.E += s.E * n
	a.Values = resize(a.Values)
	a.Values[s.I] += int8(n)
}

func (a *Accum) merge(o Accum) {
	a.M += o.M
	a.E += o.E
	a.Values = resize(a.Values)
	for i := range o.Values {
		a.Values[i] += o.Values[i]
	}
}

func (a *Accum) deduct(o Accum) {
	a.M -= o.M
	a.E -= o.E
	a.Values = resize(a.Values)
	for i := range o.Values {
		a.Values[i] -= o.Values[i]
	}
}

/////////////////////////////////////////////////////////////////////
