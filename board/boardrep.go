package board

/////////////////////////////////////////////////////////////////////
// imports

import "fmt"

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (br *BoardRep) ReportFen() string {
	buff := ""
	cumul := 0

	var file int8
	var rank int8
	for rank = 0; rank < br.NumRanks; rank++ {
		for file = 0; file < br.NumFiles; file++ {
			p := br.PieceAtSquare(Square{file, rank})

			if p.IsEmpty() {
				cumul++
			} else {
				buff += p.ToString()

				if cumul > 0 {
					buff += string([]byte{"0"[0] + byte(cumul)})
					cumul = 0
				}
			}
		}

		if cumul > 0 {
			buff += string([]byte{"0"[0] + byte(cumul)})
			cumul = 0
		}

		if rank < (br.NumRanks - 1) {
			buff += "/"
		}
	}

	return buff
}

func (br *BoardRep) Init(variant VariantKey) {
	br.NumFiles = NumFiles(variant)
	br.LastFile = br.NumFiles - 1
	br.NumRanks = NumRanks(variant)
	br.LastRank = br.NumRanks - 1

	br.Rep = make(map[Square]Piece)

	var rank int8
	var file int8
	for rank = 0; rank < br.NumRanks; rank++ {
		for file = 0; file < br.NumFiles; file++ {
			br.Rep[Square{file, rank}] = Piece{}
		}
	}
}

func (br *BoardRep) HasSquare(sq Square) bool {
	_, ok := br.Rep[sq]

	return ok
}

func (br *BoardRep) SetPieceAtSquare(sq Square, p Piece) bool {
	if br.HasSquare(sq) {
		br.Rep[sq] = p

		return true
	}

	return false
}

func (br *BoardRep) PieceAtSquare(sq Square) Piece {
	p, ok := br.Rep[sq]

	if ok {
		return p
	}

	return Piece{}
}

func (br *BoardRep) SetFromFen(fen string) {
	var file int8 = 0
	var rank int8 = 0
	for index := 0; index < len(fen); {
		chr := fen[index : index+1]
		if (chr >= "0") && (chr <= "9") {
			for cumul := chr[0] - "0"[0]; cumul > 0; cumul-- {
				br.SetPieceAtSquare(Square{file, rank}, Piece{})
				file++
			}
		} else if chr == "/" {
			rank++
			file = 0
		} else {
			pieceLetter := chr
			if (chr == "l") || (chr == "L") {
				index++
				dirFirst := fen[index : index+1]
				dirSecond := ""
				if (dirFirst == "n") || (dirFirst == "s") {
					index++
					dirSecond = fen[index : index+1]
					if (dirSecond != "w") && (dirSecond != "e") {
						dirSecond = ""
					}
				}
				pieceLetter = chr + dirFirst + dirSecond
			}
			br.SetPieceAtSquare(Square{file, rank}, PieceLetterToPiece(pieceLetter))
			file++
		}
		index++
	}
}

func (br *BoardRep) ToString() string {
	buff := ""

	var rank int8
	var file int8
	for rank = 0; rank < br.NumRanks; rank++ {
		for file = 0; file < br.NumFiles; file++ {
			piece, _ := br.Rep[Square{file, rank}]
			buff += fmt.Sprintf("%-4s", piece.ToString())
		}
		buff += "\n"
	}

	return buff
}

func (b *BoardRep) IsSquareEmpty(sq Square) bool {
	return b.PieceAtSquare(sq).Kind == NO_PIECE
}

func (br *BoardRep) Clone() BoardRep {
	clone := BoardRep{}

	clone.NumFiles = br.NumFiles
	clone.LastFile = br.LastFile
	clone.NumRanks = br.NumRanks
	clone.LastRank = br.LastRank

	clone.Rep = make(map[Square]Piece)

	for sq, piece := range br.Rep {
		clone.Rep[sq] = piece
	}

	return clone
}

func (br *BoardRep) SquaresInDirection(origSq Square, dir PieceDirection) []Square {
	sqs := make([]Square, 0)

	currentSq := origSq.Add(dir)

	for br.HasSquare(currentSq) {
		sqs = append(sqs, currentSq)

		currentSq = currentSq.Add(dir)
	}

	return sqs
}

func (br *BoardRep) WhereIsKing(color PieceColor) Square {
	for sq, p := range br.Rep {
		if (p.Kind == King) && (p.Color == color) {
			return sq
		}
	}

	return NO_SQUARE
}

/////////////////////////////////////////////////////////////////////
