package main

const (
	NOT_A_FILE  uint64 = 18374403900871474942
	NOT_H_FILE  uint64 = 9187201950435737471
	NOT_GH_FILE uint64 = 4557430888798830399
	NOT_AB_FILE uint64 = 18229723555195321596
)

const (
	white = iota
	black
)

var pawnAttacks [2][64]uint64
var knightAttacks [64]uint64
var kingAttacks [64]uint64

/**
VISUALIZATION:

Horizontal
East (->) = << 1
West (<-) = >> 1
South (DOWN) = << 8
North (UP) = >> 8
*/

func maskPawnAttacks(side int, square Square) uint64 {
	var attacks uint64 = 0

	var bitboard uint64 = 0

	bitboard = setBit(bitboard, square)

	if side == 0 {
		if ((bitboard >> 7) & NOT_A_FILE) != 0 {
			attacks |= (bitboard >> 7)
		}
		if ((bitboard >> 9) & NOT_H_FILE) != 0 {
			attacks |= (bitboard >> 9)
		}
	} else {
		if ((bitboard << 7) & NOT_H_FILE) != 0 {
			attacks |= (bitboard << 7)
		}
		if ((bitboard << 9) & NOT_A_FILE) != 0 {
			attacks |= (bitboard << 9)
		}
	}

	return attacks
}

func maskKnightAttacks(square Square) uint64 {
	var attacks uint64 = 0

	var bitboard uint64 = 0

	bitboard = setBit(bitboard, square)

	attacks |= (bitboard & NOT_H_FILE) << 17
	attacks |= (bitboard & NOT_A_FILE) << 15
	attacks |= (bitboard & NOT_H_FILE) >> 15
	attacks |= (bitboard & NOT_A_FILE) >> 17

	attacks |= (bitboard & NOT_GH_FILE) << 10
	attacks |= (bitboard & NOT_AB_FILE) << 6
	attacks |= (bitboard & NOT_GH_FILE) >> 6
	attacks |= (bitboard & NOT_AB_FILE) >> 10

	return attacks

}

func maskKingAttacks(square uint8) uint64 {
	king := uint64(1) << square

	// horizontal
	attacks := ((king & NOT_H_FILE) << 1) | // east
		((king & NOT_A_FILE) >> 1) // west

	span := king | attacks

	// vertical + diagonals
	attacks |= (span << 8) | (span >> 8)

	return attacks
}

func initLeaperAttacks() {
	for square := range 64 {
		pawnAttacks[white][square] = maskPawnAttacks(white, Square(square))
		pawnAttacks[black][square] = maskPawnAttacks(black, Square(square))
		knightAttacks[square] = maskKnightAttacks(Square(square))
		kingAttacks[square] = maskKingAttacks(uint8(square))
	}
}
