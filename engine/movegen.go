package engine

func GeneratePawnMoves(board *Board, side int, ml *MoveList) {

}

func GenerateKnightMoves(board *Board, side int, ml *MoveList) {

}

func GenerateBishopMoves(board *Board, side int, ml *MoveList) {

}

func GenerateRookMoves(board *Board, side int, ml *MoveList) {

}

func GenerateQueenMoves(board *Board, side int, ml *MoveList) {

}

func GenerateKingMoves(board *Board, side int, ml *MoveList) {

}

func GenerateCastlingMoves(board *Board, side int, ml *MoveList) {

}

func (board *Board) GenerateMoves(ml *MoveList) {
	side := board.State.SideToMove
	ownPieces := board.OccupancyBitboards[side]
	oppPieces := board.OccupancyBitboards[side^1]

	for piece := P; piece <= k; piece++ {
		bitboard := board.Bitboards[piece]

		for bitboard != 0 {
			fromIdx := getIndexOfLS1B(bitboard)
			from := Square(fromIdx)

			switch piece {
			case P:
				board.generatePawnMoves(from, -8, ml, side, ownPieces, oppPieces, 6, 1)
			case p:
				board.generatePawnMoves(from, 8, ml, side, ownPieces, oppPieces, 1, 6)
			// UNREVISED !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			case N, n:
				board.generateKnightMoves(from, ml, ownPieces, side)
			case B, b:
				attacks := bishopAttacks(from, board.OccupancyBitboards[both]) &^ ownPieces
				ml.PushSetwiseFlag(from, attacks, func(to Square) MoveFlag {
					if board.OccupiedByOpp(to, side) {
						return Capture
					}
					return Normal
				})
			case R, r:
				attacks := rookAttacks(from, board.OccupancyBitboards[both]) &^ ownPieces
				ml.PushSetwiseFlag(from, attacks, func(to Square) MoveFlag {
					if board.OccupiedByOpp(to, side) {
						return Capture
					}
					return Normal
				})
			case Q, q:
				attacks := GetQueenAttacks(from, board.OccupancyBitboards[both]) &^ ownPieces
				ml.PushSetwiseFlag(from, attacks, func(to Square) MoveFlag {
					if board.OccupiedByOpp(to, side) {
						return Capture
					}
					return Normal
				})
			case K, k:
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
		board.generateCastlingMoves(ml, side)

	}

}

func (board *Board) generateCastlingMoves(ml *MoveList, side int) {
	castle := board.State.CastlingRights
	eSquareKingside := e8
	if side == white {
		eSquareKingside = e1
	}
	eSquareQueenside := b8
	if side == white {
		eSquareQueenside = b1
	}
	switch {
	case (castle&whiteKingside != 0) || (castle&blackKingside != 0):
		// ensure square between king and king's rook are empty
		if !board.Occupied(eSquareKingside+1) && !board.Occupied(eSquareKingside+2) {
			// ensure king and the kingside squares are not under attacks
			if !board.isSquareAttacked(eSquareKingside) && !board.isSquareAttacked(eSquareKingside+1) {
				ml.Add(eSquareKingside, eSquareKingside+2, KingCastle)
			}
		}
	case (castle&whiteQueenside != 0) || (castle&blackQueenside != 0):
		// ensure square between queen and queen's rook are empty
		if !board.Occupied(eSquareQueenside) && !board.Occupied(eSquareQueenside+1) && !board.Occupied(eSquareQueenside+2) {
			// ensure king and the kingside squares are not under attacks
			if !board.isSquareAttacked(eSquareQueenside+2) && !board.isSquareAttacked(eSquareQueenside+3) {
				ml.Add(eSquareQueenside+3, eSquareQueenside+1, QueenCastle)
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
	captures := pawnAttacks[board.State.SideToMove][from] & oppPieces

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
