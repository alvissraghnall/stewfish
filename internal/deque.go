package internal

type Deque[T interface{}] struct {
	elements []T
}

// type Deque []T

func (d *Deque[T]) EnqueueRight(element T) {
	d.elements = append(d.elements, element)
}

func (d *Deque[T]) DequeueLeft() (T, bool) {
	if len(d.elements) == 0 {
		var zero T
		return zero, false
	}

	element := d.elements[0]
	var zero T
	d.elements[0] = zero
	d.elements = d.elements[1:]

	return element, true
}

func (d *Deque[T]) EnqueueLeft(element T) {
	d.elements = append([]T{element}, d.elements...)
}

func (d *Deque[T]) DequeueRight() (T, bool) {
	if len(d.elements) == 0 {
		var zero T
		return zero, false
	}
	element := d.elements[len(d.elements)-1]
	var zero T
	d.elements[len(d.elements)-1] = zero
	d.elements = d.elements[:len(d.elements)-1]
	return element, true
}

func (d *Deque[T]) Len() int {
	return len(d.elements)
}
