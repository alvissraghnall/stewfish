package engine

import "testing"

type boardSnapshot struct {
	bitboards          [12]uint64
	occupancyBitboards [3]uint64
	state              GameState
}

func snapshotBoard(board *Board) boardSnapshot {
	return boardSnapshot{
		bitboards:          board.Bitboards,
		occupancyBitboards: board.OccupancyBitboards,
		state:              *board.State,
	}
}

func assertBoardSnapshotEqual(t *testing.T, board *Board, want boardSnapshot) {
	t.Helper()

	if board.Bitboards != want.bitboards {
		t.Fatalf("bitboards = %#v, want %#v", board.Bitboards, want.bitboards)
	}
	if board.OccupancyBitboards != want.occupancyBitboards {
		t.Fatalf("occupancy bitboards = %#v, want %#v", board.OccupancyBitboards, want.occupancyBitboards)
	}
	if *board.State != want.state {
		t.Fatalf("game state = %#v, want %#v", *board.State, want.state)
	}
}

func TestAddPieceUpdatesCorrectOccupancy(t *testing.T) {
	board := NewBoard()

	board.addPiece(P, e4)

	whiteMask := uint64(1) << e4
	if board.Bitboards[P] != whiteMask {
		t.Fatalf("white pawn bitboard = %064b, want %064b", board.Bitboards[P], whiteMask)
	}
	if board.OccupancyBitboards[white] != whiteMask {
		t.Fatalf("white occupancy = %064b, want %064b", board.OccupancyBitboards[white], whiteMask)
	}
	if board.OccupancyBitboards[black] != 0 {
		t.Fatalf("black occupancy = %064b, want 0", board.OccupancyBitboards[black])
	}
	if board.OccupancyBitboards[both] != whiteMask {
		t.Fatalf("both occupancy = %064b, want %064b", board.OccupancyBitboards[both], whiteMask)
	}

	board.addPiece(q, d5)

	blackMask := uint64(1) << d5
	if board.OccupancyBitboards[white] != whiteMask {
		t.Fatalf("white occupancy after black add = %064b, want %064b", board.OccupancyBitboards[white], whiteMask)
	}
	if board.OccupancyBitboards[black] != blackMask {
		t.Fatalf("black occupancy = %064b, want %064b", board.OccupancyBitboards[black], blackMask)
	}
	if board.OccupancyBitboards[both] != whiteMask|blackMask {
		t.Fatalf("both occupancy = %064b, want %064b", board.OccupancyBitboards[both], whiteMask|blackMask)
	}
}

func TestRemovePieceClearsCorrectOccupancy(t *testing.T) {
	board := NewBoard()
	board.addPiece(P, e4)
	board.addPiece(q, d5)

	board.removePiece(P, e4)

	blackMask := uint64(1) << d5
	if board.Bitboards[P] != 0 {
		t.Fatalf("white pawn bitboard = %064b, want 0", board.Bitboards[P])
	}
	if board.OccupancyBitboards[white] != 0 {
		t.Fatalf("white occupancy = %064b, want 0", board.OccupancyBitboards[white])
	}
	if board.OccupancyBitboards[black] != blackMask {
		t.Fatalf("black occupancy = %064b, want %064b", board.OccupancyBitboards[black], blackMask)
	}
	if board.OccupancyBitboards[both] != blackMask {
		t.Fatalf("both occupancy = %064b, want %064b", board.OccupancyBitboards[both], blackMask)
	}
}

func TestMakeMoveDoublePushSetsEnPassantSquare(t *testing.T) {
	board := NewBoard()
	board.addPiece(P, e2)

	board.MakeMove(NewMove(e2, e4, DoublePawnPush))

	if board.PieceAt(e4) != P {
		t.Fatalf("piece at e4 = %v, want white pawn", board.PieceAt(e4))
	}
	if board.PieceAt(e2) != Zilch {
		t.Fatalf("piece at e2 = %v, want empty", board.PieceAt(e2))
	}
	if board.State.EnPassantSquare != e3 {
		t.Fatalf("en passant square = %v, want %v", board.State.EnPassantSquare, e3)
	}
	if board.State.SideToMove != black {
		t.Fatalf("side to move = %d, want %d", board.State.SideToMove, black)
	}
}

func TestMakeMoveEnPassantRemovesCapturedPawn(t *testing.T) {
	board := NewBoard()
	board.State.SideToMove = white
	board.State.EnPassantSquare = e6
	board.addPiece(P, d5)
	board.addPiece(p, e5)

	board.MakeMove(NewMove(d5, e6, EnPassant))

	if board.PieceAt(d5) != Zilch {
		t.Fatalf("piece at d5 = %v, want empty", board.PieceAt(d5))
	}
	if board.PieceAt(e5) != Zilch {
		t.Fatalf("piece at e5 = %v, want empty", board.PieceAt(e5))
	}
	if board.PieceAt(e6) != P {
		t.Fatalf("piece at e6 = %v, want white pawn", board.PieceAt(e6))
	}
	if board.State.EnPassantSquare != none {
		t.Fatalf("en passant square = %v, want none", board.State.EnPassantSquare)
	}
}

func TestMakeMoveCastlingMovesKingAndRook(t *testing.T) {
	board := NewBoard()
	board.State.SideToMove = white
	board.State.CastlingRights = whiteKingside | whiteQueenside
	board.addPiece(K, e1)
	board.addPiece(R, h1)

	board.MakeMove(NewMove(e1, g1, KingCastle))

	if board.PieceAt(e1) != Zilch {
		t.Fatalf("piece at e1 = %v, want empty", board.PieceAt(e1))
	}
	if board.PieceAt(g1) != K {
		t.Fatalf("piece at g1 = %v, want white king", board.PieceAt(g1))
	}
	if board.PieceAt(f1) != R {
		t.Fatalf("piece at f1 = %v, want white rook", board.PieceAt(f1))
	}
	if board.State.CastlingRights != 0 {
		t.Fatalf("castling rights = %d, want 0", board.State.CastlingRights)
	}
}

func TestMakeMovePromotionReplacesPawn(t *testing.T) {
	board := NewBoard()
	board.State.SideToMove = white
	board.addPiece(P, a7)

	board.MakeMove(NewMove(a7, a8, PromotionQ))

	if board.PieceAt(a7) != Zilch {
		t.Fatalf("piece at a7 = %v, want empty", board.PieceAt(a7))
	}
	if board.PieceAt(a8) != Q {
		t.Fatalf("piece at a8 = %v, want white queen", board.PieceAt(a8))
	}
	if board.Bitboards[P] != 0 {
		t.Fatalf("white pawn bitboard = %064b, want 0", board.Bitboards[P])
	}
}

func TestFenSetupParsesPartialCastlingRights(t *testing.T) {
	board := NewBoard()

	if err := board.FenSetup("4k3/8/8/8/8/8/8/4K3 w KQ - 0 1"); err != nil {
		t.Fatalf("FenSetup returned error: %v", err)
	}

	want := whiteKingside | whiteQueenside
	if board.State.CastlingRights != want {
		t.Fatalf("castling rights = %d, want %d", board.State.CastlingRights, want)
	}
}

func TestUndoMoveRestoresPositionAfterDoublePush(t *testing.T) {
	board := NewBoard()
	board.State.CastlingRights = whiteKingside | whiteQueenside
	board.addPiece(P, e2)

	before := snapshotBoard(board)
	move := NewMove(e2, e4, DoublePawnPush)

	board.MakeMove(move)
	board.UndoMove(move)

	assertBoardSnapshotEqual(t, board, before)
	if board.History.size() != 0 {
		t.Fatalf("history size = %d, want 0", board.History.size())
	}
}

func TestUndoMoveRestoresPositionAfterEnPassant(t *testing.T) {
	board := NewBoard()
	board.State.SideToMove = white
	board.State.EnPassantSquare = e6
	board.addPiece(P, d5)
	board.addPiece(p, e5)

	before := snapshotBoard(board)
	move := NewMove(d5, e6, EnPassant)

	board.MakeMove(move)
	board.UndoMove(move)

	assertBoardSnapshotEqual(t, board, before)
	if board.History.size() != 0 {
		t.Fatalf("history size = %d, want 0", board.History.size())
	}
}

func TestGetKingSquare (t *testing.T) {
	board := NewBoard()
	board.addPiece(K, e1)
	board.addPiece(k, d8)

	if board.getKingSquare(white) != e1 {
		println("White king square = ", board.getKingSquare(white))
		t.Fatalf("white king square = %v, want e1", board.getKingSquare(white))
	}
	if board.getKingSquare(black) != d8 {
		println("Black king square = ", board.getKingSquare(black))
		t.Fatalf("black king square = %v, want d8", board.getKingSquare(black))
	}
}

func TestGetKingSquareReturnsNoneIfNoKing(t *testing.T) {
	board := NewBoard()

	if board.getKingSquare(white) != none {
		println("White king square = ", board.getKingSquare(white))
		t.Fatalf("white king square = %v, want none", board.getKingSquare(white))
	}
	if board.getKingSquare(black) != none {
		println("Black king square = ", board.getKingSquare(black))
		t.Fatalf("black king square = %v, want none", board.getKingSquare(black))
	}
}

func TestGetKingSquareReturnsFirstKingIfMultipleKings(t *testing.T) {
	board := NewBoard()
	board.FenSetup(FenStartPosition)
	board.addPiece(K, e1)
	board.addPiece(K, d1)
	board.addPiece(k, e8)
	board.addPiece(k, d8)

	if board.getKingSquare(white) != d1 {
		println("White king square = ", board.getKingSquare(white), d1, e1, d8, e8)
		t.Fatalf("white king square = %v, want e1", BitboardSquares[board.getKingSquare(white)])
	}
	if board.getKingSquare(black) != d8 {
		println("Black king square = ", board.getKingSquare(black))
		t.Fatalf("black king square = %v, want e8", board.getKingSquare(black))
	}
	board.removePiece(K, d1)
	board.removePiece(k, d8)

	if board.getKingSquare(white) != e1 {
		println("White king square = ", board.getKingSquare(white))
		t.Fatalf("white king square = %v, want d1", BitboardSquares[board.getKingSquare(white)])
	}
	if board.getKingSquare(black) != e8 {
		println("Black king square = ", board.getKingSquare(black))
		t.Fatalf("black king square = %v, want d8", board.getKingSquare(black))
	}
	board.removePiece(K, e1)
	board.removePiece(k, e8)

	if board.getKingSquare(white) != none {
		println("White king square = ", board.getKingSquare(white))
		t.Fatalf("white king square = %v, want none", board.getKingSquare(white))
	}
	if board.getKingSquare(black) != none {
		println("Black king square = ", board.getKingSquare(black))
		t.Fatalf("black king square = %v, want none", board.getKingSquare(black))
	}
}

func TestMakeMovePromotionWithCapture(t *testing.T) {
	board := NewBoard()
	board.State.SideToMove = white
	board.addPiece(P, a7)
	board.addPiece(r, a8)

	before := snapshotBoard(board)
	move := NewMove(a7, a8, PromotionR)

	board.MakeMove(move)
	board.UndoMove(move)

	assertBoardSnapshotEqual(t, board, before)
	if board.History.size() != 0 {
		t.Fatalf("history size = %d, want 0", board.History.size())
	}
}
