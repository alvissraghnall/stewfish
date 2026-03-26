package internal

type ShadowArray[T any] struct {
	data [256]T
	count  int
}

func NewShadowArray[T any]() ShadowArray[T] {
	return ShadowArray[T]{}
}

func (s *ShadowArray[T]) Push(value T) {
	if s.count >= 256 {
		panic("Array has max length 256!")
	}
	s.data[s.count] = value
	s.count++
}

func (s *ShadowArray[T]) Get(i int) T {
	return s.data[i]
}

func (s *ShadowArray[T]) Len() int {
	return s.count
}

func (s *ShadowArray[T]) IsEmpty() bool {
	return s.count == 0
}

func (s *ShadowArray[T]) Clear() {
	s.count = 0
}

func (s *ShadowArray[T]) Slice() []T {
	return s.data[:s.count]
}