package board

/////////////////////////////////////////////////////////////////////
// imports

import "fmt"

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions / map

func (br *BoardRepMap) ReportFen() string {
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
				if cumul > 0 {
					buff += string([]byte{"0"[0] + byte(cumul)})
					cumul = 0
				}

				buff += p.ToString()
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

func (br *BoardRepMap) Init(variant VariantKey) {
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

func (br *BoardRepMap) HasSquare(sq Square) bool {
	_, ok := br.Rep[sq]

	return ok
}

func (br *BoardRepMap) SetPieceAtSquare(sq Square, p Piece) bool {
	if br.HasSquare(sq) {
		br.Rep[sq] = p

		return true
	}

	return false
}

func (br *BoardRepMap) PieceAtSquare(sq Square) Piece {
	p, ok := br.Rep[sq]

	if ok {
		return p
	}

	return Piece{}
}

func (br *BoardRepMap) SetFromFen(fen string) {
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

func (br *BoardRepMap) ToString() string {
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

func (b *BoardRepMap) IsSquareEmpty(sq Square) bool {
	return b.PieceAtSquare(sq).Kind == NO_PIECE
}

func (br *BoardRepMap) SquaresInDirection(origSq Square, dir PieceDirection) []Square {
	sqs := make([]Square, 0)

	currentSq := origSq.Add(dir)

	for br.HasSquare(currentSq) {
		sqs = append(sqs, currentSq)

		currentSq = currentSq.Add(dir)
	}

	return sqs
}

func (br *BoardRepMap) WhereIsKing(color PieceColor) Square {
	for sq, p := range br.Rep {
		if (p.Kind == King) && (p.Color == color) {
			return sq
		}
	}

	return NO_SQUARE
}

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions / slice

func (br *BoardRepSlice) ReportFen() string {
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
				if cumul > 0 {
					buff += string([]byte{"0"[0] + byte(cumul)})
					cumul = 0
				}

				buff += p.ToString()
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

func (br *BoardRepSlice) Init(variant VariantKey) {
	br.NumFiles = NumFiles(variant)
	br.LastFile = br.NumFiles - 1
	br.NumRanks = NumRanks(variant)
	br.LastRank = br.NumRanks - 1

	var rank int8
	var file int8
	for rank = 0; rank < br.NumRanks; rank++ {
		for file = 0; file < br.NumFiles; file++ {
			br.Rep[rank][file] = Piece{}
		}
	}
}

func (br *BoardRepSlice) HasSquare(sq Square) bool {
	return (sq.File >= 0) && (sq.File < br.NumFiles) && (sq.Rank >= 0) && (sq.Rank < br.NumRanks)
}

func (br *BoardRepSlice) SetPieceAtSquare(sq Square, p Piece) bool {
	if br.HasSquare(sq) {
		br.Rep[sq.Rank][sq.File] = p

		return true
	}

	return false
}

func (br *BoardRepSlice) PieceAtSquare(sq Square) Piece {
	if br.HasSquare(sq) {
		return br.Rep[sq.Rank][sq.File]
	}

	return Piece{Kind: NO_PIECE}
}

func (br *BoardRepSlice) SetFromFen(fen string) {
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

func (br *BoardRepSlice) ToString() string {
	buff := ""

	var rank int8
	var file int8
	for rank = 0; rank < br.NumRanks; rank++ {
		for file = 0; file < br.NumFiles; file++ {
			piece := br.Rep[rank][file]
			buff += fmt.Sprintf("%-4s", piece.ToString())
		}
		buff += "\n"
	}

	return buff
}

func (b *BoardRepSlice) IsSquareEmpty(sq Square) bool {
	return b.PieceAtSquare(sq).Kind == NO_PIECE
}

func (br *BoardRepSlice) SquaresInDirection(origSq Square, dir PieceDirection) []Square {
	sqs := make([]Square, 0)

	currentSq := origSq.Add(dir)

	for br.HasSquare(currentSq) {
		sqs = append(sqs, currentSq)

		currentSq = currentSq.Add(dir)
	}

	return sqs
}

func (br *BoardRepSlice) WhereIsKing(color PieceColor) Square {
	var rank int8
	var file int8
	for rank = 0; rank < br.NumRanks; rank++ {
		for file = 0; file < br.NumFiles; file++ {
			p := br.Rep[rank][file]
			if (p.Kind == King) && (p.Color == color) {
				return Square{file, rank}
			}
		}
	}

	return NO_SQUARE
}

/////////////////////////////////////////////////////////////////////
