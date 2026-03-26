package engine

import (
	"fmt"
)

type Square uint8

const (
	a8 Square = iota
	b8
	c8
	d8
	e8
	f8
	g8
	h8
	a7
	b7
	c7
	d7
	e7
	f7
	g7
	h7
	a6
	b6
	c6
	d6
	e6
	f6
	g6
	h6
	a5
	b5
	c5
	d5
	e5
	f5
	g5
	h5
	a4
	b4
	c4
	d4
	e4
	f4
	g4
	h4
	a3
	b3
	c3
	d3
	e3
	f3
	g3
	h3
	a2
	b2
	c2
	d2
	e2
	f2
	g2
	h2
	a1
	b1
	c1
	d1
	e1
	f1
	g1
	h1
	none Square = 255
)

const (
	whiteKingside = 1 << iota
	whiteQueenside
	blackKingside
	blackQueenside
)

type Piece uint8

const (
	P Piece = iota
	N
	B
	R
	Q
	K
	p
	n
	b
	r
	q
	k
	Zilch
)

var asciiPieces = []rune{'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k'}

var unicodePieces = []rune{'♙', '♘', '♗', '♖', '♕', '♔', '♟', '♞', '♝', '♜', '♛', '♚'}

var CharPieceMap = map[rune]Piece{
	'P': P,
	'N': N,
	'B': B,
	'R': R,
	'Q': Q,
	'K': K,
	'p': p,
	'n': n,
	'b': b,
	'r': r,
	'q': q,
	'k': k,
}

func PrintBitboard(bitboard uint64) {
	fmt.Println()
	for rank := range 8 {
		for file := range 8 {

			var square Square = Square(rank*8 + file)

			if file == 0 {
				fmt.Printf("  %d ", 8-rank)
			}

			bit := getBit(bitboard, square)
			fmt.Printf(" %d", bit)
		}
		fmt.Println()
	}

	// print board files
	fmt.Println("\n     a b c d e f g h")

	fmt.Println()
	fmt.Printf("     Bitboard: %d \n\n", bitboard)

}

func getBit(bitboard uint64, square Square) uint64 {
	return (bitboard >> uint64(square)) & uint64(1)
}

func setBit(bitboard uint64, square Square) uint64 {
	return bitboard | (uint64(1) << square)
}

func popBit(bitboard uint64, square Square) uint64 {
	return bitboard &^ (uint64(1) << square)
}

func toggleBit(bitboard uint64, square Square) uint64 {
	return bitboard ^ (uint64(1) << square)
}
