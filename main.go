package main

import (
	"math/bits"
)

func main() {
	// PrintBitboard(MainDiagonal)
	// println("Pop count: ", popCount(MainDiagonal))
	// println("Index: ", getIndexOfLS1B(17781434089472), " Coordinate: ", BitboardSquares[getIndexOfLS1B(17781434089472)])

	initLeaperAttacks()

	// PrintBitboard(rookAttacks(Square(e4), setBit(setBit(setBit(setBit(setBit(0, b6), c4), b4), f4), g3)))

	// attackMask := diagonalMask(d4) | antiDiagonalMask(d4)
	// for idx := range 100 {
	// 	PrintBitboard(SetOccupancy(idx, popCount(attackMask), attackMask))
	// }
	// PrintBitboard(occ)

	// PrintBitboard(generateMagicNumber())

	initMagicNumbers()

}

func trimDiagonal(diagonal uint64, block uint64) uint64 {
	return diagonal & (block - 1)
}

func trimAntiDiagonal(diagonal uint64, block uint64) uint64 {
	return diagonal & ^(block - 1 - 1)
}

func trim(ray uint64, sq Square, occ uint64) uint64 {
	lower := ray & ((1 << sq) - 1)

	// squares above sq (toward MSB)
	upper := ray & ^((1 << (sq + 1)) - 1)

	blockers := lower & occ
	// PrintBitboard(upper)
	// PrintBitboard(occ)
	// PrintBitboard(-blockers)

	if blockers != 0 {
		first := blockers & -blockers
		// PrintBitboard(first)
		lower &= (first - 1)
	}
	// PrintBitboard(lower)

	blockers = upper & occ
	if blockers != 0 {
		// PrintBitboard(blockers)
		// PrintBitboard(uint64(1) << (63 - bits.LeadingZeros64(blockers)))
		msb := uint64(1) << (63 - bits.LeadingZeros64(blockers))
		upper &= ^(msb - 1) & uint64(f3)
	}
	return lower
}
