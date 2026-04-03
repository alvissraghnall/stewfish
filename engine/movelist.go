package engine

import (
	"github.com/alvissraghnall/stewfish/internal"
)

type MoveEntry struct {
	score int
	move  Move
}

type MoveList struct {
	data internal.ShadowArray[MoveEntry]
}

func NewMoveList() *MoveList {
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

func (list *MoveList) SliceMoves() []Move {
	moves := make([]Move, list.data.Len())
	for i, entry := range list.data.Slice() {
		moves[i] = entry.move
	}
	return moves
}

func (list *MoveList) Reset() {
	list.data.Clear()
}

func (list *MoveList) Get(index int) MoveEntry {
	return list.data.Get(index)
}

func (list *MoveList) GetMove (index int) Move {
	entry := list.data.Get(index)
	return entry.move
}

func (ml *MoveList) Print(side int, board *Board) {
	var blackMoves, whiteMoves int
	for _, move := range ml.Slice() {
		println(move.move.DebugString(side, board))
		if board.PieceAt(move.move.getFrom()) <= K {
			whiteMoves += 1
		} else if board.PieceAt(move.move.getFrom()) > K || board.PieceAt(move.move.getFrom()) <= k {
			blackMoves += 1
		}
	}
	println("Total moves: ", ml.Len())
	println("White moves: ", whiteMoves)
	println("Black moves: ", blackMoves)
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


func (list *MoveList) PushSetwiseNoFlag(from Square, toBB uint64) {
	for toBB != 0 {
		to := Square(getIndexOfLS1B(toBB))
		flag := Normal
		// If there are any pieces on the target square, it's a capture
		if list.data.Len() > 0 {
			flag = Capture
		}
		list.Add(from, to, flag)
		toBB = popBit(toBB, to)
	}
}

func (list *MoveList) PushSetwiseFlag(from Square, toBB uint64, flagFn func(to Square) MoveFlag) {
	for toBB != 0 {
		to := Square(getIndexOfLS1B(toBB))
		list.Add(from, to, flagFn(to))
		toBB = popBit(toBB, to)
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
