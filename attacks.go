package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"time"
)

const (
	NotFileA         uint64 = 18374403900871474942
	NotFileH         uint64 = 9187201950435737471
	NotFileGH        uint64 = 4557430888798830399
	NotFileAB        uint64 = 18229723555195321596
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
	both
)

var rng *rand.Rand
var xorShift65Rand XORShift

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

var bishopRelevantBits [64]int = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}
var rookRelevantBits [64]int = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

var rookMagics [64]MagicEntry = [64]MagicEntry{
	{mask: 0x000101010101017E, magic: 0x0A80004000108020, shift: 52, offset: 0},
	{mask: 0x000202020202027C, magic: 0x0140100020004000, shift: 53, offset: 4096},
	{mask: 0x000404040404047A, magic: 0x01000B1041012000, shift: 53, offset: 6144},
	{mask: 0x0008080808080876, magic: 0x0480044800100080, shift: 53, offset: 8192},
	{mask: 0x001010101010106E, magic: 0x0280024800040080, shift: 53, offset: 10240},
	{mask: 0x002020202020205E, magic: 0x81800102000C0080, shift: 53, offset: 12288},
	{mask: 0x004040404040403E, magic: 0x2100090000CA0004, shift: 53, offset: 14336},
	{mask: 0x008080808080807E, magic: 0x82000205008B402C, shift: 52, offset: 16384},
	{mask: 0x0001010101017E00, magic: 0x3011002080004101, shift: 53, offset: 20480},
	{mask: 0x0002020202027C00, magic: 0x2009402000401002, shift: 54, offset: 22528},
	{mask: 0x0004040404047A00, magic: 0x0084802000100080, shift: 54, offset: 23552},
	{mask: 0x0008080808087600, magic: 0x8049000900221000, shift: 54, offset: 24576},
	{mask: 0x0010101010106E00, magic: 0x0002000806002010, shift: 54, offset: 25600},
	{mask: 0x0020202020205E00, magic: 0x2052000200041008, shift: 54, offset: 26624},
	{mask: 0x0040404040403E00, magic: 0x2144000801028410, shift: 54, offset: 27648},
	{mask: 0x0080808080807E00, magic: 0x000200040081006A, shift: 53, offset: 28672},
	{mask: 0x00010101017E0100, magic: 0x5800808000400021, shift: 53, offset: 30720},
	{mask: 0x00020202027C0200, magic: 0x0020008020400080, shift: 54, offset: 32768},
	{mask: 0x00040404047A0400, magic: 0x0000808020001000, shift: 54, offset: 33792},
	{mask: 0x0008080808760800, magic: 0x4011010008100020, shift: 54, offset: 34816},
	{mask: 0x00101010106E1000, magic: 0xA859010008001004, shift: 54, offset: 35840},
	{mask: 0x00202020205E2000, magic: 0x0022008004008002, shift: 54, offset: 36864},
	{mask: 0x00404040403E4000, magic: 0x4C02A40030018822, shift: 54, offset: 37888},
	{mask: 0x00808080807E8000, magic: 0x0003020000408104, shift: 53, offset: 38912},
	{mask: 0x000101017E010100, magic: 0x0800209280004008, shift: 53, offset: 40960},
	{mask: 0x000202027C020200, magic: 0x00415000C0002004, shift: 54, offset: 43008},
	{mask: 0x000404047A040400, magic: 0x9210200100104100, shift: 54, offset: 44032},
	{mask: 0x0008080876080800, magic: 0x1009008900100020, shift: 54, offset: 45056},
	{mask: 0x001010106E101000, magic: 0x1408008080080400, shift: 54, offset: 46080},
	{mask: 0x002020205E202000, magic: 0x0000040080800200, shift: 54, offset: 47104},
	{mask: 0x004040403E404000, magic: 0x140A0002000811D4, shift: 54, offset: 48128},
	{mask: 0x008080807E808000, magic: 0x0018014200328104, shift: 53, offset: 49152},
	{mask: 0x0001017E01010100, magic: 0x0000400021800282, shift: 53, offset: 51200},
	{mask: 0x0002027C02020200, magic: 0x0000200444401000, shift: 54, offset: 53248},
	{mask: 0x0004047A04040400, magic: 0x2060200182801001, shift: 54, offset: 54272},
	{mask: 0x0008087608080800, magic: 0x8000801002800800, shift: 54, offset: 55296},
	{mask: 0x0010106E10101000, magic: 0x2890800800800402, shift: 54, offset: 56320},
	{mask: 0x0020205E20202000, magic: 0x0050800200800400, shift: 54, offset: 57344},
	{mask: 0x0040403E40404000, magic: 0x0000225904001008, shift: 54, offset: 58368},
	{mask: 0x0080807E80808000, magic: 0x2003000841002482, shift: 53, offset: 59392},
	{mask: 0x00017E0101010100, magic: 0x80A0804000208004, shift: 53, offset: 61440},
	{mask: 0x00027C0202020200, magic: 0x0104200050024001, shift: 54, offset: 63488},
	{mask: 0x00047A0404040400, magic: 0x1510102001010040, shift: 54, offset: 64512},
	{mask: 0x0008760808080800, magic: 0x8820120040220009, shift: 54, offset: 65536},
	{mask: 0x00106E1010101000, magic: 0x01C101180091001C, shift: 54, offset: 66560},
	{mask: 0x00205E2020202000, magic: 0x4400020004008080, shift: 54, offset: 67584},
	{mask: 0x00403E4040404000, magic: 0x0A80100102040008, shift: 54, offset: 68608},
	{mask: 0x00807E8080808000, magic: 0x0001001088410002, shift: 53, offset: 69632},
	{mask: 0x007E010101010100, magic: 0xC040002050800080, shift: 53, offset: 71680},
	{mask: 0x007C020202020200, magic: 0x0400842000400C80, shift: 54, offset: 73728},
	{mask: 0x007A040404040400, magic: 0x0000430020041500, shift: 54, offset: 74752},
	{mask: 0x0076080808080800, magic: 0x0648080010008080, shift: 54, offset: 75776},
	{mask: 0x006E101010101000, magic: 0x0014040080080080, shift: 54, offset: 76800},
	{mask: 0x005E202020202000, magic: 0x0005800200040080, shift: 54, offset: 77824},
	{mask: 0x003E404040404000, magic: 0x0C01024801101400, shift: 54, offset: 78848},
	{mask: 0x007E808080808000, magic: 0x0401000082004100, shift: 53, offset: 79872},
	{mask: 0x7E01010101010100, magic: 0x1100C81080042101, shift: 52, offset: 81920},
	{mask: 0x7C02020202020200, magic: 0x0446001023008042, shift: 53, offset: 86016},
	{mask: 0x7A04040404040400, magic: 0x001200F218C06082, shift: 53, offset: 88064},
	{mask: 0x7608080808080800, magic: 0x4044090004201001, shift: 53, offset: 90112},
	{mask: 0x6E10101010101000, magic: 0x0802010410082002, shift: 53, offset: 92160},
	{mask: 0x5E20202020202000, magic: 0x400100080C00020B, shift: 53, offset: 94208},
	{mask: 0x3E40404040404000, magic: 0x0A1010D211081004, shift: 53, offset: 96256},
	{mask: 0x7E80808080808000, magic: 0x4A88008041002402, shift: 52, offset: 98304},
}

const RookMapSize uint = 102400

var bishopMagics [64]MagicEntry = [64]MagicEntry{
	{mask: 0x0040201008040200, magic: 0x0020099401828600, shift: 58, offset: 0},
	{mask: 0x0000402010080400, magic: 0x0110040800484201, shift: 59, offset: 64},
	{mask: 0x0000004020100A00, magic: 0x00210A2408400000, shift: 59, offset: 96},
	{mask: 0x0000000040221400, magic: 0x0004040888C08544, shift: 59, offset: 128},
	{mask: 0x0000000002442800, magic: 0x0082121040080000, shift: 59, offset: 160},
	{mask: 0x0000000204085000, magic: 0x0401012840380008, shift: 59, offset: 192},
	{mask: 0x0000020408102000, magic: 0x4E03081110890208, shift: 59, offset: 224},
	{mask: 0x0002040810204000, magic: 0x080A002101182104, shift: 58, offset: 256},
	{mask: 0x0020100804020000, magic: 0x04200A2004140040, shift: 59, offset: 320},
	{mask: 0x0040201008040000, magic: 0x082A041425940510, shift: 59, offset: 352},
	{mask: 0x00004020100A0000, magic: 0x0800042102120000, shift: 59, offset: 384},
	{mask: 0x0000004022140000, magic: 0x2012611041000400, shift: 59, offset: 416},
	{mask: 0x0000000244280000, magic: 0x4840040420100520, shift: 59, offset: 448},
	{mask: 0x0000020408500000, magic: 0x0030811002100200, shift: 59, offset: 480},
	{mask: 0x0002040810200000, magic: 0x2100040148080440, shift: 59, offset: 512},
	{mask: 0x0004081020400000, magic: 0x8000002888041004, shift: 59, offset: 544},
	{mask: 0x0010080402000200, magic: 0x015D001090108108, shift: 59, offset: 576},
	{mask: 0x0020100804000400, magic: 0x0143081004011400, shift: 59, offset: 608},
	{mask: 0x004020100A000A00, magic: 0xC008050C18012208, shift: 57, offset: 640},
	{mask: 0x0000402214001400, magic: 0x0028000082810180, shift: 57, offset: 768},
	{mask: 0x0000024428002800, magic: 0x2024200202010400, shift: 57, offset: 896},
	{mask: 0x0002040850005000, magic: 0x0A02000900490460, shift: 57, offset: 1024},
	{mask: 0x0004081020002000, magic: 0x0600414202422000, shift: 59, offset: 1152},
	{mask: 0x0008102040004000, magic: 0x0005009201110180, shift: 59, offset: 1184},
	{mask: 0x0008040200020400, magic: 0x02CB502068A01804, shift: 59, offset: 1216},
	{mask: 0x0010080400040800, magic: 0x0410101104041080, shift: 59, offset: 1248},
	{mask: 0x0020100A000A1000, magic: 0x000C441908002C00, shift: 57, offset: 1280},
	{mask: 0x0040221400142200, magic: 0x4814080010082008, shift: 55, offset: 1408},
	{mask: 0x0002442800284400, magic: 0x2109001091004004, shift: 55, offset: 1920},
	{mask: 0x0004085000500800, magic: 0x30008200810100A2, shift: 57, offset: 2432},
	{mask: 0x0008102000201000, magic: 0x000C20410D080290, shift: 59, offset: 2560},
	{mask: 0x0010204000402000, magic: 0x0408405005040200, shift: 59, offset: 2592},
	{mask: 0x0004020002040800, magic: 0x010410484104A080, shift: 59, offset: 2624},
	{mask: 0x0008040004081000, magic: 0x0024210400208400, shift: 59, offset: 2656},
	{mask: 0x00100A000A102000, magic: 0x0088404040080201, shift: 57, offset: 2688},
	{mask: 0x0022140014224000, magic: 0x4280020080080080, shift: 55, offset: 2816},
	{mask: 0x0044280028440200, magic: 0x8404010200140048, shift: 55, offset: 3328},
	{mask: 0x0008500050080400, magic: 0x48A1102080010040, shift: 57, offset: 3840},
	{mask: 0x0010200020100800, magic: 0x000124010000880E, shift: 59, offset: 3968},
	{mask: 0x0020400040201000, magic: 0x2001040118802100, shift: 59, offset: 4000},
	{mask: 0x0002000204081000, magic: 0x08A2421004044001, shift: 59, offset: 4032},
	{mask: 0x0004000408102000, magic: 0x0004021A10131209, shift: 59, offset: 4064},
	{mask: 0x000A000A10204000, magic: 0x2410404020821002, shift: 57, offset: 4096},
	{mask: 0x0014001422400000, magic: 0x0000002018014100, shift: 57, offset: 4224},
	{mask: 0x0028002844020000, magic: 0x4391408810401602, shift: 57, offset: 4352},
	{mask: 0x0050005008040200, magic: 0x0041090509025201, shift: 57, offset: 4480},
	{mask: 0x0020002010080400, magic: 0x2008010806A00211, shift: 59, offset: 4608},
	{mask: 0x0040004020100800, magic: 0x208128089C800100, shift: 59, offset: 4640},
	{mask: 0x0000020408102000, magic: 0x010484B410C01440, shift: 59, offset: 4672},
	{mask: 0x0000040810204000, magic: 0x0006010118824004, shift: 59, offset: 4704},
	{mask: 0x00000A1020400000, magic: 0x00000A0101210100, shift: 59, offset: 4736},
	{mask: 0x0000142240000000, magic: 0x0004182104880002, shift: 59, offset: 4768},
	{mask: 0x0000284402000000, magic: 0x00800C0420822000, shift: 59, offset: 4800},
	{mask: 0x0000500804020000, magic: 0x1220400204210000, shift: 59, offset: 4832},
	{mask: 0x0000201008040200, magic: 0xE005155002020060, shift: 59, offset: 4864},
	{mask: 0x0000402010080400, magic: 0x0008424802002000, shift: 59, offset: 4896},
	{mask: 0x0002040810204000, magic: 0x0000140202100400, shift: 58, offset: 4928},
	{mask: 0x0004081020400000, magic: 0x04002682080A4280, shift: 59, offset: 4992},
	{mask: 0x000A102040000000, magic: 0x24C0012100824100, shift: 59, offset: 5024},
	{mask: 0x0014224000000000, magic: 0xC050845502104410, shift: 59, offset: 5056},
	{mask: 0x0028440200000000, magic: 0x0000808A13020200, shift: 59, offset: 5088},
	{mask: 0x0050080402000000, magic: 0x0000118D60081640, shift: 59, offset: 5120},
	{mask: 0x0020100804020000, magic: 0x0028483304480202, shift: 59, offset: 5152},
	{mask: 0x0040201008040200, magic: 0x08C0102252819288, shift: 58, offset: 5184},
}

const BishopMapSize uint = 5248

var AttackTable []uint64 = make([]uint64, RookMapSize+BishopMapSize)

/**
VISUALIZATION:

Horizontal
East (->) = << 1
West (<-) = >> 1
South (DOWN) = << 8
North (UP) = >> 8
*/

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	xorShift65Rand = &Random{
		state: 0xFFAAB58C5833FE89,
	}
}

func maskPawnAttacks(side int, square Square) uint64 {
	var attacks uint64 = 0

	var bitboard uint64 = 0

	bitboard = setBit(bitboard, square)

	if side == 0 {
		if ((bitboard >> 7) & NotFileA) != 0 {
			attacks |= (bitboard >> 7)
		}
		if ((bitboard >> 9) & NotFileH) != 0 {
			attacks |= (bitboard >> 9)
		}
	} else {
		if ((bitboard << 7) & NotFileH) != 0 {
			attacks |= (bitboard << 7)
		}
		if ((bitboard << 9) & NotFileA) != 0 {
			attacks |= (bitboard << 9)
		}
	}

	return attacks
}

func maskKnightAttacks(square Square) uint64 {
	bb := uint64(1) << square

	return ((bb >> 17) & NotFileA) | // Up 2, Left 1
		((bb >> 15) & NotFileH) | // Up 2, Right 1
		((bb >> 10) & NotFileAB) | // Up 1, Left 2
		((bb >> 6) & NotFileGH) | // Up 1, Right 2
		((bb << 17) & NotFileH) | // Down 2, Right 1
		((bb << 15) & NotFileA) | // Down 2, Left 1
		((bb << 10) & NotFileGH) | // Down 1, Right 2
		((bb << 6) & NotFileAB) // Down 1, Left 2
}

func maskKingAttacks(square uint8) uint64 {
	king := uint64(1) << square

	// horizontal
	attacks := ((king & NotFileH) << 1) | // east
		((king & NotFileA) >> 1) // west

	span := king | attacks

	// vertical + diagonals
	attacks |= (span << 8) | (span >> 8)

	return attacks
}

func diagonalCompute(square Square) uint64 {
	diagonal := int(square&7) - int(square>>3)

	// println("square: ", square, " diag: ", diagonal)
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
	// println("square: ", square, " anti: ", diag, " file: ", file, " rank: ", rank)
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
	// PrintBitboard(reverse)
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
	// count := 0
	// for bb != 0 {
	// 	count++
	// 	bb &= bb - 1
	// }
	// return count
	return bits.OnesCount64(bb)
}

func getIndexOfLS1B(bb uint64) int {
	if bb == 0 {
		return -1
	}
	return bits.TrailingZeros64(bb)
}

// SetOccupancy generates occupancy variations for a given attack mask
// * `index`: the permutation index (0 to 2^bits_in_mask - 1)
// * `bitsInMask` : number of bits set in the attack mask
// * `attackMask`: the mask of squares that can be occupied
func SetOccupancy(index int, bitsInMask int, attackMask uint64) uint64 {
	occupancy := uint64(0)

	for count := 0; count < bitsInMask; count++ {
		// Get the index of the least significant 1st set bit (LSB)
		square := getIndexOfLS1B(attackMask)

		// If no more bits are set, break the loop
		if square == -1 {
			break
		}

		attackMask = popBit(attackMask, Square(square))

		// Ensure occupancy is on the board
		if index&(1<<count) != 0 {
			// Populate occupancy map
			occupancy |= 1 << square
		}
	}
	return occupancy
}

func distanceToEdge(square Square, direction int8) int8 {
	switch direction {
	case 1:
		return 7 - file(square)
	case 8:
		return 7 - rank(square)
	case -1:
		return file(square)
	case -8:
		return rank(square)

	case 7:
		return min(7-rank(square), file(square))
	case 9:
		return min(7-rank(square), 7-file(square))

	case -9:
		return min(rank(square), file(square))
	case -7:
		return min(rank(square), 7-file(square))
	default:
		panic("Unexpected direction: " + fmt.Sprint(direction))
	}
}

func rank(square Square) int8 {
	return int8(square) / 8
}

func file(square Square) int8 {
	return int8(square) & 7
}

func initSliderAttacks(square Square, piece SliderPiece, magic MagicEntry) {
	attackMask := magic.mask
	bitsInMask := popCount(attackMask)

	// Generate all possible occupancy variations for the attack mask
	for index := 0; index < (1 << bitsInMask); index++ {
		occupancy := SetOccupancy(index, bitsInMask, attackMask)

		var attacks uint64
		if piece == Rook {
			attacks = rookAttacks(square, occupancy)
		} else {
			attacks = bishopAttacks(square, occupancy)
		}

		magicIdx := (occupancy * magic.magic) >> magic.shift

		AttackTable[uint64(magic.offset)+magicIdx] = attacks
	}
}

func InitSliderTables() {
	for sq := range 64 {
		initSliderAttacks(Square(sq), Rook, rookMagics[sq])
		initSliderAttacks(Square(sq), Bishop, bishopMagics[sq])
	}

}

func GetAttacks(sq Square, occupancy uint64, magicEntry MagicEntry) uint64 {
	// mask the occupancy to only keep relevant squares
	// ensures bits outside the mask don't mess up the multiplication
	occupancy &= magicEntry.mask

	index := (occupancy * magicEntry.magic) >> magicEntry.shift

	return AttackTable[magicEntry.offset+uint32(index)]
}
