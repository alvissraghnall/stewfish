package main

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
	// 		if file >= 2 {
	// 			bitboard = setBit(bitboard, square)
	// 		}

	// 	}
	// }

	// PrintBitboard(bitboard)
	//
	// PrintBitboard(maskPawnAttacks(white, g4))
	// PrintBitboard(maskPawnAttacks(white, h4))
	// PrintBitboard(maskPawnAttacks(white, a2))
	// PrintBitboard(maskPawnAttacks(white, b1))
	// PrintBitboard(maskPawnAttacks(black, g4))
	// PrintBitboard(maskPawnAttacks(black, h4))
	// PrintBitboard(maskPawnAttacks(black, a2))
	// PrintBitboard(maskPawnAttacks(black, b1))
	// for i := 8; i >= 1; i-- {
	// 	fmt.Printf("a%d \nb%[1]d \nc%[1]d \nd%[1]d \ne%[1]d \nf%[1]d \ng%[1]d \nh%[1]d\n", i)
	// }
	//
	initLeaperAttacks()

	for sq := range 64 {
		PrintBitboard(kingAttacks[sq])
	}
}
