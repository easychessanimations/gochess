package butils

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (c Castle) String() string {
	if c < NoCastle || c > AnyCastle {
		return fmt.Sprintf("Castle(%d)", c)
	}
	return castleToString[c]
}

/////////////////////////////////////////////////////////////////////
