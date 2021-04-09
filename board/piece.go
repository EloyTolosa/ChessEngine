package board

import "log"

const (
	// Piece constants they're a Piece type so they don't get confused
	WhitePawn   PieceType = 1
	BlackPawn   PieceType = 2
	WhiteKnight PieceType = 3
	BlackKnight PieceType = 6
	WhiteBishop PieceType = 4
	BlackBishop PieceType = 8
	WhiteRook   PieceType = 5
	BlackRook   PieceType = 10
	WhiteKing   PieceType = 7
	BlackKing   PieceType = 14
	WhiteQueen  PieceType = 9
	BlackQueen  PieceType = 18
)

const (
	positionMask = 2016
	PieceMask    = 31
)

// A movement is a number representing the positions that the piece can take
// in a general environment. As the table, the number will be a 64-bit integer
//
// For example, the White pawn, from its starting position, can move up to 2016
// cells up. That is, with the binary representation of the table:
// ============================================================================
// 0 0 0 0 0 0 0 0
// 0 0 0 0 0 0 0 0
// 0 0 0 0 0 0 0 0
// 0 0 0 0 0 0 0 0
// 0 0 0 0 0 0 0 0
// 1 0 0 0 0 0 0 0
// 1 0 0 0 0 0 0 0
// 0 0 0 0 0 0 0 0
// Which, in decimal represents the number 8421376 (starting from top left to
// bottom right)
//
// Obviously, to this piece we have to add the position in which the piece its
// located at the moment. Because that will change the number representing the
// movement.
// In the previous case, the pawn was its initial state. In this case, we are
// going to say that the pawn has moved one cell up, and it's not the leftmost
// pawn, but the second pawn (the B2 cell pawn).
// ============================================================================
// Pawn moves 1 up => movement < 8 (left shift 8 bits)
// B2 pawn => movement > 1 (right shift 1 bit)
//
// ************************* MOVEMENT RULES ***********************************
// Pawns:
//     - Every cell up they move, shift 8 bits to the left
//     - Starting from cell A2, each pawn will shift 1 bit to the right, being
//       the A2 cell index 0, B2 index 1, etc:
//       [ A2, B2, C2, D2, E2, F2, G2, H2 ] pawn number
//       [ 0 , 1 , 2 , 3 , 4 , 5 , 6 , 7  ] bits to shift right
//
// Rook:
//     Rooks can move straight in every direction

var (
	Movements = map[PieceType]uint64{
		WhitePawn: 8421376, BlackPawn: 36169534507319296,
	}
)

// A Piece is represented with a 10 bit number. From this 10 bits, the 6 leftmost
// of them represent the position, and the 5 rightmost of them represent the Piece
// itself.
//
// For example, a White rook at E4 (which is the cell number 36, if we say that
// cell 0 is the column A8), would be represented like this
// ############################################################################
// 100100 => cell number 36, E4
// 1010   => White rook, represented by the number 5
// Piece  => (cellNumber << 5)|PieceNumber
// Piece  => (100100 << 5)|00101 => 10010000101 => 1157
//
// To get the Piece or the position, we just need to use a mask and perform the
// bitwise AND operation
//
// Piece mask    => 00000011111 => 31
// position mask => 11111100000 => 2016
type Piece int

type PieceType int

type PiecePosition int

func NewPiece(pt PieceType, pp PiecePosition) Piece {
	return Piece(int(pp<<5) | int(pt))
}

func (piece Piece) GetAvailableMovements() uint64 {
	switch piece.getPieceType() {
	case WhitePawn:
		ppos := piece.getPosition()
		bitesToShiftLeft := 6 - (ppos / 8)
		bitesToShiftRight := ppos - 48
		return (Movements[WhitePawn] >> bitesToShiftRight) << bitesToShiftLeft
	default:
		log.Printf("Not implemented yet")
		return 0
	}
}

func (Piece Piece) getPosition() PiecePosition {
	return PiecePosition((int(Piece) & positionMask) >> 5)
}

func (Piece Piece) getPieceType() PieceType {
	return PieceType(int(Piece) & PieceMask)
}
