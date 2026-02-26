package main

import "math/bits"

func main() {
	// var bitboard uint64 = 71776119061217280
	// var bitboard uint64 = 0

	// bitboard = setBit(bitboard, e4)
	// bitboard = setBit(bitboard, c3)
	// bitboard = setBit(bitboard, f2)
	// PrintBitboard(bitboard)
	// bitboard = popBit(bitboard, e4)

	// PrintBitboard(bitboard)

	// for rank := range 8 {
	// 	for file := range 8 {
	// 		square := Square(rank*8 + file)
	// 		if file == 7 {
	// 			bitboard = setBit(bitboard, square)
	// 		}

	// 	}
	// }
	//
	// corners := []uint{0, 7, 56, 63}

	// for _, pos := range corners {
	// 	bitboard |= 1 << pos
	// }
	// PrintBitboard(bitboard)

	// for u := range 8 {
	// 	PrintBitboard((MainDiagonalMasked & NotAFile) >> u)
	// }
	//
	// var mask uint64 = 0

	// for file := range 8 {
	// 	for rank := range 8 {

	// 		if file+rank == 7 {
	// 			sq := rank*8 + file
	// 			mask = setBit(mask, Square(sq))
	// 		}
	// 	}
	// }

	println("Pop count: ", popCount(MainDiagonal))
	println("Index: ", getIndexOfLS1B(17781434089472), " Coordinate: ", BitboardSquares[getIndexOfLS1B(17781434089472)])

	initLeaperAttacks()

	// for sq := range 64 {
	// PrintBitboard(rankMask(Square(sq)) | fileMask(Square(sq)))
	// PrintBitboard(maskRookAttacks(Square(sq)))
	// }

	// PrintBitboard(maskRookAttacks(Square(e4)))

	PrintBitboard(rookAttacks(Square(e4), setBit(setBit(setBit(setBit(setBit(0, b6), c4), b4), f4), g3)))

	attackMask := maskRookAttacks(a1)
	for idx := range 4096 {
		PrintBitboard(SetOccupancy(idx, popCount(attackMask), attackMask))
	}
	// PrintBitboard(occ)
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
