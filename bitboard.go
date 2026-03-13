package main

import (
	"fmt"
	"runtime"
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

const (
	P = iota
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
)

var Bitboards [12]uint64

var OccupancyBitboards [3]uint64

var SideToMove int

var EnPassantSquare Square = none

var CastlingRights int

var asciiPieces = []rune{'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k'}

var unicodePieces = []rune{'♙', '♘', '♗', '♖', '♕', '♔', '♟', '♞', '♝', '♜', '♛', '♚'}

var CharPieceMap = map[rune]int{
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

func PrintBoardWithPieces() {
	fmt.Println()
	for rank := range 8 {
		for file := range 8 {

			var square Square = Square(rank*8 + file)

			if file == 0 {
				fmt.Printf("  %d ", 8-rank)
			}

			piece := getPieceOnSquare(square)
			if piece == -1 {
				fmt.Printf(" .")
				continue
			}
			if runtime.GOOS == "windows" {
				fmt.Printf(" %c", asciiPieces[piece])
			} else {
				fmt.Printf(" %c", unicodePieces[piece])
			}
		}
		fmt.Println()
	}

	// print board files
	fmt.Println("\n     a b c d e f g h")

	fmt.Println()

	if SideToMove == white {
		fmt.Println("Side to move: White")
	} else {
		fmt.Println("Side to move: Black")
	}

	EnPassantSquare = f2
	CastlingRights = whiteKingside | blackKingside | blackQueenside

	if EnPassantSquare != none {
		fmt.Printf("En Passant Square: %s\n", BitboardSquares[EnPassantSquare])
	} else {
		fmt.Printf("En Passant Square: None\n")
	}

	fmt.Printf("Castling Rights: %s\n", getCastlingRightsString())

}

func getCastlingRightsString() string {
	var rights string
	if CastlingRights&whiteKingside != 0 {
		rights += "K"
	}
	if CastlingRights&whiteQueenside != 0 {
		rights += "Q"
	}
	if CastlingRights&blackKingside != 0 {
		rights += "k"
	}
	if CastlingRights&blackQueenside != 0 {
		rights += "q"
	}
	return rights
}

func getPieceOnSquare(square Square) int {
	for piece := range k + 1 {
		if getBit(Bitboards[piece], square) == 1 {
			return piece
		}
	}

	return -1
}

func ParseFEN(fen string) {
	Bitboards = [12]uint64{}
	OccupancyBitboards = [3]uint64{}

	SideToMove = white
	EnPassantSquare = none
	CastlingRights = 0

	var rank int = 7
	var file int = 0

	for _, char := range fen {
		if char == ' ' {
			break
		}

		if char == '/' {
			rank--
			file = 0
			continue
		}

		if char >= '1' && char <= '8' {
			file += int(char - '0')
			continue
		}

		if char == 'w' {
			SideToMove = white
			continue
		}

		if char == 'b' {
			SideToMove = black
			continue
		}

		piece, ok := CharPieceMap[char]
		if !ok {
			panic(fmt.Sprintf("Invalid FEN character: %c", char))
		}

		square := Square(rank*8 + file)
		Bitboards[piece] = setBit(Bitboards[piece], square)
		file++
	}

	fmt.Println("FEN parsed successfully!")
}
