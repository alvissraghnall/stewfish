package engine

import (
	"fmt"

	"github.com/alvissraghnall/stewfish/util"
)

type MagicEntry struct {
	mask   uint64
	magic  uint64
	shift  uint32
	size   uint32
	offset uint32
}

type SliderPiece = int

const (
	Rook SliderPiece = iota
	Bishop
)

func findMagicNumber(square Square, piece SliderPiece, random *util.Random) MagicEntry {
	var mask uint64

	if piece == Bishop {
		mask = diagonalMask(square) | antiDiagonalMask(square)
	} else {
		mask = maskRookAttacks(square)
	}

	relevantBits := 0
	if piece == Bishop {
		relevantBits = bishopRelevantBits[square]
	} else {
		relevantBits = rookRelevantBits[square]
	}

	calculatedBits := popCount(mask)
	if calculatedBits != relevantBits {
		panic(fmt.Sprintf("Calculated bits (%d) does not match expected relevant bits (%d) for square %d", calculatedBits, relevantBits, square))
	}

	occupancyIndicies := 1 << relevantBits

	// generate all possible occupancy variations and their resulting attacks
	var occupancies [4096]uint64
	var attacks [4096]uint64

	for i := range occupancyIndicies {
		occupancies[i] = SetOccupancy(i, relevantBits, mask)
		if piece == Bishop {
			attacks[i] = bishopAttacks(square, occupancies[i])
		} else {
			attacks[i] = rookAttacks(square, occupancies[i])
		}
	}

	// Search for a Magic Number
	for range 100_000_000 {
		// generate SPARSE random numbers.
		magic := random.XorShift64() & random.XorShift64() & random.XorShift64()

		// If the multiplication doesn't push bits into the upper part of the uint64,
		// the magic will likely cause collisions.
		if popCount((mask*magic)&0xFF00_0000_0000_0000) < 6 {
			continue
		}

		// reset the usedAttacks array inside the loop
		var usedAttacks [4096]uint64
		fail := false

		for j := range occupancyIndicies {
			// Calculate the index using the candidate magic
			index := (occupancies[j] * magic) >> (64 - relevantBits)

			if usedAttacks[index] == 0 {
				// spot is free, write the attack
				usedAttacks[index] = attacks[j]
			} else if usedAttacks[index] != attacks[j] {
				// 2 different board states map to the same index
				// but require different attack sets. This magic is bad.
				fail = true
				break
			}
		}

		if !fail {
			// we found a working magic number
			return MagicEntry{
				mask:  mask,
				magic: magic,
				shift: uint32(64 - relevantBits),
				size:  uint32(occupancyIndicies),
			}
		}
	}

	panic("Magic number not found!")
}

func initMagicNumbers() {
	rand := util.NewRandom(0xFFAAB58C5833FE89)

	var rookMagics [64]MagicEntry
	var bishopMagics [64]MagicEntry

	fmt.Println("// Rook Magics")
	offset := uint32(0)
	for sq := range 64 {
		entry := findMagicNumber(Square(sq), Rook, rand)
		entry.offset = offset
		rookMagics[sq] = entry

		fmt.Printf("{ mask: 0x%016X, magic: 0x%016X, shift: %d, offset: %d },\n",
			entry.mask, entry.magic, entry.shift, entry.offset)

		offset += uint32(entry.size)
	}
	fmt.Printf("Total Rook Table Size: %d\n\n", offset)
	rookOffset := offset

	fmt.Println("// Bishop Magics")
	offset = 0
	for sq := range 64 {
		entry := findMagicNumber(Square(sq), Bishop, rand)
		entry.offset = rookOffset + offset
		bishopMagics[sq] = entry

		fmt.Printf("{ mask: 0x%016X, magic: 0x%016X, shift: %d, offset: %d },\n",
			entry.mask, entry.magic, entry.shift, entry.offset)

		offset += uint32(entry.size)
	}
	fmt.Printf("Total Bishop Table Size: %d\n", offset)
}
