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
				board.generatePawnMoves(from, -8, piece, ml, side, ownPieces, oppPieces, 6, 1)
			case p:
				board.generatePawnMoves(from, 8, piece, ml, side, ownPieces, oppPieces, 1, 6)
			// UNREVISED !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			case N, n:
				attacks := knightAttacks[from] &^ ownPieces
				ml.PushSetwise(from, attacks, Normal)
			case B, b:
				attacks := bishopAttacks(from, board.OccupancyBitboards[both]) &^ ownPieces
				ml.PushSetwise(from, attacks, Normal)
			case R, r:
				attacks := rookAttacks(from, board.OccupancyBitboards[both]) &^ ownPieces
				ml.PushSetwise(from, attacks, Normal)
			case Q, q:
				attacks := GetQueenAttacks(from, board.OccupancyBitboards[both]) &^ ownPieces
				ml.PushSetwise(from, attacks, Normal)
			case K, k:
				attacks := kingAttacks[from] &^ ownPieces
				ml.PushSetwise(from, attacks, Normal)
			}

			bitboard = popBit(bitboard, from)
		}
	}

	// board.generateCastlingMoves(ml, side)
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
	piece Piece,
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
