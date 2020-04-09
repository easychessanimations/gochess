package main

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
	"time"

	chess "github.com/easychessanimations/gochess/notnil"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member function

var Nodes int
var MaxDepth int
var Start time.Time

func GetNps() (float32, float32) {
	elapsed := float32(time.Now().Sub(Start)) / float32(1e9)

	nps := float32(Nodes) / float32(elapsed)

	return nps, elapsed
}

func StopPerf() {
	nps, elapsed := GetNps()

	fmt.Printf("perf elapsed %.2f nodes %d nps %.0f", elapsed, Nodes, nps)
}

func StartPerf() {
	Nodes = 0

	Start = time.Now()
}

func PerfRecursive(pos *chess.Position, depth int) {
	Nodes++

	if depth > MaxDepth {
		return
	}

	for _, move := range pos.ValidMoves() {
		newPos := pos.Update(move)
		PerfRecursive(newPos, depth+1)
	}
}

func main() {
	pos := chess.Position{}
	pos.UnmarshalText([]byte("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"))
	MaxDepth = 5
	StartPerf()
	PerfRecursive(&pos, 0)
	StopPerf()
}

/////////////////////////////////////////////////////////////////////
