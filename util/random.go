package util

type Random struct {
	state uint64
}

type XORShift interface {
	// Generates a random 64-bit number using xorShift64 algo
	XorShift64() uint64
}

func NewRandom (state uint64) *Random {
	return &Random{
		state,
	}
}

func (rand *Random) XorShift64() uint64 {
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
