package engine

import "github.com/alvissraghnall/stewfish/internal"

const MateScore = 100_000

type Search struct {
	RootDepth int
	BestMove  Move
	nodes     uint64
	status    *internal.UciStatus
	aborted   bool
}

func (s *Search) negaMax(board *Board, alpha, beta, depth, ply int) int {
	s.nodes++
	if s.nodes&2047 == 0 && s.status.Get() != internal.StatusRunning {
		s.aborted = true
	}
	if s.aborted {
		return 0 // value is discarded by caller, see below
	}

	if depth == 0 {
		score := board.Evaluate()
		if board.State.SideToMove == black {
			score = -score
		}
		return score
	}

	oldAlpha := alpha
	legalMoves := 0
	// var bestMove Move

	var ml MoveList
	board.GenerateMoves(&ml)

	for i := range ml.SliceMoves() {
		move := ml.GetMove(i)
		if !board.MakeMove(move) {
			continue
		}
		legalMoves++
		score := -s.negaMax(board, -beta, -alpha, depth-1, ply+1)
		board.UndoMove(move)

		if s.aborted {
			return 0 // don't trust `score`, don't update alpha/bestMove, just unwind
		}

		if score > alpha {
			alpha = score
			// bestMove = move
			if depth == s.RootDepth {
				s.BestMove = move
			}
		}
		if score >= beta {
			return beta
		}
	}

	if legalMoves == 0 {
		if board.IsInCheck(board.State.SideToMove) {
			return -MateScore + ply
		}
		return 0
	}

	_ = oldAlpha
	return alpha
}

func SearchPosition(board *Board, depth int, status *internal.UciStatus) Move {
	s := Search{RootDepth: depth, status: status}
	s.negaMax(board, -MateScore, MateScore, depth, 0)
	return s.BestMove
}