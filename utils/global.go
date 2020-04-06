package utils

import "strings"

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// global functions

func VariantKeyStringToVariantKey(vks string) VariantKey {
	vk, ok := VARIANT_KEY_STRING_TO_VARIANT_KEY[vks]

	if ok {
		return vk
	}

	return VARIANT_STANDARD
}

func VariantKeyToVariantKeyString(vk VariantKey) string {
	vks, ok := VARIANT_KEY_TO_VARIANT_KEY_STRING[vk]

	if ok {
		return vks
	}

	return "standard"
}

func StartFenForVariant(vk VariantKey) string {
	fen, ok := START_FENS[vk]

	if ok {
		return fen
	}

	return STANDARD_START_FEN
}

func NumFiles(variantKey VariantKey) int8 {
	return 8
}

func NumRanks(variantKey VariantKey) int8 {
	return 8
}

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

/////////////////////////////////////////////////////////////////////
