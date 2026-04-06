package engine

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alvissraghnall/stewfish/internal"
)

type Mode uint8

const (
	Uci Mode = iota
	Cli
)

type Context struct {
	status *internal.UciStatus
	board  *Board
}

func NewContext() *Context {
	return &Context{
		status: internal.NewStatus(),
	}
}

type UciParseError struct {
	Part    string
	Message string
}

func (e *UciParseError) Error() string {
	return fmt.Sprintf("Part %s: %s", e.Part, e.Message)
}

type UCICommand struct {
	Name  string
	Args  map[string]string
	Moves []Move
}

func newUCICommand(name string) UCICommand {
	return UCICommand{
		Name:  name,
		Args:  make(map[string]string),
		Moves: []Move{},
	}
}

func ParseUCI(line string, ctx *Context) UCICommand {
	tokens := strings.Fields(line)

	if len(tokens) == 0 {
		return newUCICommand("noop")
	}

	cmdName := tokens[0]
	cmd := newUCICommand(cmdName)

	// basic commands with no args
	switch cmdName {
	case "uci", "ucinewgame", "stop", "quit", "isready":
		return cmd
	}

	// complex commands requiring argument parsing
	// we pass 'tokens[1:]' which is just the arguments
	args := tokens[1:]

	switch cmdName {
	case "position":
		parsePositionArgs(ctx, &cmd, args)
	case "go":
		parseGoArgs(&cmd, args)
	case "setoption":
	}

	return cmd
}

func parseMove(moveStr string, board *Board) (Move, bool) {
	if len(moveStr) < 4 || len(moveStr) > 5 {
		return Move(0), false
	}

	var ml MoveList
	board.GenerateMoves(&ml)
	// ml.Print(board.State.SideToMove, board)

	from := Square((rune(moveStr[0]) - 'a') + (8-(rune(moveStr[1])-'0'))*8)
	to := Square((rune(moveStr[2]) - 'a') + (8-(rune(moveStr[3])-'0'))*8)

	if from > 63 || to > 63 {
		return Move(0), false
	}

	var promotionPiece Piece

	for i := 0; i < ml.Len(); i++ {
		move := ml.GetMove(i)

		// println(move.getFrom(), from, move.getTo(), to)

		if move.getFrom() == from && move.getTo() == to {
			println(move, 78998)
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
	println(97879)
	return Move(0), false
}

// parses arguments like: "startpos moves e2e4"
func parsePositionArgs(ctx *Context, cmd *UCICommand, tokens []string) {
	board := NewBoard()

	i := 0
	for i < len(tokens) {
		token := tokens[i]

		switch token {
		case "startpos":
			board.FenSetup(FenStartPosition)
			i++
		case "fen":
			fen := strings.Join(tokens[i+1:i+7], " ")
			board.FenSetup(fen)
			i += 7
		case "moves":
			i++ // skip "moves" keyword
			// The rest of the tokens should be moves
			for i < len(tokens) {
				moveStr := tokens[i]
				move, valid := parseMove(moveStr, board)

				if valid {
					cmd.Moves = append(cmd.Moves, move)

					board.MakeMove(move)
				} else {
					// Handle invalid move? UCI spec says ignore
				}
				i++
			}
		default:
			// Unknown token, skip
			i++
		}
	}

	ctx.board = board
}

func parseGoArgs(cmd *UCICommand, tokens []string) {
	i := 0
	for i < len(tokens) {
		// most go commands are key-value pairs (ex: "depth 5")
		// Check if we have a next token
		if i+1 < len(tokens) {
			key := tokens[i]
			val := tokens[i+1]

			cmd.Args[key] = val

			i += 2
		} else {
			// Standalone token, like "ponder"
			cmd.Args[tokens[i]] = "true"
			i++
		}
	}
}

func UciMessageLoop(buffer internal.Deque[string]) {
	ctx := NewContext()

	broadcastChan := newBroadcastStream(ctx)
	var mode Mode

	if buffer.Len() == 0 {
		mode = Uci
	} else {
		mode = Cli
	}

	for {
		var message string
		if m, ok := buffer.DequeueLeft(); ok {
			message = m
		} else if mode == Uci {
			var ok bool
			if message, ok = <-broadcastChan; !ok {
				break // Channel closed (EOF)
			}
		} else {
			break
		}

		cmd := ParseUCI(message, ctx)

		switch cmd.Name {
		case "uci":
			handleUci()
			mode = Uci
		case "isready":
			fmt.Printf("readyok")
		case "go":
			depthStr, ok := cmd.Args["depth"]
			if ok {
				depth, _ := strconv.Atoi(depthStr)
				fmt.Printf("Starting search depth %d\n", depth)
			}
		case "position":
			// The context (ctx.board) is already updated by ParseUCI!
			// We just do nothing, as state is set.
		case "stop":
			ctx.status.Set(internal.StatusStopped)
		case "ucinewgame":
			parsePositionArgs(ctx, &cmd, []string{"startpos"})
		case "quit":
			// dispose of board state
			return
		}

		if mode == Cli && buffer.Len() == 0 {
			break
		}
	}
}

func newBroadcastStream(ctx *Context) <-chan string {
	ch := make(chan string)
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		defer close(ch)

		var message string

		for scanner.Scan() {
			message = strings.TrimSpace(scanner.Text())

			switch message {
			case "isready":
				println("readyok")
				continue
			case "stop":
				ctx.status.Set(internal.StatusStopped)
				continue
			case "quit":
				ctx.status.Set(internal.StatusStopped)
				select {
				case ch <- "quit":
				default: // dont blovk if no listeners
				}
				return
			default:
				if ctx.status.Get() != internal.StatusRunning {
					select {
					case ch <- message:
					default: // drop silently if no listener
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return // fail silently
		}

		// EOF
		if ctx.status.Get() != internal.StatusRunning {
			ch <- "quit"
		}
	}()

	return ch
}

func handleUci() {
	fmt.Println("id name Reckless Stewfish v", internal.Version)
	fmt.Println("id author Alviss Raghnall")
	fmt.Printf("option name Threads type spin default 1 min 1 max %d", 1)

	fmt.Println("uciok")
}
