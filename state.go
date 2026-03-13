package main

import (
	"fmt"
	"slices"
	"strings"
)

const (
	FenNumberOfParts      = 6
	Splitter         rune = '/'
	Dash             rune = '-'
	EmDash           rune = '–'
	Space            rune = ' '
	FenStartPosition      = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

type GameState struct {
	SideToMove      int
	CastlingRights  int
	HalfMoveClock   uint8
	EnPassantSquare Square
	FullmoveNumber  uint16
}

const MaxGameMoves = 2048

type Fen struct {
	board          *Board
	fenParts       []string
	currentPartIdx int
}

func NewFen(board *Board) *Fen {
	return &Fen{
		board:          board,
		fenParts:       make([]string, FenNumberOfParts),
		currentPartIdx: 0,
	}
}

type FenParser interface {
	pieces() error
	sideToMove() error
	castling() error
	EnPassant() error
	HalfMoveClock() error
	FullMoveNumber() error
}

type History struct {
	list  [MaxGameMoves]GameState
	count uint
}

type IHistory interface {
	push()
	pop() (GameState, bool, error)
	peek() GameState
	clear()
	isEmpty() bool
	isFull() bool
	size() uint
	getRef(index uint) *GameState
}

func NewHistory() *History {
	return &History{}
}

func (h *History) push(state GameState) {
	h.list[h.count] = state
	h.count++
}

func (h *History) pop() (GameState, bool, error) {
	if h.count > 0 {
		h.count--
		return h.list[h.count], true, nil
	}
	return GameState{}, false, fmt.Errorf("History is empty")
}

func (h *History) peek() GameState {
	return h.list[h.count-1]
}

func (h *History) clear() {
	h.count = 0
}

func (h *History) size() uint {
	return h.count
}

func (h *History) getRef(index uint) *GameState {
	return &h.list[index]
}

func (h *History) isEmpty() bool {
	return h.count == 0
}

func (h *History) isFull() bool {
	return h.count == MaxGameMoves
}

type FenParseError struct {
	Part    int
	Message string
}

func (e *FenParseError) Error() string {
	return fmt.Sprintf("Part %d: %s", e.Part, e.Message)
}

type BoardMethods interface {
	FenSetup(fen string) error
}

type Board struct {
	Bitboards          [12]uint64
	OccupancyBitboards [3]uint64
	State              GameState
	History            *History
	PieceList          [64]Square
}

func NewBoard() *Board {
	return &Board{
		Bitboards:          [12]uint64{},
		OccupancyBitboards: [3]uint64{},
		State:              GameState{},
		History:            NewHistory(),
	}
}

func (f *Fen) pieces() error {
	var rank int = 7
	var file int = 0

	for _, char := range f.fenParts[f.currentPartIdx] {
		switch {
		case char == Splitter:
			if file != 8 {
				return &FenParseError{
					Part:    1,
					Message: "Error parsing first part of FEN string.",
				}
			}
			rank--
			file = 0
			continue

		case char >= '1' && char <= '8':
			file += int(char - '0')
			continue

		default:
			piece, ok := CharPieceMap[char]
			if ok {
				square := Square(rank*8 + file)
				f.board.Bitboards[piece] = setBit(f.board.Bitboards[piece], square)
				file++
			}
		}
	}
	return nil
}

func (f *Fen) sideToMove() error {
	part := f.fenParts[f.currentPartIdx]

	if len(part) == 1 {
		switch part {
		case "w":
			f.board.State.SideToMove = white
		case "b":
			f.board.State.SideToMove = black
		default:
			return &FenParseError{
				Part:    2,
				Message: "Error parsing seconf part of FEN string.",
			}
		}
		return nil
	}
	return &FenParseError{
		Part:    2,
		Message: "Error parsing seconf part of FEN string.",
	}

}

// There should be 1 to 4 castling rights. If no player has castling
// rights, the character is '-'.
func (f *Fen) castling() error {
	part := f.fenParts[f.currentPartIdx]
	if len(part) <= 4 && len(part) > 0 {
		for _, char := range part {
			if char == 'K' || char == 'Q' || char == 'k' || char == 'q' {
				piece, ok := CharPieceMap[char]
				if ok {
					f.board.State.CastlingRights |= piece
				}

			}
		}
		return nil
	}
	return &FenParseError{
		Part:    3,
		Message: "Error parsing third part of FEN string.",
	}

}

func (f *Fen) EnPassant() error {
	part := f.fenParts[f.currentPartIdx]
	if len(part) == 1 && part == string(Dash) {
		return nil
	}
	if len(part) == 2 {
		sq := slices.Index(BitboardSquares[:], part)
		if sq >= 0 {
			f.board.State.EnPassantSquare = Square(sq)
			return nil
		}
	}
	return &FenParseError{
		Part:    4,
		Message: "Error parsing fourth part of FEN string.",
	}
}

func (f *Fen) HalfMoveClock() error {

	return &FenParseError{
		Part:    5,
		Message: "Error parsing fifth part of FEN string.",
	}

}

func (f *Fen) FullMoveNumber() error {

	return &FenParseError{
		Part:    5,
		Message: "Error parsing fifth part of FEN string.",
	}
}

// / This function splits the incoming FEN-string into its component parts. It does a bit of error handling, such as replacing the sometimes (mistakenly) used "em-dash" with the normal dash. Also, if the FEN-string turns out to be 4 parts long, the values 0 and 1 are assumed for the last two parts.
func splitFenString(fenString string) ([]string, error) {
	const ShortFenLength uint = 4

	if len(fenString) == 0 {
		fenString = FenStartPosition
	}

	fenString = strings.Replace(fenString, string(EmDash), string(Dash), 1)
	fenAsSlice := strings.Split(fenString, string(Space))

	if len(fenAsSlice) == int(ShortFenLength) {
		fenAsSlice = append(fenAsSlice, "0", "1")
	}

	if len(fenAsSlice) != FenNumberOfParts {
		return nil, fmt.Errorf("Invalid FEN string: %s", fenString)
	}

	return fenAsSlice, nil

}
