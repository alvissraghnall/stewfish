package engine

var pieceValues [12]int = [12]int{
	100, 350, 350, 525, 1000, 10000, // white
	-100, -350, -350, -525, -1000, 10000, // black
}

func evaluatePosition(board *Board) int {
	score := 0
	var bb uint64

	var piece Piece
	var square Square

	for bbPiece := P; bbPiece <= k; bbPiece++ {
		bb = board.Bitboards[bbPiece]
		for bb != 0 {
			square = Square(getIndexOfLS1B(bb))
			piece = Piece(bbPiece)
			score += pieceValues[piece]
			bb = popBit(bb, square)
		}
	}
	if board.State.SideToMove == black {
		score *= -1
	}
	return score

}