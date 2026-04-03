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
	CapturedPiece   Piece
}

func NewGameState() *GameState {
	return &GameState{
		SideToMove:      white,
		CastlingRights:  0,
		HalfMoveClock:   0,
		EnPassantSquare: none,
		FullmoveNumber:  1,
		CapturedPiece:   Zilch,
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

type Undo struct {
	State              GameState
	Bitboards          [12]uint64
	OccupancyBitboards [3]uint64
	PieceList          [64]Piece
}

type History struct {
	list  [MaxGameMoves]Undo
	count uint
}

type IHistory interface {
	push(state Undo)
	pop() (Undo, bool, error)
	peek() Undo
	clear()
	isEmpty() bool
	isFull() bool
	size() uint
	getRef(index uint) *Undo
}

func NewHistory() *History {
	return &History{}
}

func (h *History) push(state Undo) {
	h.list[h.count] = state
	h.count++
}

// Pops the last undo snapshot from the history and returns it. If the history is empty, it returns an error.
func (h *History) pop() (Undo, bool, error) {
	if h.count > 0 {
		h.count--
		return h.list[h.count], true, nil
	}
	return Undo{}, false, fmt.Errorf("History is empty")
}

func (h *History) peek() Undo {
	return h.list[h.count-1]
}

func (h *History) clear() {
	h.count = 0
}

func (h *History) size() uint {
	return h.count
}

func (h *History) getRef(index uint) *Undo {
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
	PieceList          [64]Piece
}

func NewBoard() *Board {
	b := &Board{
		Bitboards:          [12]uint64{},
		OccupancyBitboards: [3]uint64{},
		State:              NewGameState(),
		History:            NewHistory(),
	}
	for i := range 64 {
		b.PieceList[i] = Zilch
	}
	return b
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
				f.board.PieceList[square] = piece
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
			switch char {
			case 'K':
				f.board.State.CastlingRights |= whiteKingside
			case 'Q':
				f.board.State.CastlingRights |= whiteQueenside
			case 'k':
				f.board.State.CastlingRights |= blackKingside
			case 'q':
				f.board.State.CastlingRights |= blackQueenside
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
	return board.PieceList[square]
}

func pieceSide(piece Piece) int {
	switch {
	case piece <= K:
		return white
	case piece <= k:
		return black
	default:
		panic(fmt.Sprintf("invalid piece: %d", piece))
	}
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
func (board *Board) Occupied(square Square) bool {
	return getBit(board.OccupancyBitboards[both], square) != 0
}

// checks occupancy bitb oard of opp side to verify if
// the given square is occupied, or na-da.
func (board *Board) OccupiedByOpp(square Square, side int) bool {
	if side > 2 {
		panic("grrr grrrrr")
	}
	return getBit(board.OccupancyBitboards[side^1], square) != 0
}

func (board *Board) snapshot() Undo {
	return Undo{
		State:              *board.State,
		Bitboards:          board.Bitboards,
		OccupancyBitboards: board.OccupancyBitboards,
		PieceList:          board.PieceList,
	}
}

func (board *Board) UndoMove(move Move) {
	state, worked, err := board.History.pop()
	if !worked || err != nil {
		panic(fmt.Sprintf("failed to pop undo state from history: %v", err))
	}

	board.Bitboards = state.Bitboards
	board.OccupancyBitboards = state.OccupancyBitboards
	board.PieceList = state.PieceList
	*board.State = state.State
}

func (board *Board) MakeMove(move Move) {
	sideToMove := board.State.SideToMove
	from, to := move.getFrom(), move.getTo()
	piece := board.PieceAt(from)

	board.History.push(board.snapshot())
	board.State.CapturedPiece = Zilch

	if board.State.EnPassantSquare != none {
		board.State.EnPassantSquare = none
	}

	if move.isCapture() || pieceToChar(piece) == "p" {
		board.State.HalfMoveClock = 0
	} else {
		board.State.HalfMoveClock += 1
	}

	captured := board.PieceAt(to)

	if captured != Zilch && !move.isCastling() {
		board.removePiece(piece, from)

		board.removePiece(captured, to)

		board.addPiece(piece, to)

		board.State.CapturedPiece = captured

	} else if !move.isCastling() {
		board.removePiece(piece, from)
		board.addPiece(piece, to)
	}

	switch {
	case move.isDoublePush():
		board.State.EnPassantSquare = Square((int(from) + int(to)) / 2)
	case move.isEnPassant():
		var capturedPiece Piece

		if sideToMove == white {
			capturedPiece = p
		} else {
			capturedPiece = P
		}
		board.removePiece(capturedPiece, to^8)
	case move.isCastling():
		rookFrom, rookTo := getCastlingRookMove(to)
		var rookPiece Piece
		if sideToMove == white {
			rookPiece = R
		} else {
			rookPiece = r
		}

		// remove the rook from its original square
		board.removePiece(rookPiece, rookFrom)
		// board.removePiece(board.PieceAt(rookFrom), rookFrom)

		// remove the king
		board.removePiece(piece, from)

		// add the king to its new square
		board.addPiece(piece, to)
		board.addPiece(rookPiece, rookTo)

	case move.isPromotion():
		var promoPiece Piece = move.PromotionPiece(sideToMove)
		board.removePiece(piece, to)
		board.addPiece(promoPiece, to)
	}

	board.State.SideToMove ^= 1

	if board.State.SideToMove == white {
		board.State.FullmoveNumber += 1
	}

	board.State.CastlingRights &^= castlingRightMask(from) | castlingRightMask(to)
}

func (board *Board) IsLegal(move Move) bool {
	if move.isNull() {
		panic("Move cannot be null.")
	}

	board.MakeMove(move)
	isLegal := !board.IsInCheck(board.State.SideToMove ^ 1)
	board.UndoMove(move)

	return isLegal
}

func (board *Board) isInMultipleCheck(side int) bool {
	kingSquare := board.getKingSquare(side)
	attackers := board.getAttackersToSquare(kingSquare, side^1)
	return attackers != 0 && (attackers&(attackers-1)) != 0
}

func (board *Board) IsInCheck(side int) bool {
	kingSquare := board.getKingSquare(side)
	return board.isSquareAttacked(kingSquare, side^1)
}

func (board *Board) getKingSquare(side int) Square {
	switch side {
	case white:
		return Square(getIndexOfLS1B(board.Bitboards[K]))
	case black:
		return Square(getIndexOfLS1B(board.Bitboards[k]))
	}
	return none
}

func castlingRightMask(square Square) int {
	switch square {
	case e1:
		return whiteKingside | whiteQueenside
	case a1:
		return whiteQueenside
	case h1:
		return whiteKingside
	case e8:
		return blackKingside | blackQueenside
	case a8:
		return blackQueenside
	case h8:
		return blackKingside
	default:
		return 0
	}
}

// returns the rook's starting and ending squares for a given castling move, based on the king's destination square.
func getCastlingRookMove(kingTo Square) (Square, Square) {
	switch kingTo {
	case g1:
		return h1, f1
	case c1:
		return a1, d1
	case g8:
		return h8, f8
	case c8:
		return a8, d8
	default:
		panic(fmt.Sprintf("Invalid castling move: kingTo=%d", kingTo))
	}
}

func (board *Board) addPiece(piece Piece, square Square) {
	side := pieceSide(piece)

	board.Bitboards[piece] = setBit(board.Bitboards[piece], square)
	board.OccupancyBitboards[side] = setBit(board.OccupancyBitboards[side], square)
	board.OccupancyBitboards[both] = board.OccupancyBitboards[white] | board.OccupancyBitboards[black]
	board.PieceList[square] = piece
}

func (board *Board) removePiece(piece Piece, square Square) {
	side := pieceSide(piece)

	board.Bitboards[piece] = popBit(board.Bitboards[piece], square)
	board.OccupancyBitboards[side] = popBit(board.OccupancyBitboards[side], square)
	board.OccupancyBitboards[both] = board.OccupancyBitboards[white] | board.OccupancyBitboards[black]
	board.PieceList[square] = Zilch
}
