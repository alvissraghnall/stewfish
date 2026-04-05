package engine

import (
	"fmt"
	"strings"
)

type UciParseError struct {
	Part    string
	Message string
}

func (e *UciParseError) Error() string {
	return fmt.Sprintf("Part %s: %s", e.Part, e.Message)
}

func parseMove(moveStr string, board *Board) (Move, bool) {
	if len(moveStr) < 4 || len(moveStr) > 5 {
		return Move(0), false
	}

	var ml MoveList
	board.GenerateMoves(&ml)
	ml.Print(board.State.SideToMove, board)

	from := Square((rune(moveStr[0]) - 'a') + (8-(rune(moveStr[1])-'0'))*8)
	to := Square((rune(moveStr[2]) - 'a') + (8-(rune(moveStr[3])-'0'))*8)

	if from > 63 || to > 63 {
		return Move(0), false
	}

	var promotionPiece Piece

	for i := 0; i < ml.Len(); i++ {
		move := ml.GetMove(i)

		if move.getFrom() == from && move.getTo() == to {
			promotionPiece = move.PromotionPiece(board.State.SideToMove)

			if promotionPiece != Zilch {
				if len(moveStr) != 5 {
					return Move(0), false
				}

				println(string(promotedPieces[promotionPiece]), string(moveStr[4]))

				switch moveStr[4] {
				case 'q':
					if promotionPiece == Q || promotionPiece == q {
						return move, true
					}
				case 'r':
					if promotionPiece == R || promotionPiece == r {
						return move, true
					}
				case 'b':
					if promotionPiece == B || promotionPiece == b {
						return move, true
					}
				case 'n':
					if promotionPiece == N || promotionPiece == n {
						return move, true
					}
				}

				continue
			}

			// move is NOT a promotion
			if len(moveStr) == 5 {
				// extra promotion char on non-promotion move , so invalid
				return Move(0), false
			}

			return move, true
		}
	}

	return Move(0), false
}

func parsePosition (command string) (Board, error) {
	var board Board

	commandParts := strings.Split(command, " ")

	outer:
	for len(commandParts) > 0 {
		switch commandParts[0] {
		case "startpos":
			board.FenSetup(FenStartPosition)
			commandParts = commandParts[1:]
		case "fen":
			fen := strings.Join(commandParts[1:7], " ")
			board.FenSetup(fen)
			commandParts = commandParts[7:]
		case "moves":
			for i := 1; i < len(commandParts); i++ {
				moveStr := commandParts[i]
				move, valid := parseMove(moveStr, &board)
				if !valid {
					return Board{}, &UciParseError{Part: "position", Message: fmt.Sprintf("Invalid move in position command: %s", moveStr)}
				}
				board.MakeMove(move)
			}
			commandParts = []string{} // all parts processed
		default:
			// println("Unknown token in position command:", commandParts[0])
			// return
			commandParts = commandParts[1:] // skip unknown token
			continue outer
		}
	}

	println("Position parsed successfully. Current board state:")
	board.PrintBoardWithPieces()

	return board, nil
}
