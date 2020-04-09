package butils

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// global functions

// RankFile returns a square with rank r and file f
// r and f should be between 0 and 7
func RankFile(r, f int) Square {
	return Square(r*8 + f)
}

// SquareFromString parses a square from a string
// the string has standard chess format [a-h][1-8]
func SquareFromString(s string) (Square, error) {
	if len(s) != 2 {
		return SquareA1, fmt.Errorf("invalid square %s", s)
	}

	f, r := -1, -1
	if 'a' <= s[0] && s[0] <= 'h' {
		f = int(s[0] - 'a')
	}
	if 'A' <= s[0] && s[0] <= 'H' {
		f = int(s[0] - 'A')
	}
	if '1' <= s[1] && s[1] <= '8' {
		r = int(s[1] - '1')
	}
	if f == -1 || r == -1 {
		return SquareA1, fmt.Errorf("invalid square %s", s)
	}

	return RankFile(r, f), nil
}

/////////////////////////////////////////////////////////////////////
