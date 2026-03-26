package engine

import "github.com/alvissraghnall/stewfish/internal"

type MoveEntry struct {
	score int
	move  Move
}

type MoveList struct {
	data internal.ShadowArray[MoveEntry]
}

func newMoveList() *MoveList {
	return &MoveList{
		data: internal.NewShadowArray[MoveEntry](),
	}
}

func (list MoveList) isEmpty() bool {
	return list.data.IsEmpty()
}

func (list *MoveList) Add(from Square, to Square, flag MoveFlag) {
	list.data.Push(MoveEntry{
		move:  NewMove(from, to, flag),
		score: 0,
	})
}

func (list *MoveList) AddMove(move Move) {
	list.data.Push(MoveEntry{
		score: 0,
		move:  move,
	})
}

func (list *MoveList) Slice() []MoveEntry {
	return list.data.Slice()
}

func (list *MoveList) Len() int {
	return list.data.Len()
}

func (list *MoveList) PushSetwise(from Square, toBB uint64, flag MoveFlag) {
	for toBB != 0 {
		to := getIndexOfLS1B(toBB)
		list.Add(from, Square(to), flag)
		toBB = popBit(toBB, Square(to))
	}
}

func (list *MoveList) PushPawnsSetwise(offset int, toBB uint64, flag MoveFlag) {
	for toBB != 0 {
		to := getIndexOfLS1B(toBB)
		from := to - offset
		list.Add(Square(from), Square(to), flag)
		toBB = popBit(toBB, Square(to))
	}
}

func (list *MoveList) PushPromotionSetwise(offset int, toBB uint64) {
	if toBB != 0 {
		list.PushPawnsSetwise(offset, toBB, PromotionCaptureQ)
		list.PushPawnsSetwise(offset, toBB, PromotionCaptureR)
		list.PushPawnsSetwise(offset, toBB, PromotionCaptureB)
		list.PushPawnsSetwise(offset, toBB, PromotionCaptureN)
	}
}
