package engine

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/alvissraghnall/stewfish/internal"
)

func init() {
	Init()
}

func TestParseUci(t *testing.T) {
	board := NewBoard()
	board.FenSetup(FenStartPosition)

	// white
	tests := []struct {
		moveStr string
		want    bool
	}{
		{"e2e4", true},
		{"c2c3", true},
		{"h7h6", false}, // illegal (white to move)
		{"e7e5", false}, // illegal (white to move)
		{"g1f3", true},  // valid knight move (white to move)
		{"b8c6", false}, // illegal (white to move)
		{"e2e5", false}, // illegal move
	}

	for _, tt := range tests {
		t.Run(tt.moveStr, func(t *testing.T) {
			move, got := parseMove(tt.moveStr, board)
			if got != tt.want {
				t.Errorf("ParseMove(%q) = %v, want %v", tt.moveStr, got, tt.want)
			} else if got {
				t.Logf("Parsed move: %s", move.DebugString(board.State.SideToMove, board))
			}
		})
	}

	// black
	board.State.SideToMove = black

	testsBlack := []struct {
		moveStr string
		want    bool
	}{
		{"e7e5", true},
		{"c7c6", true},
		{"e2e4", false}, // illegal (black to move)
		{"g1f3", false}, // illegal (black to move)
		{"b8c6", true},  // valid knight move (black to move)
		{"e7e4", false}, // illegal move
	}

	for _, tt := range testsBlack {
		t.Run(tt.moveStr, func(t *testing.T) {
			move, got := parseMove(tt.moveStr, board)
			if got != tt.want {
				t.Errorf("ParseMove(%q) = %v, want %v", tt.moveStr, got, tt.want)
			} else if got {
				t.Logf("Parsed move: %s", move.DebugString(board.State.SideToMove, board))
			}
		})
	}
}

func TestUciMessageLoop(t *testing.T) {
	// mock standard input
	oldStdin := os.Stdin
	rIn, wIn, _ := os.Pipe()
	os.Stdin = rIn

	// mock standard output and error
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	wIn.WriteString("uci\n")
	wIn.WriteString("isready\n")
	wIn.WriteString("quit\n")
	wIn.Close()

	var emptyBuffer internal.Deque[string]
	UciMessageLoop(emptyBuffer)

	wOut.Close()
	wErr.Close()
	os.Stdin = oldStdin
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var outBuf, errBuf bytes.Buffer
	io.Copy(&outBuf, rOut)
	io.Copy(&errBuf, rErr)

	if !strings.Contains(outBuf.String(), "option name Threads") {
		t.Errorf("Expected threads option in stdout, got: %s", outBuf.String())
	}
	if !strings.Contains(outBuf.String(), "uciok") {
		t.Errorf("Expected uciok in stdout, got: %s", outBuf.String())
	}
}

func TestParseUciPromotion(t *testing.T) {
	board := NewBoard()
	fen := "8/P7/8/8/8/8/7p/8 w - - 0 1"
	board.FenSetup(fen)

	// white
	tests := []struct {
		moveStr string
		want    bool
	}{
		{"a7a8q", true},
		{"a7a8r", true},
		{"a7a8b", true},
		{"a7a8n", true},
		{"h2h1q", false}, // illegal (white to move)
		{"a7a8k", false}, // invalid promotion piece
	}

	for _, tt := range tests {
		t.Run(tt.moveStr, func(t *testing.T) {
			move, got := parseMove(tt.moveStr, board)
			if got != tt.want {
				t.Errorf("ParseMove(%q) = %v, want %v", tt.moveStr, got, tt.want)
			} else if got {
				t.Logf("Parsed move: %s", move.DebugString(board.State.SideToMove, board))
			}
		})
	}

	// black
	board.State.SideToMove = black

	testsBlack := []struct {
		moveStr string
		want    bool
	}{
		{"h2h1q", true},
		{"h2h1r", true},
		{"h2h1b", true},
		{"h2h1n", true},
		{"a7a8q", false}, // illegal (black to move)
		{"h2h1k", false}, // invalid promotion piece
	}

	for _, tt := range testsBlack {
		t.Run(tt.moveStr, func(t *testing.T) {
			move, got := parseMove(tt.moveStr, board)
			if got != tt.want {
				t.Errorf("ParseMove(%q) = %v, want %v", tt.moveStr, got, tt.want)
			} else if got {
				t.Logf("Parsed move: %s", move.DebugString(board.State.SideToMove, board))
			}
		})
	}
}

func TestParseUciCastling(t *testing.T) {
	board := NewBoard()
	board.FenSetup("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1")

	tests := []struct {
		moveStr string
		want    bool
	}{
		{"e1g1", true},  // white kingside castle
		{"e1c1", true},  // white queenside castle
		{"e8g8", false}, // illegal (white to move)
		{"e8c8", false}, // illegal (white to move)
	}

	for _, tt := range tests {
		t.Run(tt.moveStr, func(t *testing.T) {
			move, got := parseMove(tt.moveStr, board)
			if got != tt.want {
				t.Errorf("ParseMove(%q) = %v, want %v", tt.moveStr, got, tt.want)
			} else if got {
				t.Logf("Parsed move: %s", move.DebugString(board.State.SideToMove, board))
			}
		})
	}
}

func TestParsePosition(t *testing.T) {
	tests := []struct {
		command   string
		wantMoves int
	}{
		{"position fen r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", 0},
		{"position startpos moves e2e4", 1},
		{"position startpos moves e2e4 e7e5", 2},
		{"position startpos moves e2e4 e7e5 g1f3", 3},
		{"position startpos moves e2e4 e7e5 g1f3 axc8", 3},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			ctx := NewContext()
			cmd := ParseUCI(tt.command, ctx)
			if len(cmd.Moves) != tt.wantMoves {
				t.Fatalf("unexpected moves length: got %d, want %d", len(cmd.Moves), tt.wantMoves)
			}
			if ctx.board == nil {
				t.Fatalf("expected ctx.board to be initialized")
			}
		})
	}
}
