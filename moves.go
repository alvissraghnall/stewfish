package main

func (board *Board) GenerateMoves () {
	var fromSquare Square
	var toSquare int
	var bitboard, attacks uint64

	for piece := P; piece < k; piece++ {
		bitboard = board.Bitboards[piece]

		if board.State.SideToMove == white {
			if piece == P {
				for bitboard != 0 {
					idxLSB := getIndexOfLS1B(bitboard)
					if idxLSB == -1 {
						break
					}
					fromSquare = Square(idxLSB)
					toSquare = int(fromSquare) - 8

					if !(toSquare < 0) && (getBit(board.OccupancyBitboards[both], Square(toSquare)) != 0) {
						// pawn promotion, one step ahead or 2
						if fromSquare >= a7 && fromSquare <= h7 {
							println("PAWN promotion: ", BitboardSquares[fromSquare], BitboardSquares[toSquare])

						}
					}

					bitboard = popBit(bitboard, fromSquare)
					println(attacks)
				}
			}
		}
		
	}
}