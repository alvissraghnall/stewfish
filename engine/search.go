package engine

import "math"

func negaMax(board *Board, alpha, beta, depth int) int {
	if depth == 0 {
		// return evaluate()
		return int(psqtScore(board))
	}

	var bestMove Move
	var bestMoveSoFar Move

	legalMovesFound := 0
	oldAlpha := alpha

	var ml MoveList

	board.GenerateMoves(&ml)

	for i := range ml.SliceMoves() {
		board.State.HalfMoveClock++
		if board.IsLegal(ml.GetMove(i)) {
			board.State.HalfMoveClock--
			legalMovesFound++
			continue
		}
		score := -negaMax(board, -beta, -alpha, depth-1)
		if score >= beta { 
			return beta // fail hard beta cutoff
		}
		if score > alpha {
			alpha = score
			if board.State.HalfMoveClock == 0 {
				bestMoveSoFar = ml.GetMove(i)
			}
		}
	}

	if oldAlpha != alpha {
		bestMove = bestMoveSoFar
	}

	println("Best move: ", bestMove)

	return alpha
}

func search (board *Board, depth int) {

	score := negaMax(board, int(math.Inf(-1)), 999_999, depth)

	println("score: ", score)
}