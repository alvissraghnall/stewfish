package engine

var between = make([]uint64, 64*65/2)

func triangularIndex(a, b int) int {
	if a < b {
		a, b = b, a
	}
	return a + ((b * (127 - b)) >> 1)
}

func initLookups() {
	for i := range 64 {
		for j := range 64 {
			a, b := Square(i), Square(j)
			var inBetween uint64 = 0

			if checkBit(GetAttacks(a, 0, rookMagics[a]), b) {
				inBetween = GetAttacks(a, 1<<uint64(b), rookMagics[a]) & GetAttacks(b, 1<<uint64(a), rookMagics[b])
			} else if checkBit(GetAttacks(a, 0, bishopMagics[a]), b) {
				inBetween = GetAttacks(a, 1<<uint64(b), bishopMagics[a]) & GetAttacks(b, 1<<uint64(a), bishopMagics[b])
			}
			idx := triangularIndex(i, j)
			between[idx] = inBetween
		}
	}
}

func inBetween(from Square, to Square) uint64 {
	return between[triangularIndex(int(from), int(to))]
}
