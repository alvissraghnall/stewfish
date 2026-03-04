package main

type Random struct {
	state uint64
}

type XORShift interface {
	// Generates a random 64-bit number using xorShift64 algo
	xorShift64() uint64
}

func (rand *Random) xorShift64() uint64 {
	x := rand.state
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17

	// x ^= x << 13
	// x ^= x >> 17
	// x ^= x << 5
	rand.state = x

	return x
}
