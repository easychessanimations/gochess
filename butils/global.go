package butils

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
	"strconv"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// global functions

// ColorFigure returns a piece with col and fig
func ColorFigure(col Color, fig Figure) Piece {
	return Piece(fig<<1) + Piece(col>>1)
}

// RankFile returns a square with rank r and file f
// r and f should be between 0 and 7
func RankFile(r, f int) Square {
	return Square(r*8 + f)
}

// MakeMove constructs a move
func MakeMove(moveType MoveType, from, to Square, capture, target Piece) Move {
	piece := target
	if moveType == Promotion {
		piece = ColorFigure(target.Color(), Pawn)
	}

	return Move(from)<<0 +
		Move(to)<<8 +
		Move(moveType)<<16 +
		Move(target)<<20 +
		Move(capture)<<24 +
		Move(piece)<<28
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

// CastlingRook returns the rook moved during castling
// together with starting and stopping squares
func CastlingRook(kingEnd Square) (Piece, Square, Square) {
	// Explanation how rookStart works for king on E1.
	// if kingEnd == C1 == b010, then rookStart == A1 == b000
	// if kingEnd == G1 == b110, then rookStart == H1 == b111
	// So bit 3 will set bit 2 and bit 1.
	//
	// Explanation how rookEnd works for king on E1.
	// if kingEnd == C1 == b010, then rookEnd == D1 == b011
	// if kingEnd == G1 == b110, then rookEnd == F1 == b101
	// So bit 3 will invert bit 2. bit 1 is always set.
	piece := Piece(Rook<<1) + (1 - Piece(kingEnd>>5))
	rookStart := kingEnd&^3 | (kingEnd & 4 >> 1) | (kingEnd & 4 >> 2)
	rookEnd := kingEnd ^ (kingEnd & 4 >> 1) | 1
	return piece, rookStart, rookEnd
}

// NewPosition returns a new position representing an empty board.
func NewPosition() *Position {
	pos := &Position{
		fullmoveCounter: 1,
		states:          make([]state, 1, 4),
	}
	pos.curr = &pos.states[pos.Ply]
	return pos
}

// PositionFromFEN parses fen and returns the position.
//
// fen must contain the position using Forsyth–Edwards Notation
// http://en.wikipedia.org/wiki/Forsyth%E2%80%93Edwards_Notation
func PositionFromFEN(fen string) (*Position, error) {
	// Split fen into 6 fields.
	// Same as string.Fields() but creates much less garbage.
	// The optimization is important when a huge number of positions
	// need to be evaluated.
	f, p := [6]string{}, 0
	for i := 0; i < len(fen); {
		// Find the start and end of the token.
		for ; i < len(fen) && fen[i] == ' '; i++ {
		}
		start := i
		for ; i < len(fen) && fen[i] != ' '; i++ {
		}
		limit := i

		if start == limit {
			continue
		}
		if p >= len(f) {
			return nil, fmt.Errorf("FEN has too many fields")
		}
		f[p] = fen[start:limit]
		p++
	}
	if p != 4 && p != 6 {
		return nil, fmt.Errorf("FEN has wrong number of fields, expected 4 or 6")
	}

	// Parse each field.
	pos := NewPosition()
	if err := ParsePiecePlacement(f[0], pos); err != nil {
		return nil, err
	}
	if err := ParseSideToMove(f[1], pos); err != nil {
		return nil, err
	}
	if err := ParseCastlingAbility(f[2], pos); err != nil {
		return nil, err
	}
	if err := ParseEnpassantSquare(f[3], pos); err != nil {
		return nil, err
	}
	if p == 6 {
		var err error
		if pos.curr.HalfmoveClock, err = strconv.Atoi(f[4]); err != nil {
			return nil, err
		}
		if pos.fullmoveCounter, err = strconv.Atoi(f[5]); err != nil {
			return nil, err
		}
	} else {
		// Despite being required the last two fields of the FEN string
		// are often omitted. If the FEN is incomplete, provide default
		// values for halfmove clock and full move counter.
		pos.curr.HalfmoveClock = 0
		pos.fullmoveCounter = 1
	}

	pos.Ply = (pos.fullmoveCounter - 1) * 2
	if pos.Us() == Black {
		pos.Ply++
	}
	return pos, nil
}

// ParsePiecePlacement parse pieces from str (FEN like) into pos.
func ParsePiecePlacement(str string, pos *Position) error {
	r, f := 0, 0
	for _, p := range str {
		if p == '/' {
			if r == 7 {
				return fmt.Errorf("expected 8 ranks")
			}
			if f != 8 {
				return fmt.Errorf("expected 8 squares per rank, got %d", f)
			}
			r, f = r+1, 0
			continue
		}

		if '1' <= p && p <= '8' {
			f += int(p) - int('0')
			continue
		}
		pi := symbolToPiece[p]
		if pi == NoPiece {
			return fmt.Errorf("expected piece or number, got %s", string(p))
		}
		if f >= 8 {
			return fmt.Errorf("rank %d too long (%d cells)", 8-r, f)
		}

		// 7-r because FEN describes the table from 8th rank.
		pos.Put(RankFile(7-r, f), pi)
		f++

	}

	if f < 8 {
		return fmt.Errorf("rank %d too short (%d cells)", r+1, f)
	}
	return nil
}

// FormatPiecePlacement converts a position to FEN piece placement.
func FormatPiecePlacement(pos *Position) string {
	s := ""
	for r := 7; r >= 0; r-- {
		space := 0
		for f := 0; f < 8; f++ {
			sq := RankFile(r, f)
			pi := pos.Get(sq)
			if pi == NoPiece {
				space++
			} else {
				if space != 0 {
					s += strconv.Itoa(space)
					space = 0
				}
				s += pieceToSymbol[pi:][:1]
			}
		}

		if space != 0 {
			s += strconv.Itoa(space)
		}
		if r != 0 {
			s += "/"
		}
	}
	return s
}

// ParseEnpassantSquare parses the en passant square from str.
func ParseEnpassantSquare(str string, pos *Position) error {
	if str[:1] == "-" {
		pos.SetEnpassantSquare(SquareA1)
		return nil
	}
	sq, err := SquareFromString(str)
	if err != nil {
		return err
	}
	pos.SetEnpassantSquare(sq)
	return nil
}

// FormatEnpassantSquare converts position's castling ability to string.
func FormatEnpassantSquare(pos *Position) string {
	if pos.EnpassantSquare() != SquareA1 {
		return pos.EnpassantSquare().String()
	}
	return "-"
}

// ParseSideToMove sets side to move for pos from str.
func ParseSideToMove(str string, pos *Position) error {
	if col, ok := symbolToColor[str]; ok {
		pos.SetSideToMove(col)
		return nil
	}
	return fmt.Errorf("invalid color %s", str)
}

// FormatSideToMove returns "w" for white to play or "b" for black to play.
func FormatSideToMove(pos *Position) string {
	return colorToSymbol[pos.Us():][:1]
}

// ParseCastlingAbility sets castling ability for pos from str.
func ParseCastlingAbility(str string, pos *Position) error {
	if str == "-" {
		pos.SetCastlingAbility(NoCastle)
		return nil
	}

	ability := NoCastle
	for _, p := range str {
		info, ok := symbolToCastleInfo[p]
		if !ok {
			return fmt.Errorf("invalid castling ability %s", str)
		}
		ability |= info.Castle
		for i := 0; i < 2; i++ {
			if info.Piece[i] != pos.Get(info.Square[i]) {
				return fmt.Errorf("expected %v at %v, got %v",
					info.Piece[i], info.Square[i], pos.Get(info.Square[i]))
			}
		}
	}
	pos.SetCastlingAbility(ability)
	return nil
}

// FormatCastlingAbility returns a string specifying the castling ability
// using standard FEN format.
func FormatCastlingAbility(pos *Position) string {
	return pos.CastlingAbility().String()
}

// Pawns return the set of pawns of the given color.
func Pawns(pos *Position, us Color) Bitboard {
	return pos.ByPiece(us, Pawn)
}

// Knights return the set of knights of the given color.
func Knights(pos *Position, us Color) Bitboard {
	return pos.ByPiece(us, Knight)
}

// Bishops return the set of bishops of the given color.
func Bishops(pos *Position, us Color) Bitboard {
	return pos.ByPiece(us, Bishop)
}

// Rooks return the set of rooks of the given color.
func Rooks(pos *Position, us Color) Bitboard {
	return pos.ByPiece(us, Rook)
}

// Queens return the set of queens of the given color.
func Queens(pos *Position, us Color) Bitboard {
	return pos.ByPiece(us, Queen)
}

// Kings return the set of kings of the given color.
// Normally there is exactly on king for each side.
func Kings(pos *Position, us Color) Bitboard {
	return pos.ByPiece(us, King)
}

// PawnThreats returns the squares threatened by our pawns.
func PawnThreats(pos *Position, us Color) Bitboard {
	ours := Pawns(pos, us)
	return Forward(us, East(ours)|West(ours))
}

// BackwardPawns returns the our backward pawns.
// A backward pawn is a pawn that has no pawns behind them on its file or
// adjacent file, it's not isolated and cannot advance safely.
func BackwardPawns(pos *Position, us Color) Bitboard {
	ours := Pawns(pos, us)
	behind := ForwardFill(us, East(ours)|West(ours))
	doubled := BackwardSpan(us, ours)
	isolated := IsolatedPawns(pos, us)
	return ours & Backward(us, PawnThreats(pos, us.Opposite())) &^ behind &^ doubled &^ isolated
}

// DoubledPawns returns a bitboard with our doubled pawns.
func DoubledPawns(pos *Position, us Color) Bitboard {
	ours := Pawns(pos, us)
	return ours & Backward(us, ours)
}

// IsolatedPawns returns a bitboard with our isolated pawns.
func IsolatedPawns(pos *Position, us Color) Bitboard {
	ours := Pawns(pos, us)
	wings := East(ours) | West(ours)
	return ours &^ Fill(wings)
}

// PassedPawns returns a bitboard with our passed pawns.
func PassedPawns(pos *Position, us Color) Bitboard {
	// From white's POV: w - white pawn, b - black pawn, x - non-passed pawns.
	// ........
	// .....w..
	// .....x..
	// ..b..x..
	// .xxx.x..
	// .xxx.x..
	ours := Pawns(pos, us)
	theirs := pos.ByPiece(us.Opposite(), Pawn)
	theirs |= East(theirs) | West(theirs)
	block := BackwardSpan(us, theirs|ours)
	return ours &^ block
}

// ConnectedPawns returns a bitboad with our connected pawns.
func ConnectedPawns(pos *Position, us Color) Bitboard {
	ours := Pawns(pos, us)
	wings := East(ours) | West(ours)
	return ours & (North(wings) | wings | South(wings))
}

// RammedPawns returns pawns on ranks 2, 3 for white
// and rank 6 and 7 blocking an advanced enemy pawn.
func RammedPawns(pos *Position, us Color) Bitboard {
	var bb Bitboard
	if us == White {
		bb = BbRank2 | BbRank3
	} else if us == Black {
		bb = BbRank7 | BbRank6
	}
	return Pawns(pos, us) & Backward(us, pos.ByPiece(us.Opposite(), Pawn)) & bb
}

// Minors returns a bitboard with our knights and bishops.
func Minors(pos *Position, us Color) Bitboard {
	return pos.ByPiece2(us, Knight, Bishop)
}

// Majors returns a bitboard with our rooks and queens.
func Majors(pos *Position, us Color) Bitboard {
	return pos.ByPiece2(us, Rook, Queen)
}

// MinorsAndMajors returns a bitboard with minor and major pieces.
func MinorsAndMajors(pos *Position, col Color) Bitboard {
	return pos.ByColor(col) &^ pos.ByFigure(Pawn) &^ pos.ByFigure(King)
}

// OpenFiles returns our fully set files with no pawns.
func OpenFiles(pos *Position, us Color) Bitboard {
	pawns := pos.ByFigure(Pawn)
	return ^Fill(pawns)
}

// SemiOpenFiles returns our fully set files with enemy pawns, but no friendly pawns.
func SemiOpenFiles(pos *Position, us Color) Bitboard {
	ours := Pawns(pos, us)
	theirs := pos.ByPiece(us.Opposite(), Pawn)
	return Fill(theirs) &^ Fill(ours)
}

// KingArea returns an area around king.
func KingArea(pos *Position, us Color) Bitboard {
	bb := pos.ByPiece(us, King)
	bb = East(bb) | bb | West(bb)
	bb = North(bb) | bb | South(bb)
	return bb
}

// PawnPromotionSquare returns the propotion square of a col pawn on sq.
// Undefined behaviour if col is not White or Black.
func PawnPromotionSquare(col Color, sq Square) Square {
	if col == White {
		return sq | 0x38
	}
	if col == Black {
		return sq &^ 0x38
	}
	return sq
}

var homeRank = [ColorArraySize]int{0, 7, 0}

// HomeRank returns the rank of the king at the begining of the game.
// By construction HomeRank(col)^1 returns the pawn rank.
// Result is undefined if c is not White or Black.
func HomeRank(col Color) int {
	return homeRank[col]
}

/////////////////////////////////////////////////////////////////////
