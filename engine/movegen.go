package engine

func (board *Board) GenerateMoves(ml *MoveList) {
	side := board.State.SideToMove
	var ownPieces, oppPieces uint64

	for piece := P; piece <= k; piece++ {
		bitboard := board.Bitboards[piece]

		for bitboard != 0 {
			fromIdx := getIndexOfLS1B(bitboard)
			from := Square(fromIdx)
			switch piece {
			case P:
				// println(from / 8, BitboardSquares[from])
				// if (from / 8) == 3 {
				// 	PrintBitboard(oppPieces)
				// 	PrintBitboard(pawnAttacks[board.State.SideToMove][from] & oppPieces)
				// }
				ownPieces = board.OccupancyBitboards[white]
				oppPieces = board.OccupancyBitboards[black]
				board.generatePawnMoves(from, -8, ml, white, ownPieces, oppPieces, 1, 6)
			case p:
				// println(from / 8, BitboardSquares[from])
				// if (from / 8) == 5 {
				// 	PrintBitboard(oppPieces)
				// 	PrintBitboard(pawnAttacks[black][from] & oppPieces)
				// 	PrintBitboard(pawnAttacks[black][from] & ownPieces) // THIS
				// }
				ownPieces = board.OccupancyBitboards[black]
				oppPieces = board.OccupancyBitboards[white]
				board.generatePawnMoves(from, 8, ml, black, ownPieces, oppPieces, 6, 1)
				// UNREVISED !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			case N:
				ownPieces = board.OccupancyBitboards[white]
				board.generateKnightMoves(from, ml, ownPieces, white)
			case n:
				// PrintBitboard(maskKnightAttacks(b6))
				// PrintBitboard(knightAttacks[b6])
				ownPieces = board.OccupancyBitboards[black]
				board.generateKnightMoves(from, ml, ownPieces, black)
			case B:
				ownPieces = board.OccupancyBitboards[white]
				attacks := bishopAttacks(from, board.OccupancyBitboards[both]) &^ ownPieces
				ml.PushSetwiseFlag(from, attacks, func(to Square) MoveFlag {
					if board.OccupiedByOpp(to, side) {
						return Capture
					}
					return Normal
				})
			case b:
				ownPieces = board.OccupancyBitboards[black]
				attacks := bishopAttacks(from, board.OccupancyBitboards[both]) &^ ownPieces
				ml.PushSetwiseFlag(from, attacks, func(to Square) MoveFlag {
					if board.OccupiedByOpp(to, side) {
						return Capture
					}
					return Normal
				})
			case R:
				ownPieces = board.OccupancyBitboards[white]
				attacks := rookAttacks(from, board.OccupancyBitboards[both]) &^ ownPieces
				ml.PushSetwiseFlag(from, attacks, func(to Square) MoveFlag {
					if board.OccupiedByOpp(to, side) {
						return Capture
					}
					return Normal
				})
			case r:
				ownPieces = board.OccupancyBitboards[black]
				attacks := rookAttacks(from, board.OccupancyBitboards[both]) &^ ownPieces
				ml.PushSetwiseFlag(from, attacks, func(to Square) MoveFlag {
					if board.OccupiedByOpp(to, side) {
						return Capture
					}
					return Normal
				})
			case Q:
				ownPieces = board.OccupancyBitboards[white]
				attacks := GetQueenAttacks(from, board.OccupancyBitboards[both]) &^ ownPieces
				ml.PushSetwiseFlag(from, attacks, func(to Square) MoveFlag {
					if board.OccupiedByOpp(to, side) {
						return Capture
					}
					return Normal
				})
			case q:
				ownPieces = board.OccupancyBitboards[black]
				attacks := GetQueenAttacks(from, board.OccupancyBitboards[both]) &^ ownPieces
				ml.PushSetwiseFlag(from, attacks, func(to Square) MoveFlag {
					if board.OccupiedByOpp(to, side) {
						return Capture
					}
					return Normal
				})
			case K:
				ownPieces = board.OccupancyBitboards[white]
				attacks := kingAttacks[from] &^ ownPieces
				ml.PushSetwiseFlag(from, attacks, func(to Square) MoveFlag {
					if board.OccupiedByOpp(to, side) {
						return Capture
					}
					return Normal
				})
			case k:
				ownPieces = board.OccupancyBitboards[black]
				attacks := kingAttacks[from] &^ ownPieces
				ml.PushSetwiseFlag(from, attacks, func(to Square) MoveFlag {
					if board.OccupiedByOpp(to, side) {
						return Capture
					}
					return Normal
				})
			}

			bitboard = popBit(bitboard, from)
		}

	}
	board.generateCastlingMoves(ml, side)

}

func (board *Board) generateCastlingMoves(ml *MoveList, side int) {
	castle := board.State.CastlingRights
	eSquareQueenside := b8
	if side == white {
		eSquareQueenside = b1
	}
	if (castle & whiteKingside) != 0 {
		if !board.Occupied(f1) && !board.Occupied(g1) {
			if !board.isSquareAttacked(e1, black) && !board.isSquareAttacked(f1, black) {
				ml.Add(e1, g1, KingCastle)
			}
		}
	}

	if (castle & blackKingside) != 0 {
		println(castle&blackKingside, 67998, board.isSquareAttacked(f8, white), board.isSquareAttacked(e8, white))
		if !board.Occupied(f8) && !board.Occupied(g8) {
			if !board.isSquareAttacked(e8, white) && !board.isSquareAttacked(f8, white) {
				ml.Add(e8, g8, KingCastle)
			}
		}
	}

	if (castle & whiteQueenside) != 0 {
		if !board.Occupied(eSquareQueenside+2) && !board.Occupied(eSquareQueenside+1) && !board.Occupied(eSquareQueenside) {
			if !board.isSquareAttacked(eSquareQueenside+3, black) && !board.isSquareAttacked(eSquareQueenside+2, black) {
				ml.Add(eSquareQueenside+3, eSquareQueenside+1, QueenCastle)
			}
		}
	}

	if (castle & blackQueenside) != 0 {
		if !board.Occupied(d8) && !board.Occupied(c8) && !board.Occupied(b8) {
			if !board.isSquareAttacked(e8, white) && !board.isSquareAttacked(d8, white) {
				ml.Add(e8, c8, QueenCastle)
			}
		}
	}
}

// Generate both quiet pawn moves and captures and appends to
// movelist accordingly.
//
//   - 'from' is thesquare the piece is currently on.
//
//   - 'direction' is to help with white and black sides. e.g. yt pawns move forward one or teo places, and on our board, we could interpret that as moving 'from' -8 or -16 accordingly.
//
//   - 'piece' is the piece type to be moved.\
//
//   - 'ml' is the generalized move Move List
//
//   - 'side' is side to move, which upon introspection feels unnecessary rn but oh well-
//
//   - 'own' & 'opp' pieces are pretty self explnatory (from occupancy bitboards obv)
//
//   - 'promotionRank' is the rank from which promotion is considered: 6 (from 0) for yt, and 1 for black.
//
//   - the rest vars are to help with dealing with white and black sides.
//
//   - i intend to make this function a general solution and not dependent on whatever side, so these vars help setup constraints as regards moves, e.g. white promotion rank is from a7->h7 anything else is useless, pretty much.
func (board *Board) generatePawnMoves(
	from Square,
	direction int,
	ml *MoveList,
	side int,
	ownPieces, oppPieces uint64,
	promotionRank int,
	doublePushRank int,
) {
	oneStep := from + Square(direction)
	twoStep := oneStep + Square(direction)
	rank := int(from) / 8

	if oneStep < a8 || oneStep > h1 {
		return
	}

	// QUIET MOVESSS
	if !board.Occupied(oneStep) {
		if rank == promotionRank {
			// promotion pushes
			ml.Add(from, oneStep, PromotionQ)
			ml.Add(from, oneStep, PromotionR)
			ml.Add(from, oneStep, PromotionB)
			ml.Add(from, oneStep, PromotionN)
		} else {
			// one step push
			ml.Add(from, oneStep, Normal)
			// double push
			if rank == doublePushRank && !board.Occupied(twoStep) {
				ml.Add(from, twoStep, DoublePawnPush)
			}
		}
	}

	// CAPTURESSSS
	captures := pawnAttacks[side][from] & oppPieces

	for captures != 0 {
		to := Square(getIndexOfLS1B(captures))
		// Pawn Capture Promotion
		if rank == promotionRank {
			ml.Add(from, to, PromotionCaptureQ)
			ml.Add(from, to, PromotionCaptureR)
			ml.Add(from, to, PromotionCaptureB)
			ml.Add(from, to, PromotionCaptureN)
		} else {
			// normal pawn capture
			ml.Add(from, to, Capture)
		}
		captures = popBit(captures, to)
	}

	// en passant
	if board.State.EnPassantSquare != none {
		enPassantAttacks := pawnAttacks[side][from] & (1 << board.State.EnPassantSquare)

		if enPassantAttacks != 0 {
			targetEnPassant := getIndexOfLS1B(enPassantAttacks)
			ml.Add(from, Square(targetEnPassant), EnPassant)
		}
	}

}

func (board *Board) generateKnightMoves(
	from Square,
	ml *MoveList,
	ownPieces uint64,
	side int,
) {
	attacks := knightAttacks[from] &^ ownPieces
	ml.PushSetwiseFlag(from, attacks, func(to Square) MoveFlag {
		if board.OccupiedByOpp(to, side) {
			return Capture
		}
		return Normal
	})
}