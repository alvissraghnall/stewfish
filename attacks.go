package main

import "math/bits"

const (
	NotAFile         uint64 = 18374403900871474942
	NotHFile         uint64 = 9187201950435737471
	NotGhFile        uint64 = 4557430888798830399
	NotAbFile        uint64 = 18229723555195321596
	MainDiagonal     uint64 = 9241421688590303745
	AntiMainDiagonal uint64 = 72624976668147840
	Rank1            uint64 = 255                  // rank 8 acc
	Rank8            uint64 = 18374686479671623680 // rank 1 acc
	FileA            uint64 = 72340172838076673
	FileH            uint64 = 9259542123273814144
	Edges                   = FileA | FileH | Rank1 | Rank8
	Vertices         uint64 = 9295429630892703873

	MainDiagonalMasked = MainDiagonal & ^Edges
	AntiDiagonalMasked = AntiMainDiagonal & ^Edges
)

const (
	white = iota
	black
)

var BitboardSquares = [64]string{
	"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
	"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
	"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
	"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
	"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
}

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
		if ((bitboard >> 7) & NotAFile) != 0 {
			attacks |= (bitboard >> 7)
		}
		if ((bitboard >> 9) & NotHFile) != 0 {
			attacks |= (bitboard >> 9)
		}
	} else {
		if ((bitboard << 7) & NotHFile) != 0 {
			attacks |= (bitboard << 7)
		}
		if ((bitboard << 9) & NotAFile) != 0 {
			attacks |= (bitboard << 9)
		}
	}

	return attacks
}

func maskKnightAttacks(square Square) uint64 {
	bb := uint64(1) << square

	return ((bb >> 17) & NotAFile) | // Up 2, Left 1
		((bb >> 15) & NotHFile) | // Up 2, Right 1
		((bb >> 10) & NotAbFile) | // Up 1, Left 2
		((bb >> 6) & NotGhFile) | // Up 1, Right 2
		((bb << 17) & NotHFile) | // Down 2, Right 1
		((bb << 15) & NotAFile) | // Down 2, Left 1
		((bb << 10) & NotGhFile) | // Down 1, Right 2
		((bb << 6) & NotAbFile) // Down 1, Left 2
}

func maskKingAttacks(square uint8) uint64 {
	king := uint64(1) << square

	// horizontal
	attacks := ((king & NotHFile) << 1) | // east
		((king & NotAFile) >> 1) // west

	span := king | attacks

	// vertical + diagonals
	attacks |= (span << 8) | (span >> 8)

	return attacks
}

func diagonalCompute(square Square) uint64 {
	diagonal := int(square&7) - int(square>>3)

	println("square: ", square, " diag: ", diagonal)
	var attacks uint64
	if diagonal >= 0 {
		attacks = (MainDiagonal) >> uint(diagonal*8)
	} else {
		attacks = MainDiagonal << uint(-diagonal*8)
	}
	// PrintBitboard(MainDiagonalMasked)
	// PrintBitboard(attacks & ^Edges)
	return attacks & ^(1 << uint(square))
}

func diagonalMask(square Square) uint64 {
	return diagonalCompute(square) & ^Edges
}

func antiDiagonalCompute(square Square) uint64 {
	file := square & 7
	rank := square >> 3
	var diag int = 7 - int(file) - int(rank)
	println("square: ", square, " anti: ", diag, " file: ", file, " rank: ", rank)
	if diag >= 0 {
		return (AntiMainDiagonal >> uint(diag*8)) & ^(1 << uint(square))
	}
	return (AntiMainDiagonal << uint(-diag*8)) & ^(1 << uint(square))
}

func antiDiagonalMask(square Square) uint64 {
	return antiDiagonalCompute(square) & ^Edges
}

func diagAttacks(square Square, occ uint64, diagMask uint64) uint64 {
	var mask uint64 = 1 << square
	forward := occ & diagMask
	reverse := bits.ReverseBytes64(forward)
	PrintBitboard(reverse)
	forward -= mask
	reverse -= bits.ReverseBytes64(mask)
	forward ^= bits.ReverseBytes64(reverse)
	forward &= diagMask
	return forward
}

func bishopAttacks(sq Square, occ uint64) uint64 {
	return (diagAttacks(sq, occ, diagonalCompute(sq)) | diagAttacks(sq, occ, antiDiagonalCompute(sq)))
}

func fileCompute(sq Square) uint64 {
	file := sq & 7
	return (FileA << file) & ^(1 << sq)
}

func fileMask(sq Square) uint64 {
	attacks := fileCompute(sq)

	if file := sq & 7; file != 0 {
		attacks &= ^FileA
	}
	if file := sq & 7; file != 7 {
		attacks &= ^FileH
	}
	return attacks
}

func rankCompute(sq Square) uint64 {
	rank := sq >> 3
	return (Rank1 << (8 * rank)) & ^(1 << sq)
}

func rankMask(sq Square) uint64 {
	attacks := rankCompute(sq)

	if rank := sq >> 3; rank != 0 {
		attacks &= ^Rank1
	}
	if rank := sq >> 3; rank != 7 {
		attacks &= ^Rank8
	}
	return attacks
}

func maskRookAttacks(sq Square) uint64 {
	file := sq & 7
	rank := sq >> 3

	attacks := (FileA << file) | (Rank1 << (8 * rank))
	attacks &= ^(uint64(1) << uint(sq)) // remove self

	if file != 0 {
		attacks &= ^FileA
	}
	if file != 7 {
		attacks &= ^FileH
	}
	if rank != 0 {
		attacks &= ^Rank1
	}
	if rank != 7 {
		attacks &= ^Rank8
	}

	return attacks
}

func fileAttacks(sq Square, occ uint64) uint64 {
	file := fileCompute(sq)       // With edges for attack generation
	occFile := occ & fileMask(sq) // Without edges for occupancy

	var mask uint64 = 1 << sq
	forward := occFile
	reverse := bits.ReverseBytes64(forward)
	forward -= mask
	reverse -= bits.ReverseBytes64(mask)
	forward ^= bits.ReverseBytes64(reverse)
	forward &= file
	return forward
}

func rankAttacks(sq Square, occ uint64) uint64 {
	rank := rankCompute(sq)       // With edges for attack generation
	occRank := occ & rankMask(sq) // Without edges for occupancy

	var mask uint64 = 1 << sq
	forward := occRank
	reverse := bits.Reverse64(forward) // Bit reverse for ranks
	forward -= mask
	reverse -= bits.Reverse64(mask)
	forward ^= bits.Reverse64(reverse)
	forward &= rank
	return forward
}

func rookAttacks(sq Square, occ uint64) uint64 {
	file := sq & 7
	rank := sq >> 3

	// All squares on file and rank (INCLUDING edges) for attacks
	allFile := FileA << file & ^(1 << sq)
	allRank := Rank1 << (8 * rank) & ^(1 << sq)

	// Occupancy - start with all pieces
	occFile := occ
	occRank := occ

	// Remove EDGE squares from occupancy ONLY
	// File occupancy: remove edge files
	if file != 0 {
		occFile &= ^FileA
	}
	if file != 7 {
		occFile &= ^FileH
	}
	// Also remove edge ranks from file occupancy
	if rank != 0 {
		occFile &= ^Rank1
	}
	if rank != 7 {
		occFile &= ^Rank8
	}

	// Rank occupancy: remove edge ranks
	if rank != 0 {
		occRank &= ^Rank1
	}
	if rank != 7 {
		occRank &= ^Rank8
	}
	// akso remove edge files from rank occupancy
	if file != 0 {
		occRank &= ^FileA
	}
	if file != 7 {
		occRank &= ^FileH
	}

	// File attacks (vertical)... using allFile (WITH edges)
	var mask uint64 = 1 << sq
	forward := occFile & allFile // occupancy only on non-edge file squares
	reverse := bits.ReverseBytes64(forward)
	forward -= mask
	reverse -= bits.ReverseBytes64(mask)
	forward ^= bits.ReverseBytes64(reverse)
	fileAttacks := forward & allFile // Result includes edges

	// Rank attacks (horizontal) ... using allRank (WITH edges)
	forward = occRank & allRank // occupancy only on non-edge rank squares
	reverse = bits.Reverse64(forward)
	forward -= mask
	reverse -= bits.Reverse64(mask)
	forward ^= bits.Reverse64(reverse)
	rankAttacks := forward & allRank // Result includes edges

	return fileAttacks | rankAttacks
}

func initLeaperAttacks() {
	for square := range 64 {
		pawnAttacks[white][square] = maskPawnAttacks(white, Square(square))
		pawnAttacks[black][square] = maskPawnAttacks(black, Square(square))
		knightAttacks[square] = maskKnightAttacks(Square(square))
		kingAttacks[square] = maskKingAttacks(uint8(square))
	}
}

func popCount(bb uint64) int {
	count := 0
	for bb != 0 {
		count++
		bb &= bb - 1
	}
	return count
}

func getIndexOfLS1B(bb uint64) int {
	if bb == 0 {
		return -1
	}
	return bits.TrailingZeros64(bb)
}

// SetOccupancy generates occupancy variations for a given attack mask
// index: the permutation index (0 to 2^bits_in_mask - 1)
// bitsInMask: number of bits set in the attack mask
// attackMask: the mask of squares that can be occupied
func SetOccupancy(index int, bitsInMask int, attackMask uint64) uint64 {
	occupancy := uint64(0)

	for count := range bitsInMask {
		square := getIndexOfLS1B(attackMask)

		// pop LSB in attack mask
		attackMask = popBit(attackMask, Square(square))

		// make sure occupancy is on board
		if index&(1<<count) != 0 {
			// populate occupancy map
			occupancy |= 1 << square
		}
	}
	return occupancy
}
