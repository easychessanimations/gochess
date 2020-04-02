package board

/////////////////////////////////////////////////////////////////////
// imports

import (
	"strings"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// global functions

func PieceLetterToPiece(pieceLetter string) Piece {
	pieceKind, _ := PIECE_LETTER_TO_PIECE_KIND[strings.ToLower((pieceLetter[0:1]))]

	color := WHITE
	if pieceLetter >= "a" {
		color = BLACK
	}

	dirStr := ""
	if pieceKind == Lancer {
		dirStr = pieceLetter[1:]
	}

	return Piece{pieceKind, color, DirectionStringToPieceDirection(dirStr)}
}

func DirectionStringToPieceDirection(dirStr string) PieceDirection {
	dir, ok := DIRECTION_STRING_TO_PIECE_DIRECTION[dirStr]

	if ok {
		return dir
	}

	return PieceDirection{}
}

func NumFiles(variantKey VariantKey) int8 {
	return 8
}

func NumRanks(variantKey VariantKey) int8 {
	return 8
}

/////////////////////////////////////////////////////////////////////
