package engine

// / Represents a chess move containing the from and to squares, as well as flags for special moves.
// / The information encoded as a 16-bit integer, 6 bits for the from/to square and 4 bits for the flags.
// /
// / See [Encoding Moves](https://www.chessprogramming.org/Encoding_Moves) for more information.
type Move uint16

type MoveFlag uint8

const (
	Normal         MoveFlag = 0b0000
	DoublePawnPush MoveFlag = 0b0001
	KingCastle     MoveFlag = 0b0010
	QueenCastle    MoveFlag = 0b0011

	Capture   MoveFlag = 0b0100
	EnPassant MoveFlag = 0b0101

	PromotionN MoveFlag = 0b1000
	PromotionB MoveFlag = 0b1001
	PromotionR MoveFlag = 0b1010
	PromotionQ MoveFlag = 0b1011

	PromotionCaptureN MoveFlag = 0b1100
	PromotionCaptureB MoveFlag = 0b1101
	PromotionCaptureR MoveFlag = 0b1110
	PromotionCaptureQ MoveFlag = 0b1111
)

var promotedPieces = map[Piece]rune{
	Q: 'q',
	R: 'r',
	B: 'b',
	N: 'n',
	q: 'q',
	r: 'r',
	b: 'b',
	n: 'n',
}

func (f MoveFlag) IsCapture() bool {
	return f&Capture != 0
}

func (f MoveFlag) IsCastle() bool {
	return f == KingCastle || f == QueenCastle
}

func NewMove(from Square, to Square, flag MoveFlag) Move {
	// FLAG - TO - FROM....in that order
	return Move(uint16(from) | (uint16(to) << 6) | (uint16(flag) << 12))
}

func (move Move) DebugString(side int) string {
	from := BitboardSquares[move.getFrom()]
	to := BitboardSquares[move.getTo()]
	flag := move.getFlag()

	str := from + to

	if move.isPromotion() {
		str += pieceToChar(move.PromotionPiece(side))
	}

	if move.isCapture() {
		str += " (capture)"
	}

	if flag == EnPassant {
		str += " (ep)"
	}

	if flag == KingCastle {
		str += " (O-O)"
	}

	if flag == QueenCastle {
		str += " (O-O-O)"
	}

	return str
}

// taking last 6 bits: ergo, 0b0011_1111
func (move Move) getFrom() Square {
	return Square(move & 0x3f)
}

// shifting right by 6 bits to eliminate 'from' bits
// then masking only rightmost 6 bits: ergo, 0b0011_1111
func (move Move) getTo() Square {
	return Square((move >> 6) & 0x3f)
}

// shifting right by 12 bits to eliminate 'from' & 'to' bits
// then masking only rightmost 4 bits: ergo, 0b1111
func (move Move) getFlag() MoveFlag {
	return MoveFlag((move >> 12) & 0x0f) // we might as well not mask yunno
}

func (move Move) isPresent() bool {
	return !move.isNull()
}

func (move Move) isNull() bool {
	return move == 0
}

func (move Move) isNoisy() bool {
	flag := move.getFlag()
	return flag == Capture ||
		flag == EnPassant ||
		flag == PromotionQ ||
		flag == PromotionCaptureN ||
		flag == PromotionCaptureB ||
		flag == PromotionCaptureQ ||
		flag == PromotionCaptureR
}

func (move Move) isQuiet() bool {
	return move.isPresent() && !move.isNoisy()
}

// right shift move 16 bit integer first by 6, invariably getting
// rid of 'to' bits, then another 6 for 'from' bits, then 2 LSB
// from 'flag' bits, so we're left with 2 MSB from 'flag'.
// only capture & en-passant captures satisfy curr value & 1.
func (move Move) isCapture() bool {
	return (move>>14)&1 != 0
}

// same as isCapture, but shifting 'flag' by 3 rather than 2 and
// leaving just the MSB. since all promotions have MSB set (ergo 1)
// it's pretty straightforward.
func (move Move) isPromotion() bool {
	return (move >> 15) != 0
}

func (move Move) isEnPassant() bool {
	return move.getFlag() == EnPassant
}

func (move Move) isCastling() bool {
	return ((move.getFlag() == KingCastle) || (move.getFlag() == QueenCastle))
}

func (move Move) isDoublePush() bool {
	return move.getFlag() == DoublePawnPush
}

func (move Move) PromotionPiece(side int) Piece {
	if !move.isPromotion() {
		return Zilch
	}
	promoType := Piece(move.getFlag()&0b11) + N // +N offset maps 0->N,1->B,2->R,3->Q
	if side == black {
		promoType += 6 // offset for black pieces (p..k)
	}
	return promoType
}

func (move Move) getEncoded() uint {
	return uint(move & 0b0000_1111_1111_1111)
}

// func (board *Board) GenerateMoves() {
// 	var bitboard, attacks uint64
// 	var from Square

// 	for piece := P; piece < k; piece++ {
// 		bitboard = board.Bitboards[piece]
// 		from = Square(getIndexOfLS1B(bitboard))
// 		for bitboard != 0 {

// 			if board.State.SideToMove == white {
// 				if piece == P {
// 					generateQuietPawnMoves(board, bitboard, from, -8, a7, h7, a2, h2)
// 				}
// 			} else {
// 				if piece == p {
// 					generateQuietPawnMoves(board, bitboard, from, 8, a2, h2, a7, h7)
// 				}
// 			}
// 			bitboard = popBit(bitboard, from)

// 		}
// 		println(attacks)

// 	}
// }

func generatePawnCaptures(
	board *Board,
	bitboard uint64,
	from, to Square,
	attacks uint64,
	direction int,
	promotionRankStart, promotionRankEnd Square,

) {
	// to := int(from) + direction

	attacks = pawnAttacks[board.State.SideToMove][from] & board.OccupancyBitboards[black]

	for attacks != 0 {
		targetSquare := getIndexOfLS1B(attacks)
		if from >= promotionRankStart && from <= promotionRankEnd {
			println("PAWN capture promotion:", BitboardSquares[from], BitboardSquares[targetSquare])
		} else {
			println("PAWN capture:", BitboardSquares[from], BitboardSquares[targetSquare])
		}
		attacks = popBit(attacks, Square(targetSquare))
	}
}

func generateQuietPawnMoves(
	board *Board,
	bitboard uint64,
	from Square,
	direction int,
	promotionRankStart, promotionRankEnd Square,
	doublePushRankStart, doublePushRankEnd Square,
	// moves *[]Move,
) {
	to := int(from) + direction

	if to < int(a8) || to > int(h1) {
		bitboard = popBit(bitboard, from)
		return
	}

	// square gotta be empty
	if getBit(board.OccupancyBitboards[both], Square(to)) != 0 {
		bitboard = popBit(bitboard, from)
		return
	}

	// promotion
	if from >= promotionRankStart && from <= promotionRankEnd {
		println("PAWN promotion:", BitboardSquares[from], BitboardSquares[to])
	} else {
		println("PAWN push:", BitboardSquares[from], BitboardSquares[to])

		// double push
		doubleTo := to + direction
		if from >= doublePushRankStart && from <= doublePushRankEnd &&
			getBit(board.OccupancyBitboards[both], Square(doubleTo)) == 0 {

			println("2 step PAWN push:", BitboardSquares[from], BitboardSquares[doubleTo])
		}
	}

}
