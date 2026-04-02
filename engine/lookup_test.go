package engine

import (
	"testing"
)

func TestTriangularIndex(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{0, 0, 0},
		{1, 0, 1},
		{1, 1, 64},
		{2, 0, 2},
		{2, 1, 65},
		{2, 2, 127},
		{3, 0, 3},
		{3, 1, 66},
		{3, 2, 128},
		{3, 3, 189},
		{63, 63, 2079},
		{0, 63, 63},
		{63, 0, 63},
	}

	for _, tt := range tests {
		result := triangularIndex(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("triangularIndex(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestTriangularIndexSymmetric(t *testing.T) {
	for a := range 64 {
		for b := range 64 {
			idx1 := triangularIndex(a, b)
			idx2 := triangularIndex(b, a)
			if idx1 != idx2 {
				t.Errorf("triangularIndex(%d, %d) = %d != triangularIndex(%d, %d) = %d", a, b, idx1, b, a, idx2)
			}
		}
	}
}

func TestTriangularIndexUnique(t *testing.T) {
	seen := make(map[int]bool)
	for a := 0; a < 64; a++ {
		for b := a; b < 64; b++ {
			idx := triangularIndex(a, b)
			if idx < 0 || idx >= 64*65/2 {
				t.Errorf("triangularIndex(%d, %d) = %d out of range", a, b, idx)
			}
			if seen[idx] {
				t.Errorf("duplicate index %d for (%d, %d)", idx, a, b)
			}
			seen[idx] = true
		}
	}
	if len(seen) != 64*65/2 {
		t.Errorf("expected %d unique indices, got %d", 64*65/2, len(seen))
	}
}

func TestInBetween(t *testing.T) {
	Init()

	cases := []struct {
		from, to Square
		expected uint64
	}{
		{from: a1, to: h1, expected: Rank8 & NotFileA & NotFileH}, // b1-g1
		{from: a1, to: a8, expected: FileA & ^Rank1 & ^Rank8}, // a2-a7
		{from: a1, to: h8, expected: AntiDiagonalMasked}, // b2-g7
		{from: a1, to: b2, expected: 0}, // next to each other so no squares in between
		{from: a1, to: c2, expected: 0}, // not aligned
		{from: d4, to: g4, expected: setBit(setBit(0, e4), f4)}, // e4-f4
		{from: d4, to: d7, expected: setBit(setBit(0, d5), d6)}, // d5-d6
		{from: d4, to: g7, expected: setBit(setBit(0, e5), f6)}, // e5-f6
		{from: a1, to: e5, expected: setBit(setBit(setBit(0, d4), c3), b2)}, // b2-d4
	}
	PrintBitboard(setBit(setBit(0, d4), b2))

	for _, c := range cases {
		if got := inBetween(c.from, c.to); got != c.expected {
			t.Errorf("inBetween(%v,%v) = 0x%016x; want 0x%016x", c.from, c.to, got, c.expected)
		}
		if got := inBetween(c.to, c.from); got != c.expected {
			t.Errorf("inBetween(%v,%v) = 0x%016x; want 0x%016x (symmetry)", c.to, c.from, got, c.expected)
		}
	}
}

func TestInBetweenDebugRookA8A1(t *testing.T) {
	Init()
	a8 := Square(0)
	a1 := Square(56)
	attA8Empty := GetAttacks(a8, 0, rookMagics[a8])
	t.Logf("a8 rook attack mask empty = 0x%016x", attA8Empty)
	if !checkBit(attA8Empty, a1) {
		t.Fatalf("a8 does not see a1 in empty occupancy; a1 bit: %v", checkBit(attA8Empty, a1))
	}
	atkA8 := GetAttacks(a8, 1<<uint64(a1), rookMagics[a8])
	atkA1 := GetAttacks(a1, 1<<uint64(a8), rookMagics[a1])
	intersect := atkA8 & atkA1
	want := FileA & ^Rank1 & ^Rank8
	if intersect != want {
		t.Fatalf("rook a8,a1 inbetween computed=0x%016x expect=0x%016x", intersect, want)
	}
	if bt := inBetween(a8,a1); bt != want {
		t.Fatalf("inBetween(a8,a1)=0x%016x expect=0x%016x", bt, want)
	}
}


