package engine

import (
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

const (
	FenNumberOfParts      = 6
	Splitter         rune = '/'
	Dash             rune = '-'
	EmDash           rune = '–'
	Space            rune = ' '
	FenStartPosition      = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	MaxMoveRule           = 100
)

type GameState struct {
	SideToMove      int
	CastlingRights  int
	HalfMoveClock   uint8
	EnPassantSquare Square
	FullmoveNumber  uint16
}

func NewGameState() *GameState {
	return &GameState{
		SideToMove:      white,
		CastlingRights:  0,
		HalfMoveClock:   0,
		EnPassantSquare: none,
		FullmoveNumber:  1,
	}
}

const MaxGameMoves = 2048
const EnPassantSquaresWhiteStart Square = 40
const EnPassantSquaresWhiteEnd Square = 47
const EnPassantSquaresBlackStart Square = 16
const EnPassantSquaresBlackEnd Square = 23

type BoardSetup interface {
	FenSetup(fen string) error
}

type SquareAttacked interface {
	isSquareAttacked(sq Square, side int) bool
}

type PrintBoard interface {
	PrintBoardWithPieces()
	GetPieceOnSquare() int
	getCastlingRightsString() string
}

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

type Board struct {
	Bitboards          [12]uint64
	OccupancyBitboards [3]uint64
	State              *GameState
	History            *History
	PieceList          [64]Square
}

func NewBoard() *Board {
	return &Board{
		Bitboards:          [12]uint64{},
		OccupancyBitboards: [3]uint64{},
		State:              NewGameState(),
		History:            NewHistory(),
	}
}

func (f *Fen) pieces() error {
	var rank int = 0
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
			rank++
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
					f.board.State.CastlingRights |= int(piece)
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
		if (sq >= int(EnPassantSquaresWhiteStart) && sq <= int(EnPassantSquaresWhiteEnd)) ||
			(sq >= int(EnPassantSquaresBlackStart) && sq <= int(EnPassantSquaresBlackEnd)) {

			println(sq)
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
	part := f.fenParts[f.currentPartIdx]
	if l := len(part); l >= 1 && l <= 3 {
		x, err := strconv.ParseUint(part, 10, 8)
		if err == nil && x <= MaxMoveRule {
			f.board.State.HalfMoveClock = uint8(x)
			return nil
		}
	}
	return &FenParseError{
		Part:    5,
		Message: "Error parsing fifth part of FEN string.",
	}

}

func (f *Fen) FullMoveNumber() error {
	part := f.fenParts[f.currentPartIdx]
	if l := len(part); l >= 1 && l <= 4 {

		x, err := strconv.ParseUint(part, 10, 16)

		if err == nil && x <= MaxGameMoves {
			f.board.State.FullmoveNumber = uint16(x)
			return nil
		}
	}

	return &FenParseError{
		Part:    6,
		Message: "Error parsing sixth part of FEN string.",
	}
}

func (board *Board) FenSetup(fen string) error {
	fenParts, err := splitFenString(fen)
	if err != nil {
		return err
	}
	tempBoard := NewBoard()
	f := NewFen(tempBoard)
	f.fenParts = fenParts

	fenSteps := []func() error{
		f.pieces,
		f.sideToMove,
		f.castling,
		f.EnPassant,
		f.HalfMoveClock,
		f.FullMoveNumber,
	}

	for i, step := range fenSteps {
		f.currentPartIdx = i
		if err := step(); err != nil {
			return err
		}
	}

	for piece := P; piece <= K; piece++ {
		tempBoard.OccupancyBitboards[white] |= f.board.Bitboards[piece]
	}
	
	for piece := p; piece <= k; piece++ {
		tempBoard.OccupancyBitboards[black] |= f.board.Bitboards[piece]
	}
	
	tempBoard.OccupancyBitboards[both] |= tempBoard.OccupancyBitboards[white] | tempBoard.OccupancyBitboards[black]

	*board = *tempBoard
	return nil
}

// / This function splits the incoming FEN-string into its component parts.
// It replaces the sometimes (mistakenly) used "em-dash" with the normal dash.
// Also, if the FEN-string turns out to be 4 parts long, the values 0 and 1 are assumed for the last two parts.
func splitFenString(fenString string) ([]string, error) {
	const ShortFenLength uint = 4

	if len(fenString) == 0 {
		fenString = FenStartPosition
	}

	fenString = strings.TrimSpace(fenString)
	fenString = strings.Replace(fenString, string(EmDash), string(Dash), 1)
	fenAsSlice := strings.Split(fenString, string(Space))

	if len(fenAsSlice) == int(ShortFenLength) {
		fenAsSlice = append(fenAsSlice, "0", "1")
	}

	if len(fenAsSlice) != FenNumberOfParts {
		return nil, fmt.Errorf("Invalid FEN string: %s", fenString)
	}

	println(fenString)
	fmt.Printf("%#v\n", fenAsSlice)
	return fenAsSlice, nil

}

func (board *Board) PrintBoardWithPieces() {
	fmt.Println()
	for rank := range 8 {
		for file := range 8 {

			var square Square = Square(rank*8 + file)

			if file == 0 {
				fmt.Printf("  %d ", 8-rank)
			}

			piece := board.PieceAt(square)
			if piece == Zilch {
				fmt.Printf(" .")
				continue
			}
			if runtime.GOOS == "windows" {
				fmt.Printf(" %c", asciiPieces[piece])
			} else {
				fmt.Printf(" %c", unicodePieces[piece])
			}
		}
		fmt.Println()
	}

	// print board files
	fmt.Println("\n     a b c d e f g h")

	fmt.Println()

	if board.State.SideToMove == white {
		fmt.Println("Side to move: White")
	} else {
		fmt.Println("Side to move: Black")
	}

	if board.State.EnPassantSquare != none {
		fmt.Printf("En Passant Square: %s\n", BitboardSquares[board.State.EnPassantSquare])
	} else {
		fmt.Printf("En Passant Square: None\n")
	}

	fmt.Printf("Castling Rights: %s\n", board.getCastlingRightsString())

}

func (board *Board) PieceAt(square Square) Piece {
	for piece := range k + 1 {
		if getBit(board.Bitboards[piece], square) == 1 {
			return piece
		}
	}
	return Zilch
}

func (board *Board) getCastlingRightsString() string {
	var rights string
	if board.State.CastlingRights&whiteKingside != 0 {
		rights += "K"
	}
	if board.State.CastlingRights&whiteQueenside != 0 {
		rights += "Q"
	}
	if board.State.CastlingRights&blackKingside != 0 {
		rights += "k"
	}
	if board.State.CastlingRights&blackQueenside != 0 {
		rights += "q"
	}
	return rights
}

// checks occupancy bitb oard of both sides to verify if
// the given square is occupied, or na-da.
func (board *Board) Occupied (square Square) bool {
	return getBit(board.OccupancyBitboards[both], square) != 0
}

// checks occupancy bitb oard of opp side to verify if
// the given square is occupied, or na-da.
func (board *Board) OccupiedByOpp (square Square, side int) bool {
	if side > 2 {
		panic("grrr grrrrr")
	}
	return getBit(board.OccupancyBitboards[side^1], square) != 0
}
