package butils

/////////////////////////////////////////////////////////////////////
// imports

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

// FormatDisabledMove formats disabled move as FEN field
func (pos *Position) FormatDisabledMove() string {
	if !pos.curr.HasDisabledMove {
		return "-"
	}

	return pos.curr.DisableFromSquare.String() + pos.curr.DisableToSquare.String()
}

// String returns position in FEN format
// for table format use PrettyPrint
func (pos *Position) String() string {
	s := FormatPiecePlacement(pos)
	s += " " + FormatSideToMove(pos)
	s += " " + FormatCastlingAbility(pos)
	s += " " + FormatEnpassantSquare(pos)
	s += " " + strconv.Itoa(pos.curr.HalfmoveClock)
	s += " " + strconv.Itoa(pos.fullmoveCounter)
	s += " " + pos.FormatDisabledMove()
	return s
}

// popState pops one ply
func (pos *Position) popState() {
	len := len(pos.states) - 1
	pos.states = pos.states[:len]
	pos.curr = &pos.states[len-1]
	pos.Ply--
}

// pushState adds one ply
func (pos *Position) pushState() {
	len := len(pos.states)
	pos.states = append(pos.states, pos.states[len-1])
	pos.curr = &pos.states[len]
	pos.Ply++
}

// FullmoveCounter returns the number of full moves, starts from 1
func (pos *Position) FullmoveCounter() int {
	return pos.fullmoveCounter
}

// SetFullmoveCounter sets the number of full moves
func (pos *Position) SetFullmoveCounter(n int) {
	pos.fullmoveCounter = n
}

// HalfmoveClock returns the number of halfmoves since the last capture or pawn advance
func (pos *Position) HalfmoveClock() int {
	return pos.curr.HalfmoveClock
}

// SetHalfmoveClock sets the number of halfmoves since the last capture or pawn advance
func (pos *Position) SetHalfmoveClock(n int) {
	pos.curr.HalfmoveClock = n
}

// Us returns the current player to move
// Us/Them is based on Glaurung terminology
func (pos *Position) Us() Color {
	return pos.sideToMove
}

// Them returns the player awaiting to move
func (pos *Position) Them() Color {
	return pos.sideToMove.Opposite()
}

// IsEnpassantSquare returns true if sq is the en passant square
func (pos *Position) IsEnpassantSquare(sq Square) bool {
	return sq != SquareA1 && sq == pos.EnpassantSquare()
}

// EnpassantSquare returns the en passant square
// if none, return SquareA1
// this uses the polyglot definition: if the en. passant square is
// not attacked by the enemy, then EnpassantSquare() returns SquareA1
func (pos *Position) EnpassantSquare() Square {
	return pos.curr.EnpassantSquare
}

// CastlingAbility returns kings' castling ability
func (pos *Position) CastlingAbility() Castle {
	return pos.curr.CastlingAbility
}

// LastMove returns the last move played, if any
func (pos *Position) LastMove() Move {
	return pos.curr.Move
}

// Zobrist returns the zobrist key of the position, never returns 0
func (pos *Position) Zobrist() uint64 {
	if pos.curr.Zobrist != 0 {
		return pos.curr.Zobrist
	}
	return 0x4204fa763da3abeb
}

// IsPseudoLegal returns true if m is a pseudo legal move for pos
// it returns true iff m can be executed even if own king is in check
// after the move, NullMove is not a valid move
// assumes that there exists a position for which this move is valid,
// e.g. not a rook moving diagonally or a pawn promoting on 4th rank
func (pos *Position) IsPseudoLegal(m Move) bool {
	if m == NullMove ||
		m.Color() != pos.Us() ||
		pos.Get(m.From()) != m.Piece() ||
		pos.Get(m.CaptureSquare()) != m.Capture() {
		return false
	}

	from, to := m.From(), m.To()
	all := pos.ByColor(White) | pos.ByColor(Black)

	// use base figure here instead of figure
	switch m.Figure().BaseFigure() {
	case Pawn:
		// pawn move is tested above, promotion is always correct
		if m.MoveType() == Enpassant && !pos.IsEnpassantSquare(m.To()) {
			return false
		}
		if BbPawnStartRank.Has(m.From()) && BbPawnDoubleRank.Has(m.To()) && pos.Get((m.From()+m.To())/2) != NoPiece {
			return false
		}
		return true
	case Knight: // knight jumps around
		return true
	case Bishop, Rook, Queen:
		// bbSuperAttack contains queen's mobility on an empty board
		// intersecting mobility from `from` and from `to` we get
		// the diagonal, rank or file on which the piece moved; if the
		// intersection is empty we are sure that no other piece was in the way
		if bbSuperAttack[from]&bbSuperAttack[to]&all == BbEmpty {
			return true
		}
		switch m.Figure() {
		case Bishop:
			return BishopMobility(from, all).Has(to)
		case Rook:
			return RookMobility(from, all).Has(to)
		case Queen:
			return QueenMobility(from, all).Has(to)
		}
	case King:
		if m.MoveType() == Normal {
			return bbKingAttack[from].Has(to)
		}

		// king must be castling, final square is empty
		// m.MoveType() == Castling

		if m.Color() == White && m.To() == SquareG1 {
			if pos.CastlingAbility()&WhiteOO == 0 || all.Has(SquareF1) {
				return false
			}
		}
		if m.Color() == White && m.To() == SquareC1 {
			if pos.CastlingAbility()&WhiteOOO == 0 || all.Has(SquareD1) || all.Has(SquareB1) {
				return false
			}
		}
		if m.Color() == Black && m.To() == SquareG8 {
			if pos.CastlingAbility()&BlackOO == 0 || all.Has(SquareF8) {
				return false
			}
		}
		if m.Color() == Black && m.To() == SquareC8 {
			if pos.CastlingAbility()&BlackOOO == 0 || all.Has(SquareD8) || all.Has(SquareB8) {
				return false
			}
		}
		rook, start, end := CastlingRook(m.To())
		if pos.Get(start) != rook {
			return false
		}
		them := m.Color().Opposite()
		if pos.GetAttacker(m.From(), them) != NoFigure ||
			pos.GetAttacker(end, them) != NoFigure ||
			pos.GetAttacker(m.To(), them) != NoFigure {
			return false
		}
	case Lancer:
		return LancerMobility(from, m.Figure().LancerDirection(), pos.UsBb(), pos.ThemBb()).Has(to)
	case Sentry:
		// TODO: check if sentry can capture
		return BishopMobility(from, all).Has(to)
	case Jailer:
		return JailerMobility(from, pos.UsBb(), pos.ThemBb()).Has(to)
	default:
		fmt.Println(m.Figure().BaseFigure())
		panic("unreachable")
	}

	return true
}

// verify checks the validity of the position
// mostly used for debugging purposes
func (pos *Position) Verify() error {
	if bb := pos.ByColor(White) & pos.ByColor(Black); bb != 0 {
		sq := bb.Pop()
		return fmt.Errorf("Square %v is both White and Black", sq)
	}
	// check that there is at most one king
	// catches castling issues
	for col := ColorMinValue; col <= ColorMaxValue; col++ {
		bb := pos.ByPiece(col, King)
		sq := bb.Pop()
		if bb != 0 {
			sq2 := bb.Pop()
			return fmt.Errorf("More than one King for %v at %v and %v", col, sq, sq2)
		}
	}

	// verifies that pieces have the right color
	for col := ColorMinValue; col <= ColorMaxValue; col++ {
		for bb := pos.ByColor(col); bb != 0; {
			sq := bb.Pop()
			pi := pos.Get(sq)
			if pi.Color() != col {
				return fmt.Errorf("Expected color %v, got %v", col, pi)
			}
		}
	}

	// verifies that no two pieces sit on the same cell
	for pi1 := PieceMinValue; pi1 <= PieceMaxValue; pi1++ {
		for pi2 := pi1 + 1; pi2 <= PieceMaxValue; pi2++ {
			if pos.ByPiece(pi1.Color(), pi1.Figure())&pos.ByPiece(pi2.Color(), pi2.Figure()) != 0 {
				return fmt.Errorf("%v and %v overlap", pi1, pi2)
			}
		}
	}

	// verifies that en passant square is empty
	if sq := pos.curr.EnpassantSquare; sq != SquareA1 && pos.Get(sq) != NoPiece {
		return fmt.Errorf("Expected empty en passant square %v, got %v", sq, pos.Get(sq))
	}

	// verifies that the position has two kings
	if pos.ByPiece(White, King).Count() != 1 || pos.ByPiece(Black, King).Count() != 1 {
		return fmt.Errorf("Expected one king of each color")
	}

	return nil
}

// SetCastlingAbility sets the side to move, correctly updating the Zobrist key
func (pos *Position) SetCastlingAbility(castle Castle) {
	pos.curr.Zobrist ^= zobristCastle[pos.curr.CastlingAbility]
	pos.curr.CastlingAbility = castle
	pos.curr.Zobrist ^= zobristCastle[pos.curr.CastlingAbility]
}

// SetSideToMove sets the side to move, correctly updating the Zobrist key
func (pos *Position) SetSideToMove(col Color) {
	pos.curr.Zobrist ^= zobristColor[pos.sideToMove]
	pos.sideToMove = col
	pos.curr.Zobrist ^= zobristColor[pos.sideToMove]
}

// SetEnpassantSquare sets the en passant square correctly updating the Zobrist key
func (pos *Position) SetEnpassantSquare(epsq Square) {
	if epsq != SquareA1 {
		// in polyglot the hash key for en passant is updated only if
		// an en passant capture is possible next move; in other words
		// if there is an enemy pawn next to the end square of the move
		var theirs Bitboard
		var sq Square
		if epsq.Rank() == 2 { // White
			theirs, sq = pos.ByPiece(Black, Pawn), RankFile(3, epsq.File())
		} else if epsq.Rank() == 5 { // Black
			theirs, sq = pos.ByPiece(White, Pawn), RankFile(4, epsq.File())
		} else {
			panic("bad en passant square")
		}

		if (sq.File() == 0 || !theirs.Has(sq-1)) && (sq.File() == 7 || !theirs.Has(sq+1)) {
			epsq = SquareA1
		}
	}

	pos.curr.Zobrist ^= zobristEnpassant[pos.curr.EnpassantSquare]
	pos.curr.EnpassantSquare = epsq
	pos.curr.Zobrist ^= zobristEnpassant[pos.curr.EnpassantSquare]
}

// ByColor returns the bitboard occupied by color col
func (pos *Position) ByColor(col Color) Bitboard {
	return pos.curr.ByColor[col]
}

// ByFigure returns the bitboard occupied by figure fig
func (pos *Position) ByFigure(fig Figure) Bitboard {
	return pos.curr.ByFigure[fig]
}

// ByPiece is a shortcut for ByColor(col)&ByFigure(fig)
func (pos *Position) ByPiece(col Color, fig Figure) Bitboard {
	return pos.ByColor(col) & pos.ByFigure(fig)
}

// ByPiece2 is a shortcut for ByColor(col)&(ByFigure(fig0)|ByFigure(fig1))
func (pos *Position) ByPiece2(col Color, fig0, fig1 Figure) Bitboard {
	return pos.ByColor(col) & (pos.ByFigure(fig0) | pos.ByFigure(fig1))
}

// put puts a piece on the board
// does nothing if pi is NoPiece, does not validate input
func (pos *Position) Put(sq Square, pi Piece) {
	if pi != NoPiece {
		bb := sq.Bitboard()
		pos.curr.Zobrist ^= zobristPiece[pi][sq]
		pos.curr.ByColor[pi.Color()] |= bb
		pos.curr.ByFigure[pi.Figure()] |= bb
		pos.pieces[sq] = pi
	}
}

// remove removes a piece from the table
// does nothing if pi is NoPiece, does not validate input
func (pos *Position) Remove(sq Square, pi Piece) {
	if pi != NoPiece {
		bb := ^sq.Bitboard()
		pos.curr.Zobrist ^= zobristPiece[pi][sq]
		pos.curr.ByColor[pi.Color()] &= bb
		pos.curr.ByFigure[pi.Figure()] &= bb
		pos.pieces[sq] = NoPiece
	}
}

// Get returns the piece at sq
func (pos *Position) Get(sq Square) Piece {
	return pos.pieces[sq]
}

// HasLegalMoves returns true if current side has any legal moves
// this function is very expensive
func (pos *Position) HasLegalMoves() bool {
	var moves []Move
	pos.GenerateMoves(Violent|Quiet, &moves)
	for _, m := range moves {
		pos.DoMove(m)
		checked := pos.IsChecked(pos.Them())
		pos.UndoMove()
		if !checked { // check if move didn't leave the player in check
			return true
		}
	}
	return false
}

// LegalMoves returns all legal moves from the position
func (pos *Position) LegalMoves() []Move {
	var moves []Move

	pos.GenerateMoves(Violent|Quiet, &moves)

	legalMoves := []Move{}

	for _, m := range moves {
		pos.DoMove(m)
		checked := pos.IsChecked(pos.Them())
		pos.UndoMove()
		if !checked { // check if move didn't leave the player in check
			legalMoves = append(legalMoves, m)
		}
	}

	return legalMoves
}

func (pos *Position) PerftRec(depth int, maxDepth int) {
	pos.Nodes++

	if depth >= maxDepth {
		return
	}

	lms := pos.LegalMoves()

	for _, move := range lms {
		pos.DoMove(move)
		pos.PerftRec(depth+1, maxDepth)
		pos.UndoMove()
	}
}

func (pos *Position) Perft(maxDepth int, verbose bool) int {
	pos.Nodes = 0

	lms := pos.LegalMoves()

	timeStart := time.Now()

	for _, move := range lms {
		pos.DoMove(move)
		if verbose {
			fmt.Printf("perft %s : ", move.UCI())
		}
		nodesOrig := pos.Nodes
		pos.PerftRec(0, maxDepth-1)
		if verbose {
			fmt.Printf("%d\n", pos.Nodes-nodesOrig)
		}
		pos.UndoMove()
	}

	elapsed := float32(time.Now().Sub(timeStart)) / 1e9
	nps := float32(pos.Nodes) / elapsed

	if verbose {
		fmt.Printf("total %d elapsed %.0f nps %.0f\n", pos.Nodes, elapsed, nps)
	}

	return pos.Nodes
}

// CreateLegalMoveBuff creates a move buffer for all legal moves, with SAN and UCI
// sorted by SAN
func (pos *Position) CreateLegalMoveBuff() {
	lms := pos.LegalMoves()
	pos.LegalMoveBuff = MoveBuff{}
	for _, lm := range lms {
		mbi := MoveBuffItem{
			Move:  lm,
			San:   lm.LAN(), // use LAN for meaningful initialization
			Algeb: lm.UCI(),
			Lan:   lm.LAN(),
		}

		pos.LegalMoveBuff = append(pos.LegalMoveBuff, mbi)

		sort.Sort(MoveBuffBySan(pos.LegalMoveBuff))
	}
}

// InitMoveToSan should be called before batch calls to MoveToSanBatch
func (pos *Position) InitMoveToSan() {
	pos.CreateLegalMoveBuff()
	for i, mbi := range pos.LegalMoveBuff {
		mbi.San = pos.MoveToSanBatch(mbi.Move)
		pos.LegalMoveBuff[i] = mbi
	}
	sort.Sort(MoveBuffBySan(pos.LegalMoveBuff))
}

// MoveToSanBatch returns the move in SAN notation
// provided that InitMoveToSan was called for the position
func (pos *Position) MoveToSanBatch(move Move) string {
	canditates := MoveBuff{}

	seenSquares := make(map[Square]bool)

	for _, mbi := range pos.LegalMoveBuff {
		if (mbi.Move.Piece() == move.Piece()) && (mbi.Move.To() == move.To()) {
			_, seen := seenSquares[mbi.Move.From()]
			if !seen {
				seenSquares[mbi.Move.From()] = true
				canditates = append(canditates, mbi)
			}
		}
	}

	if len(canditates) == 0 {
		// move not found among legal moves
		return "-"
	}

	sameFile := false
	sameRank := false
	files := map[int]bool{}
	ranks := map[int]bool{}

	for _, candidate := range canditates {
		file := candidate.Move.From().File()
		rank := candidate.Move.From().Rank()

		_, hasFile := files[file]
		_, hasRank := ranks[rank]

		if hasFile {
			sameFile = true
		} else {
			files[file] = true
		}

		if hasRank {
			sameRank = true
		} else {
			ranks[rank] = true
		}
	}

	// full qualifier
	qualifier := move.From().String()

	if len(canditates) == 1 {
		// no qualifier for only one move
		qualifier = ""
	} else {
		if (!sameFile) && (!sameRank) {
			// default is qualify by file
			qualifier = qualifier[0:1]
		} else if sameFile && sameRank {
			// nothing to do; qualifier is should be left full, as already initialized
		} else if sameFile {
			// has same files, needs to be qualified by rank
			qualifier = qualifier[1:2]
		} else {
			// same rank, has to be qualified by file
			qualifier = qualifier[0:1]
		}
	}

	if move.Piece().Figure() == Pawn && !pos.IsSquareEmpty(move.To()) {
		qualifier = move.From().String()[0:1]
	}

	letter := move.Piece().SanLetter()
	if move.Figure() == Pawn {
		letter = ""
	}
	capture := ""
	if move.Capture() != NoPiece {
		capture = "x"
	}
	toAlgeb := move.To().String()
	prom := ""
	if move.Promotion() != NoPiece {
		prom = "=" + move.Promotion().SanSymbol()
	}

	promSq := ""

	if move.MoveType() == SentryPush {
		promSq = "@" + move.PromotionSquare().String()
	}

	san := letter + qualifier + capture + toAlgeb + prom + promSq

	check := ""

	pos.DoMove(move)
	if pos.IsChecked(pos.Us()) {
		check = "+"
		lms := pos.LegalMoves()
		if len(lms) == 0 {
			check = "#"
		}
	}
	pos.UndoMove()

	if move.MoveType() == Castling {
		san = "O-O"

		if move.To().File() < 4 {
			san = "O-O-O"
		}
	}

	return san + check
}

// MoveToSan returns the move in SAN notation
func (pos *Position) MoveToSan(move Move) string {
	pos.InitMoveToSan()
	return pos.MoveToSanBatch(move)
}

// InsufficientMaterial returns true if the position is theoretical draw
func (pos *Position) InsufficientMaterial() bool {
	// K vs K is draw
	noKings := (pos.ByColor(White) | pos.ByColor(Black)) &^ pos.ByFigure(King)
	if noKings == 0 {
		return true
	}
	// KN vs K is theoretical draw
	if n := pos.ByFigure(Knight); noKings == n && n&(n-1) == 0 {
		return true
	}
	// KB* vs KB* is theoretical draw if all bishops are on the same square color
	if bishops := pos.ByFigure(Bishop); noKings == bishops {
		if bishops&BbWhiteSquares == bishops || bishops&BbBlackSquares == bishops {
			return true
		}
	}
	return false
}

// ThreeFoldRepetition returns whether current position was seen three times already
// returns minimum between 3 and the actual number of repetitions
func (pos *Position) ThreeFoldRepetition() int {
	c, z := 0, pos.Zobrist()
	for i := 0; i < len(pos.states) && i <= pos.curr.HalfmoveClock; i += 2 {
		j := len(pos.states) - 1 - i
		if j != 0 && pos.states[j].Move == NullMove {
			// search uses NullMove for Null Move Pruning heuristic; a couple of consecutive
			// NullMoves results in a repeated position, but it's not actually a repeat
			break
		}
		if pos.states[j].Zobrist == z {
			if c++; c == 3 {
				break
			}
		}
	}
	return c
}

// FiftyMoveRule returns true if 50 moves (on each side) were made
// without any capture of pawn move
//
// if FiftyMoveRule returns true, the position is a draw
func (pos *Position) FiftyMoveRule() bool {
	return pos.curr.HalfmoveClock >= 100
}

// WhereIsKing tells the position of the king for color
func (pos *Position) WhereIsKing(color Color) Square {
	return pos.ByPiece(color, King).AsSquare()
}

// WhereIsOurKing tells the position of our king
func (pos *Position) WhereIsOurKing() Square {
	return pos.ByPiece(pos.Us(), King).AsSquare()
}

// WhereIsTheirKing tells the position of their king
func (pos *Position) WhereTheirKing() Square {
	return pos.ByPiece(pos.Them(), King).AsSquare()
}

// IsChecked returns true if side's king is checked
func (pos *Position) IsChecked(col Color) bool {
	if pos.Us() == col && pos.curr.IsCheckedKnown {
		return pos.curr.IsChecked
	}
	kingSq := pos.WhereIsKing(col)
	isChecked := pos.GetAttacker(kingSq, col.Opposite()) != NoFigure
	if pos.Us() == col {
		pos.curr.IsCheckedKnown = true
		pos.curr.IsChecked = isChecked
	}
	return isChecked
}

// GivesCheck returns true if the opposite side is in check after m is executed
// assumes that the position is legal and opposite side is not already in check
func (pos *Position) GivesCheck(m Move) bool {
	if m.MoveType() == Castling {
		// TODO: bail out on castling because it can check via rook and king
		pos.curr.GivesCheckMove, pos.curr.GivesCheckResult = NullMove, false
		pos.DoMove(m)
		givesCheck := pos.IsChecked(pos.Us())
		pos.UndoMove()
		pos.curr.GivesCheckMove, pos.curr.GivesCheckResult = m, givesCheck
		return givesCheck
	}

	us := pos.Us()
	king := pos.ByPiece(pos.Them(), King)
	fig := m.Target().Figure()
	pos.curr.GivesCheckMove = m
	pos.curr.GivesCheckResult = true

	// test whether pawn advance gives check
	if fig == Pawn {
		bb := Forward(us, m.To().Bitboard())
		bb = East(bb) | West(bb)
		if bb&king != 0 {
			return true
		}
	}
	// test whether the knight move gives check
	// there is no such thing as discovered knight check so the check must be from this move
	if fig == Knight && KnightMobility(m.To())&king != 0 {
		return true
	}

	// fast check whether king can be attacked by a Bishop, Rook, Queen, King
	// using the moves of a Queen on an empty table
	kingSq := king.AsSquare()
	ours := pos.ByColor(us)&^pos.ByPiece2(us, Pawn, Knight)&^m.From().Bitboard() | m.To().Bitboard()
	if ours&bbSuperAttack[kingSq] == 0 {
		pos.curr.GivesCheckResult = false
		return false
	}

	// test bishop, rook, queen and king
	all := pos.ByColor(White) | pos.ByColor(Black)
	all = all&^m.From().Bitboard()&^m.CaptureSquare().Bitboard() | m.To().Bitboard()
	mob := BishopMobility(kingSq, all) &^ m.From().Bitboard()
	if mob&pos.ByPiece2(us, Bishop, Queen) != 0 ||
		mob.Has(m.To()) && (fig == Bishop || fig == Queen) {
		return true
	}
	mob = RookMobility(kingSq, all) &^ m.From().Bitboard()
	if mob&pos.ByPiece2(us, Rook, Queen) != 0 ||
		mob.Has(m.To()) && (fig == Rook || fig == Queen) {
		return true
	}
	// king checking another king is an illegal move,
	// but make the result consistent with IsChecked
	mob = KingMobility(kingSq) &^ m.From().Bitboard()
	if mob&pos.ByPiece(us, King) != 0 ||
		mob.Has(m.To()) && fig == King {
		return true
	}

	pos.curr.GivesCheckResult = false
	return false
}

// LegalMovesString lists the legal moves fromt the position numbered and sorted by SAN as string
func (pos *Position) LegalMovesString() string {
	buff := ""
	pos.InitMoveToSan()
	for i, mbi := range pos.LegalMoveBuff {
		//buff += fmt.Sprintf("%d. %s [ %s , %s ] ", i+1, mbi.San, mbi.Algeb, mbi.Lan)
		buff += fmt.Sprintf("%d. %s ", i+1, mbi.San)
	}
	return buff
}

// PrettyPrint pretty prints the current position to string
func (pos *Position) PrettyPrintString() string {
	buff := ""

	for r := 7; r >= 0; r-- {
		line := ""
		for f := 0; f < 8; f++ {
			sq := RankFile(r, f)
			if pos.IsEnpassantSquare(sq) {
				line += ","
			} else {
				line += pos.Get(sq).PrettySymbol()
			}
		}
		if r == HomeRank(pos.Us()) {
			line += " *"
		}
		buff += line + "\n"
	}

	//buff += fmt.Sprintf("zobrist = %v\n", pos.Zobrist())
	buff += fmt.Sprintf("\n%v\n", pos.String())
	buff += fmt.Sprintf("\n%v", pos.LegalMovesString())

	return buff
}

// PrettyPrint pretty prints the current position
func (pos *Position) PrettyPrint() {
	fmt.Println(pos.PrettyPrintString())
}

// UCIToMove parses a move given in UCI format
// s can be "a2a4" or "h7h8Q" for pawn promotion
func (pos *Position) UCIToMove(s string) (Move, error) {
	if len(s) < 4 {
		return NullMove, fmt.Errorf("%s is too short", s)
	}

	from, err := SquareFromString(s[0:2])
	if err != nil {
		return NullMove, err
	}
	to, err := SquareFromString(s[2:4])
	if err != nil {
		return NullMove, err
	}

	moveType := Normal
	capt := pos.Get(to)
	target := pos.Get(from)

	pi := pos.Get(from)
	if pi.Figure() == Pawn && pos.IsEnpassantSquare(to) {
		moveType = Enpassant
		capt = ColorFigure(pos.Them(), Pawn)
	}
	if pi == WhiteKing && from == SquareE1 && (to == SquareC1 || to == SquareG1) {
		moveType = Castling
	}
	if pi == BlackKing && from == SquareE8 && (to == SquareC8 || to == SquareG8) {
		moveType = Castling
	}
	if pi.Figure() == Pawn && (to.Rank() == 0 || to.Rank() == 7) {
		if len(s) != 5 {
			return NullMove, fmt.Errorf("%s doesn't have a promotion piece", s)
		}
		moveType = Promotion
		target = ColorFigure(pos.Us(), SymbolToFigure(s[4:5]))
	} else {
		if len(s) != 4 {
			return NullMove, fmt.Errorf("%s move is too long", s)
		}
	}

	move := MakeMove(moveType, from, to, target, capt, pi, NO_SQUARE, NoPiece)
	if !pos.IsPseudoLegal(move) {
		return NullMove, fmt.Errorf("%s is not a valid move", s)
	}
	return move, nil
}

// DoMove executes a legal move
func (pos *Position) DoMove(move Move) {
	pos.pushState()
	curr := pos.curr

	// update castling rights
	pi := move.Piece()
	if pi != NoPiece { // nullmove cannot change castling ability
		pos.SetCastlingAbility(curr.CastlingAbility &^ lostCastleRights[move.From()] &^ lostCastleRights[move.To()])
		if move.MoveType() == SentryPush {
			pos.SetCastlingAbility(curr.CastlingAbility &^ lostCastleRights[move.PromotionSquare()])
		}
	}
	// update fullmove counter
	if pos.Us() == Black {
		pos.fullmoveCounter++
	}
	// update halfmove clock
	curr.HalfmoveClock++
	if pi.Figure() == Pawn || move.Capture() != NoPiece {
		curr.HalfmoveClock = 0
	}
	// set Enpassant square for capturing
	if pi.Figure() == Pawn && move.From().Rank()^move.To().Rank() == 2 {
		pos.SetEnpassantSquare((move.From() + move.To()) / 2)
	} else if pos.EnpassantSquare() != SquareA1 {
		pos.SetEnpassantSquare(SquareA1)
	}
	// move rook on castling
	if move.MoveType() == Castling {
		rook, start, end := CastlingRook(move.To())
		pos.Remove(start, rook)
		pos.Put(end, rook)
	}

	// delete any former disabled move
	pos.curr.HasDisabledMove = false

	// update the pieces on the chess board
	if move.MoveType() == SentryPush {
		pos.Remove(move.From(), pi)

		pos.Remove(move.To(), move.Capture())
		pos.Put(move.To(), pi)

		promSq := move.PromotionSquare()
		pos.Remove(promSq, move.PromotionCapture())
		pos.Put(promSq, move.Target())

		// set disabled move
		pos.curr.HasDisabledMove = true
		pos.curr.DisableFromSquare = move.PromotionSquare()
		pos.curr.DisableToSquare = move.To()
	} else {
		pos.Remove(move.From(), pi)

		pos.Remove(move.CaptureSquare(), move.Capture())
		pos.Put(move.To(), move.Target())
	}

	// invert side to move
	pos.InvertSideToMove()

	curr.Move = move
	curr.IsCheckedKnown = move != NullMove && pos.curr.GivesCheckMove == move
	curr.IsChecked = curr.IsCheckedKnown && pos.curr.GivesCheckResult
	curr.GivesCheckMove = NullMove
	curr.GivesCheckResult = false

	// calculate jailed squares
	pos.calcJailedSquares()
}

// InvertSideToMove inverts the side to move
func (pos *Position) InvertSideToMove() {
	pos.SetSideToMove(pos.Them())
}

// UndoMove takes back the last move, there should be at least one move on the stack
func (pos *Position) UndoMove() {
	move := pos.LastMove()
	pos.InvertSideToMove()

	if move != NullMove {
		pos.pieces[move.From()] = move.Piece()
		pos.pieces[move.To()] = NoPiece
		pos.pieces[move.CaptureSquare()] = move.Capture()

		if move.MoveType() == SentryPush {
			pos.pieces[move.PromotionSquare()] = move.PromotionCapture()
		}
	}
	if move.MoveType() == Castling {
		rook, start, end := CastlingRook(move.To())
		pos.pieces[start] = rook
		pos.pieces[end] = NoPiece
	}

	if pos.Us() == Black {
		pos.fullmoveCounter--
	}
	pos.popState()
}

// UndoMoveSafe takes back the last move, does nothing if there is no move on the stack
func (pos *Position) UndoMoveSafe() {
	if len(pos.states) <= 1 {
		return
	}

	pos.UndoMove()
}

// JailedForColor tells squares jailed for color
func (pos *Position) JailedForColor(color Color) Bitboard {
	return pos.curr.JailedForColor[color]
}

// JailedForUs tells squares jailed for us
func (pos *Position) JailedForUs() Bitboard {
	return pos.curr.JailedForColor[pos.Us()]
}

// JailedForThem tells squares jailed for them
func (pos *Position) JailedForThem() Bitboard {
	return pos.curr.JailedForColor[pos.Them()]
}

// IsSquareJailedForColor tells whether a square is jailed for color
func (pos *Position) IsSquareJailedForColor(sq Square, color Color) bool {
	return (pos.JailedForColor(color) & sq.Bitboard()) != 0
}

// IsSquareJailedForUs tells whether a square is jailed for us
func (pos *Position) IsSquareJailedForUs(sq Square) bool {
	return (pos.JailedForColor(pos.Us()) & sq.Bitboard()) != 0
}

// IsSquareJailedForThem tells whether a square is jailed for them
func (pos *Position) IsSquareJailedForThem(sq Square) bool {
	return (pos.JailedForColor(pos.Them()) & sq.Bitboard()) != 0
}

// AppendMove appends a move to a move list
// if the move is a disabled move, it does nothing
func (pos *Position) AppendMove(move Move, moves *[]Move) {
	if (!pos.curr.HasDisabledMove) || move.Piece().Figure() == Pawn {
		// position has no disabled move, or figure is pawn which cannot be disabled
		// append and return
		*moves = append(*moves, move)
		return
	}

	disableFromSq := pos.curr.DisableFromSquare
	disableToSq := pos.curr.DisableToSquare

	// check from square
	if move.From() != disableFromSq {
		// form square is not disabled, append and return
		*moves = append(*moves, move)
		return
	}

	// check exact match
	if move.To() == disableToSq {
		// exact match, return without appending
		return
	}

	// no exact match, so any non sliding piece's move can be appended
	fig := move.Piece().Figure()

	if fig == Knight || fig == King {
		// non sliding piece with free move, append and return
		*moves = append(*moves, move)
		return
	}

	ndm, err := NormalizedDelta(move.From(), move.To())
	if err != nil {
		// move has no normalized delta
		// this should not happen with normal moves
		// TODO: consider panic here
		fmt.Println("warning, move has no normalized delta", move)
		return
	}

	ndd, err := NormalizedDelta(disableFromSq, disableToSq)
	if err != nil {
		// disabled move has no normalized delta, this should be an error
		panic(fmt.Sprintf("disabled move has no normalized delta, from %v to %v curr %v", disableFromSq, disableToSq, pos.curr))
	}

	if ndm != ndd {
		// normalized deltas differ, safe to append move
		*moves = append(*moves, move)
		return
	}

	// return without appending
	return
}

// genPawnPromotions generates pawn promotions of kind with from squares limited to limitFrom
func (pos *Position) genPawnPromotions(kind int, moves *[]Move, limitFrom Bitboard) {
	promFigures := []Figure{}
	if kind&Violent != 0 {
		promFigures = append(promFigures, PROMOTION_FIGURES_VIOLENT...)
	}
	if kind&Quiet != 0 {
		promFigures = append(promFigures, PROMOTION_FIGURES_QUIET...)
	}

	// get the pawns that can be promoted
	us, them := pos.Us(), pos.Them()
	all := pos.ByColor(White) | pos.ByColor(Black)
	ours := pos.ByPiece(us, Pawn) & limitFrom &^ pos.JailedForUs()
	theirs := pos.ByColor(them) // their pieces

	forward := Square(0)
	if us == White {
		ours &= BbRank7
		forward = RankFile(+1, 0)
	} else {
		ours &= BbRank2
		forward = RankFile(-1, 0)
	}

	for ours != 0 {
		from := ours.Pop()
		to := from + forward

		if !all.Has(to) { // advance front
			for _, p := range promFigures {
				pos.AppendMove(MakeMove(Promotion, from, to, ColorFigure(us, p), NoPiece, ColorFigure(us, Pawn), NO_SQUARE, NoPiece), moves)
			}
		}
		for fileDelta := -1; fileDelta <= 1; fileDelta += 2 {
			captFile := to.File() + fileDelta
			if captFile >= 0 && captFile <= 7 {
				captSq := RankFile(to.Rank(), captFile)
				if theirs.Has(captSq) {
					capt := pos.Get(captSq)
					for _, p := range promFigures {
						pos.AppendMove(MakeMove(Promotion, from, captSq, ColorFigure(us, p), capt, ColorFigure(us, Pawn), NO_SQUARE, NoPiece), moves)
					}
				}
			}
		}
	}
}

// genPawnAdvanceMoves generates pawns forward one or two squares moves of kind masked by mask
// with from squares limited to limitFrom
// does not generate promotions nor attacks
// disablePromotions when true disables promotions, false only when generating sentry pushes
func (pos *Position) genPawnAdvanceMoves(kind int, mask Bitboard, moves *[]Move, limitFrom Bitboard, allowPushByTwo bool, disablePromotion bool) {
	if kind&Quiet == 0 {
		return
	}

	ours := pos.ByPiece(pos.Us(), Pawn) & limitFrom &^ pos.JailedForUs()
	occu := pos.ByColor(White) | pos.ByColor(Black)
	pawn := ColorFigure(pos.Us(), Pawn)

	var forward Square
	if pos.Us() == White {
		ours = ours &^ South(occu)
		if disablePromotion {
			ours &^= BbRank7
		}
		forward = RankFile(+1, 0)
	} else {
		ours = ours &^ North(occu)
		if disablePromotion {
			ours &^= BbRank2
		}
		forward = RankFile(-1, 0)
	}

	for ours != 0 {
		from := ours.Pop()
		to := from + forward
		if mask.Has(to) {
			pos.AppendMove(MakeMove(Normal, from, to, pawn, NoPiece, pawn, NO_SQUARE, NoPiece), moves)
		}
		if allowPushByTwo {
			to += forward
			if mask.Has(to) && from.Rank() == HomeRank(pos.Us())^1 && !occu.Has(to) {
				pos.AppendMove(MakeMove(Normal, from, to, pawn, NoPiece, pawn, NO_SQUARE, NoPiece), moves)
			}
		}
	}
}

func (pos *Position) pawnCapture(to Square) (MoveType, Piece) {
	if pos.IsEnpassantSquare(to) {
		return Enpassant, ColorFigure(pos.Them(), Pawn)
	}
	return Normal, pos.Get(to)
}

// generate pawn attacks moves of kind
// does not generate promotions
// from squares limited to limitFrom
// disablePromotions when true disables promotions, false only when generating sentry pushes
func (pos *Position) genPawnAttackMoves(kind int, moves *[]Move, limitFrom Bitboard, disablePromotions bool) {
	if kind&Violent == 0 {
		return
	}

	theirs := pos.ByColor(pos.Them())
	if pos.curr.EnpassantSquare != SquareA1 {
		theirs |= pos.curr.EnpassantSquare.Bitboard()
	}

	forward := 0
	pawn := ColorFigure(pos.Us(), Pawn)
	ours := pos.ByPiece(pos.Us(), Pawn) & limitFrom &^ pos.JailedForUs()
	if pos.Us() == White {
		if disablePromotions {
			ours &^= BbRank7
		}
		theirs = South(theirs)
		forward = +1
	} else {
		if disablePromotions {
			ours &^= BbRank2
		}
		theirs = North(theirs)
		forward = -1
	}

	// left
	att := RankFile(forward, -1)
	for bbl := ours & East(theirs); bbl > 0; {
		from := bbl.Pop()
		to := from + att
		mt, capt := pos.pawnCapture(to)
		pos.AppendMove(MakeMove(mt, from, to, pawn, capt, pawn, NO_SQUARE, NoPiece), moves)
	}

	// right
	att = RankFile(forward, +1)
	for bbr := ours & West(theirs); bbr > 0; {
		from := bbr.Pop()
		to := from + att
		mt, capt := pos.pawnCapture(to)
		pos.AppendMove(MakeMove(mt, from, to, pawn, capt, pawn, NO_SQUARE, NoPiece), moves)
	}
}

func (pos *Position) genBitboardMoves(pi Piece, from Square, att Bitboard, moves *[]Move) {
	for att != 0 {
		to := att.Pop()
		pos.AppendMove(MakeMove(Normal, from, to, pi, pos.Get(to), pi, NO_SQUARE, NoPiece), moves)
	}
}

func (pos *Position) getMask(kind int) Bitboard {
	mask := Bitboard(0)
	if kind&Violent != 0 {
		// generate all attacks, promotions are handled specially
		mask |= pos.ByColor(pos.Them())
	}
	if kind&Quiet != 0 {
		// generate all non-attacks
		mask |= ^(pos.ByColor(White) | pos.ByColor(Black))
	}
	if pos.curr.IsCheckedKnown && pos.curr.IsChecked {
		// if the king is in check we can only move to block or avoid the check
		king := pos.ByPiece(pos.Us(), King).AsSquare()
		mask &= (pos.ByFigure(Knight) & bbKnightAttack[king]) | bbSuperAttack[king]
	}
	// minor promotions and castling are handled specially
	return mask
}

// genPieceMoves generate piece moves for fig with occupance mask for squares limitFrom
func (pos *Position) genPieceMoves(fig Figure, mask Bitboard, moves *[]Move, limitFrom Bitboard) {
	pi := ColorFigure(pos.Us(), fig)
	all := pos.ByColor(White) | pos.ByColor(Black)
	squares := pos.ByPiece(pos.Us(), fig) & limitFrom &^ pos.JailedForUs()
	for bb := squares; bb != 0; {
		from := bb.Pop()
		var att Bitboard
		switch fig {
		case Knight:
			att = KnightMobility(from)
		case Bishop:
			att = BishopMobility(from, all)
		case Rook:
			att = RookMobility(from, all)
		case Queen:
			att = QueenMobility(from, all)
		case King:
			att = KingMobility(from)
		case Jailer:
			att = JailerMobility(from, pos.UsBb(), pos.ThemBb())
		}
		pos.genBitboardMoves(pi, from, att&mask, moves)
	}
}

// UsBb returns the bitboard for us
func (pos *Position) UsBb() Bitboard {
	return pos.ByColor(pos.Us())
}

// ThemBb returns the bitboard for them
func (pos *Position) ThemBb() Bitboard {
	return pos.ByColor(pos.Them())
}

// genLancerMoves generates lancer moves for lancer masked by mask
// with from squares limited to limitFrom
// keepDirection should be true if the only promotion is the original piece
func (pos *Position) genLancerMoves(lancer Figure, mask Bitboard, moves *[]Move, limitFrom Bitboard, keepDirection bool) {
	ld := lancer.LancerDirection()
	pi := ColorFigure(pos.Us(), lancer)
	squares := pos.ByPiece(pos.Us(), lancer) & limitFrom &^ pos.JailedForUs()
	for bb := squares; bb != 0; {
		from := bb.Pop()
		if pos.curr.HasDisabledMove && pos.curr.DisableFromSquare == from {
			// nudged lancer
			for ldn := 0; ldn < NUM_LANCER_DIRECTIONS; ldn++ {
				att := LancerMobility(from, ldn, pos.UsBb(), pos.ThemBb()) & mask
				for att != 0 {
					to := att.Pop()
					targetLancer := MakeLancer(pos.Us(), ldn)
					if ldn == ld {
						// lancer's own direction
						if keepDirection {
							pos.AppendMove(MakeLancerMove(from, to, pi, pos.Get(to), targetLancer), moves)
						} else {
							for ldi := 0; ldi < NUM_LANCER_DIRECTIONS; ldi++ {
								pos.AppendMove(MakeLancerMove(from, to, pi, pos.Get(to), MakeLancer(pos.Us(), ldi)), moves)
							}
						}
					} else {
						// other direction, this can be only one move
						pos.AppendMove(MakeLancerMove(from, to, pi, pos.Get(to), targetLancer), moves)
					}
				}
			}
		} else {
			// regular lancer
			att := LancerMobility(from, ld, pos.UsBb(), pos.ThemBb()) & mask
			for att != 0 {
				to := att.Pop()
				if keepDirection {
					pos.AppendMove(MakeLancerMove(from, to, pi, pos.Get(to), pi), moves)
				} else {
					for ldi := 0; ldi < NUM_LANCER_DIRECTIONS; ldi++ {
						pos.AppendMove(MakeLancerMove(from, to, pi, pos.Get(to), MakeLancer(pos.Us(), ldi)), moves)
					}
				}
			}
		}
	}
}

// IsSquareEmpty tells whether a square is empty
func (pos *Position) IsSquareEmpty(sq Square) bool {
	return pos.Get(sq) == NoPiece
}

const ALLOW_SENTRY_PUSH = true

// genSentryMoves generates snetry moves for sentry masked by mask
// with from squares limited to limitFrom
func (pos *Position) genSentryMoves(mask Bitboard, moves *[]Move, limitFrom Bitboard) {
	pi := ColorFigure(pos.Us(), Sentry)
	all := pos.ByColor(White) | pos.ByColor(Black)
	squares := pos.ByPiece(pos.Us(), Sentry) & limitFrom &^ pos.JailedForUs()
	for bb := squares; bb != 0; {
		from := bb.Pop()
		att := BishopMobility(from, all) & mask
		for att != 0 {
			to := att.Pop()
			top := pos.Get(to)
			if top == NoPiece {
				pos.AppendMove(MakeMove(Normal, from, to, pi, NoPiece, pi, NO_SQUARE, NoPiece), moves)
			} else if ALLOW_SENTRY_PUSH {
				// sentry push
				// remove sentry so that pushed piece can move to its square
				pos.Remove(from, pi)
				// replace to piece with its inverse
				pos.Remove(to, top)
				pos.Put(to, top.ColorInverse())
				// save jailed squares for us
				jailedSquaresOld := pos.curr.JailedForColor[pos.Us()]
				// remove jailing for pushed piece
				pos.curr.JailedForColor[pos.Us()] &^= to.Bitboard()

				// generate possible moves of pushed piece
				var pushMoves []Move
				if top.Figure() == Sentry {
					// pushed sentry cannot push, so generate only quite moves
					pos.GenerateFigureMoves(Sentry, Quiet, &pushMoves, to.Bitboard())
				} else if top.Figure() == Pawn {
					// pushed pawn cannot move by two
					pos.genPawnAdvanceMoves(Quiet, BbFull, &pushMoves, to.Bitboard(), false, false)
					pos.genPawnAttackMoves(Violent, &pushMoves, to.Bitboard(), false)
					// cannot push pawn to promotion
				} else if top.BaseFigure() == Lancer {
					pos.genLancerMoves(top.Figure(), pos.getMask(Violent|Quiet), &pushMoves, to.Bitboard(), true)

					for ld := 0; ld < NUM_LANCER_DIRECTIONS; ld++ {
						// only generate nudges of direction different from the lancer's own direction
						// as the former one has already been generated
						if ld != top.LancerDirection() {
							nudgeTo, err := to.AddDelta(LANCER_DIRECTION_TO_DELTA[ld])

							if err == nil {
								// for nudge to adjacent square the square has to be empty
								if pos.IsSquareEmpty(nudgeTo) {
									pos.AppendMove(MakeMove(SentryPush, from, to, MakeLancer(pos.Them(), ld), top, pi, nudgeTo, NoPiece), moves)
								}
							}
						}
					}
				} else {
					pos.GenerateFigureMoves(top.Figure(), Violent|Quiet, &pushMoves, to.Bitboard())
				}

				for _, pushMove := range pushMoves {
					promCapture := pos.Get(pushMove.To())
					pos.AppendMove(MakeMove(SentryPush, from, to, top, top, pi, pushMove.To(), promCapture), moves)
				}

				// undo all removals / replacements
				pos.Remove(to, top.ColorInverse())
				pos.Put(to, top)
				pos.Put(from, pi)
				pos.curr.JailedForColor[pos.Us()] = jailedSquaresOld
			}
		}
	}
}

// genAllLancerMoves generates lancer moves for all directions masked by mask
// with from squares limited to limitFrom
func (pos *Position) genAllLancerMoves(mask Bitboard, moves *[]Move, limitFrom Bitboard) {
	for lancer := LancerMinValue; lancer <= LancerMaxValue; lancer++ {
		pos.genLancerMoves(lancer, mask, moves, limitFrom, false)
	}
}

// IsKingJailed tells whether king is jailed for color
func (pos *Position) IsKingJailed(color Color) bool {
	wk := pos.WhereIsKing(color)

	return pos.IsSquareJailedForColor(wk, color)
}

// IsOurKingJailed tells whether our king is jailed
func (pos *Position) IsOurKingJailed() bool {
	wk := pos.WhereIsKing(pos.Us())

	return pos.IsSquareJailedForColor(wk, pos.Us())
}

// IsTheirKingJailed tells whether their king is jailed
func (pos *Position) IsTheirKingJailed() bool {
	wk := pos.WhereIsKing(pos.Them())

	return pos.IsSquareJailedForColor(wk, pos.Them())
}

func (pos *Position) genKingCastles(kind int, moves *[]Move) {
	// no castling for jailed king
	if pos.IsOurKingJailed() {
		return
	}
	// skip if we only generate violent or evasion moves
	if kind&Quiet == 0 || pos.curr.IsChecked {
		return
	}

	rank := HomeRank(pos.Us())
	oo, ooo := WhiteOO, WhiteOOO
	if pos.Us() == Black {
		oo, ooo = BlackOO, BlackOOO
	}

	// castle king side
	if pos.curr.CastlingAbility&oo != 0 {
		r5 := RankFile(rank, 5)
		r6 := RankFile(rank, 6)
		if pos.Get(r5) == NoPiece && pos.Get(r6) == NoPiece {
			r4 := RankFile(rank, 4)
			if pos.GetAttacker(r4, pos.Them()) == NoFigure &&
				pos.GetAttacker(r5, pos.Them()) == NoFigure &&
				pos.GetAttacker(r6, pos.Them()) == NoFigure {
				pos.AppendMove(MakeMove(Castling, r4, r6, ColorFigure(pos.Us(), King), NoPiece, ColorFigure(pos.Us(), King), NO_SQUARE, NoPiece), moves)
			}
		}
	}

	// castle queen side
	if pos.curr.CastlingAbility&ooo != 0 {
		r3 := RankFile(rank, 3)
		r2 := RankFile(rank, 2)
		r1 := RankFile(rank, 1)
		if pos.Get(r3) == NoPiece && pos.Get(r2) == NoPiece && pos.Get(r1) == NoPiece {
			r4 := RankFile(rank, 4)
			if pos.GetAttacker(r4, pos.Them()) == NoFigure &&
				pos.GetAttacker(r3, pos.Them()) == NoFigure &&
				pos.GetAttacker(r2, pos.Them()) == NoFigure {
				pos.AppendMove(MakeMove(Castling, r4, r2, ColorFigure(pos.Us(), King), NoPiece, ColorFigure(pos.Us(), King), NO_SQUARE, NoPiece), moves)
			}
		}
	}
}

// SaveSideToMove saves side to move
// for switching turn for move generation
func (pos *Position) SaveSideToMove() {
	pos.oldSideToMove = pos.sideToMove
}

// RetrieveSideToMove retrieves side to move
// for switching turn for move generation
func (pos *Position) RetrieveSideToMove() {
	pos.SetSideToMove(pos.oldSideToMove)
}

// GetAttacker returns the smallest figure of color them that attacks sq
func (pos *Position) GetAttacker(sq Square, them Color) Figure {
	enemy := pos.ByColor(them) &^ pos.JailedForThem()
	if PawnThreats(pos, them).Has(sq) {
		return Pawn
	}
	if enemy&bbKnightAttack[sq]&pos.ByFigure(Knight) != 0 {
		return Knight
	}
	// quick test of queen's attack on an empty board
	// exclude pawns and knights because they were already tested
	enemy &^= pos.ByFigure(Pawn)
	enemy &^= pos.ByFigure(Knight)
	if enemy&bbSuperAttack[sq] == 0 {
		return NoFigure
	}
	all := pos.ByColor(White) | pos.ByColor(Black)
	bishop := BishopMobility(sq, all)
	if enemy&pos.ByFigure(Bishop)&bishop != 0 {
		return Bishop
	}
	rook := RookMobility(sq, all)
	if enemy&pos.ByFigure(Rook)&rook != 0 {
		return Rook
	}

	// save side to move for lancer and sentry checks and set it to them
	pos.SaveSideToMove()
	pos.SetSideToMove(them)

	// lancer checks, expensive
	// TODO: find less expensive way
	var moves []Move
	pos.genAllLancerMoves(pos.getMask(Violent), &moves, BbFull)
	for _, move := range moves {
		if move.To() == sq {
			// retrieve side to move
			pos.RetrieveSideToMove()
			return Lancer
		}
	}

	// sentry checks, expensive ( doubt there is a much better way )
	moves = []Move{}
	pos.genSentryMoves(pos.getMask(Violent), &moves, BbFull)
	for _, move := range moves {
		if move.PromotionSquare() == sq {
			// retrieve side to move
			pos.RetrieveSideToMove()
			return Sentry
		}
	}

	// retrieve side to move
	pos.RetrieveSideToMove()

	if enemy&pos.ByFigure(Queen)&(bishop|rook) != 0 {
		return Queen
	}
	if enemy&bbKingAttack[sq]&pos.ByFigure(King) != 0 {
		return King
	}
	return NoFigure
}

// calcJailedSquares calculates jailed squares for all colors
func (pos *Position) calcJailedSquares() {
	for color := ColorMinValue; color <= ColorMaxValue; color++ {
		pos.curr.JailedForColor[color] = BbEmpty

		squares := pos.ByPiece(color.Opposite(), Jailer)
		for bb := squares; bb != 0; {
			sq := bb.Pop()

			jailed := bbJailerAdjacent[sq]

			pos.curr.JailedForColor[color] |= jailed
		}
	}
}

// generateMoves appends to moves all moves valid from pos
// the generated moves are pseudo-legal, i.e. they can leave the king in check.
// kind is Quiet or Violent, or both
func (pos *Position) GenerateMoves(kind int, moves *[]Move) {
	mask := pos.getMask(kind)
	// order of the moves is important because the last quiet
	// moves will be reduced less.
	pos.genPawnPromotions(kind, moves, BbFull)
	pos.genPieceMoves(King, mask, moves, BbFull)
	pos.genKingCastles(kind, moves)
	pos.genKingPassMove(kind, moves)
	pos.genPieceMoves(Queen, mask, moves, BbFull)

	pos.genAllLancerMoves(mask, moves, BbFull)

	pos.genSentryMoves(mask, moves, BbFull)

	pos.genPieceMoves(Jailer, mask, moves, BbFull)

	pos.genPieceMoves(Rook, mask, moves, BbFull)
	pos.genPieceMoves(Bishop, mask, moves, BbFull)
	pos.genPieceMoves(Knight, mask, moves, BbFull)
	pos.genPawnAdvanceMoves(kind, mask, moves, BbFull, true, true)
	pos.genPawnAttackMoves(kind, moves, BbFull, true)
}

// GenerateFigureMoves generate moves of kind for a given figure
// the generated moves are pseudo-legal, i.e. they can leave the king in check.
// kind is Quiet or Violent, or both
// limitFrom limits from squares
func (pos *Position) GenerateFigureMoves(fig Figure, kind int, moves *[]Move, limitFrom Bitboard) {
	mask := pos.getMask(kind)
	switch fig.BaseFigure() {
	case Pawn:
		pos.genPawnAdvanceMoves(kind, mask, moves, limitFrom, true, true)
		pos.genPawnAttackMoves(kind, moves, limitFrom, true)
		pos.genPawnPromotions(kind, moves, limitFrom)
		return
	case Knight, Bishop, Rook, Queen, Jailer:
		pos.genPieceMoves(fig, mask, moves, limitFrom)
		return
	case King:
		pos.genPieceMoves(King, mask, moves, limitFrom)
		pos.genKingCastles(kind, moves)
		pos.genKingPassMove(kind, moves)
		return
	case Lancer:
		pos.genLancerMoves(fig, mask, moves, limitFrom, false)
		return
	case Sentry:
		pos.genSentryMoves(mask, moves, limitFrom)
		return
	}
}

func (pos *Position) genKingPassMove(kind int, moves *[]Move) {
	if kind&Quiet == 0 {
		return
	}
	if pos.IsOurKingJailed() {
		king := ColorFigure(pos.Us(), King)
		wk := pos.WhereIsOurKing()
		pos.AppendMove(MakeMove(Normal, wk, wk, king, king, king, NO_SQUARE, NoPiece), moves)
	}
}

func init() {
	//fmt.Println("position init")
	r := rand.New(rand.NewSource(5))
	f := func() uint64 { return uint64(r.Int63())<<32 ^ uint64(r.Int63()) }
	initZobristPiece(f)
	initZobristEnpassant(f)
	initZobristCastle(f)
	initZobristColor(f)
	//fmt.Println("position init done")
}

func initZobristPiece(f func() uint64) {
	for pi := PieceMinValue; pi <= PieceMaxValue; pi++ {
		for sq := SquareMinValue; sq <= SquareMaxValue; sq++ {
			zobristPiece[pi][sq] = f()
		}
	}
}

func initZobristEnpassant(f func() uint64) {
	for i := 0; i < 8; i++ {
		zobristEnpassant[SquareA3+Square(i)] = f()
		zobristEnpassant[SquareA6+Square(i)] = f()
	}
}

func initZobristCastle(f func() uint64) {
	r := [...]uint64{f(), f(), f(), f()}
	for i := CastleMinValue; i <= CastleMaxValue; i++ {
		if i&WhiteOO != 0 {
			zobristCastle[i] ^= r[0]
		}
		if i&WhiteOOO != 0 {
			zobristCastle[i] ^= r[1]
		}
		if i&BlackOO != 0 {
			zobristCastle[i] ^= r[2]
		}
		if i&BlackOOO != 0 {
			zobristCastle[i] ^= r[3]
		}
	}
}

func initZobristColor(f func() uint64) {
	zobristColor[White] = f()
}

/////////////////////////////////////////////////////////////////////
