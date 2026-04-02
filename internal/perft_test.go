package internal_test

import (
	"fmt"
	"time"

	"github.com/alvissraghnall/stewfish/engine"
)

func perfTest (board *engine.Board, depth int) int {
	start := time.Now()
	// nodes, index := 0, 0

	ml := engine.NewMoveList()
	board.GenerateMoves(ml)

	// for _, move := range ml.Slice() {

	// }

	
	end := time.Now()

	fmt.Printf("Move generation took %s\n", end.Sub(start))
	return 0
}