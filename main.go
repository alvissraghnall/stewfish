package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/alvissraghnall/stewfish/engine"
)

func init() {
	engine.Init()
}

func main() {
	// PrintBitboard(MainDiagonal)
	// println("Pop count: ", popCount(MainDiagonal))
	// println("Index: ", getIndexOfLS1B(17781434089472), " Coordinate: ", BitboardSquares[getIndexOfLS1B(17781434089472)])

	// PrintBitboard(rookAttacks(Square(e4), setBit(setBit(setBit(setBit(setBit(0, b6), c4), b4), f4), g3)))

	// attackMask := diagonalMask(d4) | antiDiagonalMask(d4)
	// for idx := range 100 {
	// 	PrintBitboard(SetOccupancy(idx, popCount(attackMask), attackMask))
	// }
	// PrintBitboard(occ)

	// PrintBitboard(generateMagicNumber())

	// initMagicNumbers()

	// PrintBoardWithPieces()

	// fen := "r1bqkbnr/pppppppp/2n5/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq e3 0 2"
	// fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	// fen := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPpP/R3K2R w KQkq a3 0 1 "
	fen := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1 "
	// fen := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPpP/R3K2R b KQkq - 0 1 "

	var board *engine.Board = engine.NewBoard()

	err := board.FenSetup(fen)
	if err != nil {
		panic(err)
	}
	board.PrintBoardWithPieces()

	var ml engine.MoveList

	start := time.Now()
	board.GenerateMoves(&ml)

	end := time.Now()

	ml.Print(0, board)

	fmt.Printf("Move generation took %s\n", end.Sub(start))

}

func PrintJSON(obj any) {
	bytes, _ := json.MarshalIndent(obj, "\t", "\t")
	fmt.Println(string(bytes))
}
