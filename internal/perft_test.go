package internal_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/alvissraghnall/stewfish/engine"
)

func init() {
	engine.Init()
}

func TestPerftPosition1(t *testing.T) {
	board := engine.NewBoard()
	board.FenSetup("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	expected := []uint64{
		20, 400, 8902, 197281, 4865609, 119060324, 3195901860, 84998978956, 2439530234167,
	}

	for depth, want := range expected {
		t.Run(fmt.Sprintf("Depth%d", depth+1), func(t *testing.T) {
			if depth+1 > 5 && testing.Short() {
				t.Skip("skipping deep perft in short mode")
			}
			nodes := perftTestInternal(depth+1, board)
			if nodes != want {
				t.Errorf("perft depth %d = %d, want %d", depth+1, nodes, want)
			}
		})
	}
}

func TestPerftPosition2(t *testing.T) {
	board := engine.NewBoard()
	board.FenSetup("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")

	expected := []uint64{
		48, 2039, 97862, 4085603, 193690690, 8031647685,
	}

	for depth, want := range expected {
		t.Run(fmt.Sprintf("Depth%d", depth+1), func(t *testing.T) {
			if depth+1 > 4 && testing.Short() {
				t.Skip("skipping deep perft in short mode")
			}
			nodes := perftTestInternal(depth+1, board)
			if nodes != want {
				t.Errorf("perft depth %d = %d, want %d", depth+1, nodes, want)
			}
		})
	}
}

func TestPerftPosition3(t *testing.T) {
	board := engine.NewBoard()
	board.FenSetup("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")

	expected := []uint64{
		14, 191, 2812, 43238, 674624, 11030083, 178633661, 3009794393,
	}

	for depth, want := range expected {
		t.Run(fmt.Sprintf("Depth%d", depth+1), func(t *testing.T) {
			if depth+1 > 5 && testing.Short() {
				t.Skip("skipping deep perft in short mode")
			}
			nodes := perftTestInternal(depth+1, board)
			if nodes != want {
				t.Errorf("perft depth %d = %d, want %d", depth+1, nodes, want)
			}
		})
	}
}

func TestPerftPosition4(t *testing.T) {
	board := engine.NewBoard()
	board.FenSetup("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")

	expected := []uint64{
		6, 264, 9467, 422333, 15833292, 706045033,
	}

	for depth, want := range expected {
		t.Run(fmt.Sprintf("Depth%d", depth+1), func(t *testing.T) {
			if depth+1 > 4 && testing.Short() {
				t.Skip("skipping deep perft in short mode")
			}
			nodes := perftTestInternal(depth+1, board)
			if nodes != want {
				t.Errorf("perft depth %d = %d, want %d", depth+1, nodes, want)
			}
		})
	}
}

func TestGenerateMovesStartingPosition(t *testing.T) {
	board := engine.NewBoard()
	board.FenSetup("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	fmt.Printf("white knights bits: %064b\n", board.Bitboards[engine.N])
	fmt.Printf("white bishops bits: %064b\n", board.Bitboards[engine.B])
	fmt.Printf("white rooks bits: %064b\n", board.Bitboards[engine.R])

	var ml engine.MoveList
	board.GenerateMoves(&ml)
	for _, move := range ml.SliceMoves() {
		t.Log(move.DebugString(board.State.SideToMove, board))
	}
	if got, want := ml.Len(), 20; got != want {
		t.Fatalf("generatemoves initial position = %d, want %d", got, want)
	}
}

func perfTest(board *engine.Board, depth int) {
	divider := strings.Repeat("-", 75)
	fmt.Println(divider)
	fmt.Printf("%5s | %15s | %12s | %15s | %15s\n", "Index", "Move", "Nodes", "Elapsed", "MNPS")
	fmt.Println(divider)

	start := time.Now()
	nodes, index := 0, 0

	var ml engine.MoveList
	board.GenerateMoves(&ml)

	for _, move := range ml.SliceMoves() {
		now := time.Now()

		board.MakeMove(move)
		count := perftTestInternal(depth-1, board)
		nodes += int(count)
		board.UndoMove(move)

		elapsed := time.Since(now)
		mnps := 0.0
		if elapsed.Seconds() > 0 {
			mnps = float64(count) / elapsed.Seconds() / 1e6
		}

		fmt.Printf("%5d | %15s | %12d | %15s | %15.2f\n", index, move.DebugString(board.State.SideToMove, board), count, elapsed, mnps)

		index++
	}

	end := time.Since(start)
	mnps := 0.0
	if end.Seconds() > 0 {
		mnps = float64(nodes) / end.Seconds() / 1e6
	}

	fmt.Println(divider)
	fmt.Printf("%5s | %15s | %12d | %15s | %15.2f\n", "", "Total", nodes, end, mnps)
	fmt.Println(divider)

}

func TestPerftDivide(t *testing.T) {
	board := engine.NewBoard()
	board.FenSetup(engine.FenStartPosition)

	perfTest(board, 6)
}

func perftTestInternal(depth int, board *engine.Board) uint64 {
	if depth == 0 {
		return 1
	}

	var ml engine.MoveList
	board.GenerateMoves(&ml)

	var nodes uint64

	if depth == 1 {
		for i := 0; i < ml.Len(); i++ {
			move := ml.GetMove(i)
			board.MakeMove(move)
			if !board.IsInCheck(board.State.SideToMove ^ 1) {
				nodes++
			}
			board.UndoMove(move)
		}
		return nodes
	}

	for i := 0; i < ml.Len(); i++ {
		move := ml.GetMove(i)
		board.MakeMove(move)
		if !board.IsInCheck(board.State.SideToMove ^ 1) {
			nodes += perftTestInternal(depth-1, board)
		}
		board.UndoMove(move)
	}

	return nodes
}

func BenchmarkPerft(b *testing.B) {
	benchmarks := []struct {
		name  string
		fen   string
		depth int
	}{
		{"StartPosition_Depth4", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 4},
		{"Kiwipete_Depth4", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 4},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			board := engine.NewBoard()
			board.FenSetup(bm.fen)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				perftTestInternal(bm.depth, board)
			}
		})
	}
}
