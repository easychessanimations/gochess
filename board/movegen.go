package board

import "github.com/easychessanimations/gochess/utils"

/////////////////////////////////////////////////////////////////////
// imports

/////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////
// member functions

func (b *Board) IsSquareAttackedByPiece(sq utils.Square, p utils.Piece) bool {
	attacks := b.AttacksOnSquareByPiece(sq, p, STOP_AT_FIRST)

	return len(attacks) > 0
}

func (b *Board) IsInCheck(color utils.PieceColor) bool {
	wk := b.WhereIsKing(color)

	if wk == utils.NO_SQUARE {
		// missing king is considered check
		return true
	}

	if b.IS_ATOMIC() {
		if b.IsExploded(color.Inverse()) {
			// no check if opponent king exploded but our king not
			return false
		}

		if b.KingsAdjacent() {
			// no check when kings adjacent
			return false
		}
	}

	return b.IsSquareAttackedByColor(wk, color.Inverse())
}

func (b *Board) IsSquareAttackedByColor(sq utils.Square, color utils.PieceColor) bool {
	apks := b.AttackingPieceKinds()

	for _, apk := range apks {
		if b.IsSquareAttackedByPiece(sq, utils.Piece{Kind: apk, Color: color}) {
			return true
		}
	}

	return false
}

func (b *Board) AttacksOnSquareBySentry(sq utils.Square, color utils.PieceColor, stopAtFirst bool) []utils.Move {
	sentry := utils.Piece{
		Kind:  utils.Sentry,
		Color: color,
	}

	ssqs := b.SquaresForPiece(sentry)

	attacks := []utils.Move{}

	for _, ssq := range ssqs {
		splms := utils.MoveList(b.PslmsForPieceAtSquare(sentry, ssq))

		splms = splms.Filter(utils.SentryPush)

		for _, splm := range splms {
			if splm.PromotionSquare.EqualTo(sq) {
				attack := utils.Move{
					FromSq: splm.FromSq,
					ToSq:   splm.PromotionSquare,
				}

				attacks = append(attacks, attack)
			}
		}
	}

	return attacks
}

// TODO: these attacks don't work in move to san

func (b *Board) AttacksOnSquareByPiece(sq utils.Square, p utils.Piece, stopAtFirst bool) []utils.Move {
	if p.Kind == utils.Pawn {
		return b.AttacksOnSquareByPawn(sq, p.Color, stopAtFirst)
	}

	if p.Kind == utils.Sentry {
		return b.AttacksOnSquareBySentry(sq, p.Color, stopAtFirst)
	}

	return b.AttacksOnSquareByVectorPiece(sq, p, stopAtFirst)
}

func (b *Board) AttacksOnSquareByPawn(sq utils.Square, color utils.PieceColor, stopAtFirst bool) []utils.Move {
	attacks := make([]utils.Move, 0)

	rdir := -b.PawnRankDir(color)

	var df int8
	for df = -1; df <= 1; df += 2 {
		testsq := sq.Add(utils.PieceDirection{df, rdir})

		if b.HasSquare(testsq) {
			testp := b.PieceAtSquare(testsq)

			if (testp.Kind == utils.Pawn) && (testp.Color == color) {
				attacks = append(attacks, utils.Move{
					FromSq: testsq,
					ToSq:   sq,
				})

				if stopAtFirst {
					return attacks
				}
			}
		}
	}

	return attacks
}

func (b *Board) AttacksOnSquareByVectorPiece(sq utils.Square, p utils.Piece, stopAtFirst bool) []utils.Move {
	attacks := make([]utils.Move, 0)

	testp := p.ColorInverse()

	pslms := b.PslmsForVectorPieceAtSquare(testp, sq)

	for _, pslm := range pslms {
		if pslm.IsCapture() {
			testp := b.PieceAtSquare(pslm.ToSq)
			if testp.KindColorEqualTo(p) {
				attack := utils.Move{
					FromSq: pslm.ToSq,
					ToSq:   pslm.FromSq,
				}

				attacks = append(attacks, attack)

				if stopAtFirst {
					return attacks
				}
			}
		}
	}

	return attacks
}

func (b *Board) LancerMovesToSquare(lancer utils.Piece, fromSq utils.Square, toSq utils.Square, nudge bool) []utils.Move {
	lms := []utils.Move{}

	for _, ld := range utils.LANCER_DIRECTIONS {
		if (!nudge) || ld.EqualTo(lancer.Direction) {
			move := utils.Move{
				FromSq: fromSq,
				ToSq:   toSq,
				PromotionPiece: utils.Piece{
					Kind:      utils.Lancer,
					Color:     lancer.Color,
					Direction: ld,
				},
				PromotionSquare: utils.NO_SQUARE,
			}

			lms = append(lms, move)
		}
	}

	return lms
}

func (b *Board) PslmsForVectorPieceAtSquare(p utils.Piece, sq utils.Square) []utils.Move {
	pslms := make([]utils.Move, 0)

	pdesc, ok := utils.PIECE_KIND_TO_PIECE_DESCRIPTOR[p.Kind]

	if !ok {
		return pslms
	}

	currentSq := sq

	directions := pdesc.Directions

	nudge := false

	if p.Kind == utils.Lancer {
		if (b.DisabledMove == NO_MOVE) || (!b.DisabledMove.FromSq.EqualTo(sq)) {
			// lancer normally can only go in itw own direction
			directions = []utils.PieceDirection{p.Direction}
		} else {
			// nudged lancer
			nudge = true
		}
	}

	for _, dir := range directions {
		ok := true

		currentSq = sq.Add(dir)

		for ok {
			if b.HasSquare(currentSq) {
				top := b.PieceAtSquare(currentSq)

				capture := false
				add := true

				if top != utils.NO_PIECE {
					// non empty target square is capture
					capture = true

					if top.Color == p.Color {
						// cannot capture own piece
						add = false

						if pdesc.CanJumpOverOwnPiece {
							// for pieces that can jump over their own piece just skip this move
							capture = false
						}
					} else {
						if p.Kind == utils.Sentry {
							// sentry push
							// add manually
							add = false
							// no more moves for sentry
							ok = false
							if p.PushDisabled {
								// pushed sentry cannot push
							} else {
								top := b.PieceAtSquare(currentSq)

								topInv := top.ColorInverse()

								if top.Kind == utils.Sentry {
									// disable push for pushed sentry
									top.PushDisabled = true
								}

								pushes := utils.MoveList(b.PslmsForPieceAtSquare(topInv, currentSq))

								pushes = pushes.Filter(utils.NonPawnPushByTwo)

								for _, pslm := range pushes {
									move := utils.Move{
										FromSq:          sq,
										ToSq:            currentSq,
										SentryPush:      true,
										PromotionPiece:  top,
										PromotionSquare: pslm.ToSq,
									}

									pslms = append(pslms, move)
								}
							}
						}
					}
				}

				pslm := utils.Move{
					FromSq: sq,
					ToSq:   currentSq,
				}

				if !pdesc.Sliding {
					ok = false
				}

				if capture {
					ok = false

					pslm.Capture = capture

					if !pdesc.CanCapture {
						add = false
					}
				}

				if add {
					if p.Kind == utils.Lancer {
						pslms = append(pslms, b.LancerMovesToSquare(p, sq, currentSq, nudge)...)
					} else {
						pslms = append(pslms, pslm)
					}
				}
			} else {
				ok = false
			}

			currentSq = currentSq.Add(dir)
		}
	}

	return pslms
}

func (b *Board) PslmsForPawnAtSquare(p utils.Piece, sq utils.Square) []utils.Move {
	pslms := make([]utils.Move, 0)

	rankDir := b.PawnRankDir(p.Color)

	pushOneSq := sq.Add(utils.PieceDirection{0, rankDir})

	if b.HasSquare(pushOneSq) {
		if b.IsSquareEmpty(pushOneSq) {
			if pushOneSq.Rank == b.PromotionRank(p.Color) {
				promotionMoves := b.CreatePromotionMoves(
					sq,        // from
					pushOneSq, // to
					false,     // pawn capture
					true,      // push by one
					p.Color,   // color
				)

				pslms = append(pslms, promotionMoves...)
			} else {
				move := utils.Move{
					FromSq:        sq,
					ToSq:          pushOneSq,
					PawnPushByOne: true,
				}

				pslms = append(pslms, move)
			}

			pushTwoSq := pushOneSq.Add(utils.PieceDirection{0, rankDir})

			if b.HasSquare(pushTwoSq) {
				if b.IsSquareEmpty(pushTwoSq) {
					epsq := utils.NO_SQUARE

					var df int8
					for df = -1; df <= 1; df += 2 {
						testsq := pushTwoSq.Add(utils.PieceDirection{df, 0})
						if b.HasSquare(testsq) {
							tp := b.PieceAtSquare(testsq)

							if (tp.Kind == utils.Pawn) && (tp.Color != p.Color) {
								epsq = pushOneSq
							}
						}
					}

					plm := utils.Move{
						FromSq:        sq,
						ToSq:          pushTwoSq,
						PawnPushByTwo: true,
						EpSquare:      epsq,
					}

					pslms = append(pslms, plm)
				}
			}
		}
	}

	var fileDir int8
	for fileDir = -1; fileDir <= 1; fileDir += 2 {
		captureSquare := sq.Add(utils.PieceDirection{fileDir, rankDir})

		if b.HasSquare(captureSquare) {
			top := b.PieceAtSquare(captureSquare)

			if (top != utils.NO_PIECE) && (top.Color != p.Color) {
				if pushOneSq.Rank == b.PromotionRank(p.Color) {
					promotionMoves := b.CreatePromotionMoves(
						sq,            // from
						captureSquare, // to
						true,          // pawn capture
						false,         // push by one
						p.Color,       // color
					)

					pslms = append(pslms, promotionMoves...)
				} else {
					plm := utils.Move{
						FromSq:      sq,
						ToSq:        captureSquare,
						PawnCapture: true,
					}

					pslms = append(pslms, plm)
				}
			}

			if b.Pos.EpSquare == captureSquare {
				plm := utils.Move{
					FromSq:        sq,
					ToSq:          captureSquare,
					EpCapture:     true,
					EpClearSquare: captureSquare.Add(utils.PieceDirection{0, -rankDir}),
				}

				pslms = append(pslms, plm)
			}
		}
	}

	return pslms
}

func (b *Board) PslmsForPieceAtSquareInner(p utils.Piece, sq utils.Square) []utils.Move {
	if p.Kind == utils.Pawn {
		return b.PslmsForPawnAtSquare(p, sq)
	}

	return b.PslmsForVectorPieceAtSquare(p, sq)
}

func (b *Board) PslmsForPieceAtSquare(p utils.Piece, sq utils.Square) []utils.Move {
	if b.IsSquareJailedForColor(sq, p.Color) {
		// jailed pieces have no pseudo legal moves
		if p.Kind == utils.King {
			// except for king which can pass
			passMove := utils.Move{
				FromSq:  sq,
				ToSq:    sq,
				Capture: true,
			}

			return []utils.Move{passMove}
		}

		return []utils.Move{}
	}

	return b.PslmsForPieceAtSquareInner(p, sq)
}

func (b *Board) PslmsForAllPiecesOfColor(color utils.PieceColor) []utils.Move {
	pslms := make([]utils.Move, 0)

	var rank int8
	var file int8
	for rank = 0; rank < b.NumRanks; rank++ {
		for file = 0; file < b.NumFiles; file++ {
			sq := utils.Square{file, rank}
			p := b.PieceAtSquare(sq)
			if (p.Color == color) && (p != utils.NO_PIECE) {
				pslms = append(pslms, b.PslmsForPieceAtSquare(p, sq)...)
			}
		}
	}

	wk := b.WhereIsKing(color)

	for side := utils.QUEEN_SIDE; side <= utils.KING_SIDE; side++ {
		cs := b.Pos.CastlingRights[color][side]

		if cs.Free(b) {
			move := utils.Move{
				FromSq:        wk,
				ToSq:          cs.RookOrigSquare,
				Castling:      true,
				CastlingSide:  side,
				RookOrigPiece: cs.RookOrigPiece,
			}

			pslms = append(pslms, move)
		}
	}

	return pslms
}

func (b *Board) PickLegalMovesFrom(pslms []utils.Move, color utils.PieceColor) []utils.Move {
	lms := make([]utils.Move, 0)

	for _, pslm := range pslms {
		b.Push(pslm, !ADD_SAN)
		check := b.IsInCheck(color)
		b.Pop()

		if !check {
			lms = append(lms, pslm)
		}
	}

	return lms
}

func (b *Board) LegalMovesForAllPiecesOfColor(color utils.PieceColor) []utils.Move {
	pslms := b.PslmsForAllPiecesOfColor(color)

	lms := b.PickLegalMovesFrom(pslms, color)

	return lms
}

func (b *Board) HasLegalMoveColor(color utils.PieceColor) bool {
	return len(b.LegalMovesForAllPiecesOfColor(color)) > 0
}

func (b *Board) LegalMovesForAllPieces() []utils.Move {
	return b.LegalMovesForAllPiecesOfColor(b.Pos.Turn)
}

/////////////////////////////////////////////////////////////////////
