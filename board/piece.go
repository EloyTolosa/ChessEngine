package board

import (
	"ChessEngine/globals"
	"log"
	"math"
)

const (
	// Piece constants they're a Piece type so they don't get confused
	WhitePawn PieceType = iota
	BlackPawn
	WhiteKnight
	BlackKnight
	WhiteBishop
	BlackBishop
	WhiteRook
	BlackRook
	WhiteKing
	BlackKing
	WhiteQueen
	BlackQueen
)

const (
	positionMask uint16 = 65280
	pieceMask    uint16 = 255
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
// 		- Pawns can move one or two cells up. Pawns can move two cells in case of
// 			the pawn being in the initial spot.
// 			Pawns can also move one cell in diagonal, in the case of capturing.
//			Also, pawns have a special move called "an-passant", meaning that you
// 			can capture an oponent pawn without the pawn being strictly in diagonal.
//     	- Every cell up they move, shift 8 bits to the left
//     	- Starting from cell A2, each pawn will shift 1 bit to the right, being
// 			the A2 cell index 0, B2 index 1, etc:
//      	[ A2, B2, C2, D2, E2, F2, G2, H2 ] pawn number
//      	[ 0 , 1 , 2 , 3 , 4 , 5 , 6 , 7  ] bits to shift right
//
// Rook:
// 		- Rooks move in straight lines, either up/down, or left/right, and can move
// 			the ammount of cells they want to.

// A Piece is represented with a 16 bit number. The 8 leftmost bits represent the
// position, and the 8 rightmost of them represent the type of the piece.
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
// Piece mask    => 0000000011111111 => (2^8)-1 => 255
// position mask => 1111111100000000 => 65280
type Piece uint16

type PieceType uint8

type PiecePosition uint8

func (pos PiecePosition) getX() int {
	return int(pos % globals.TableDim)
}

func (pos PiecePosition) getY() int {
	return int(pos / globals.TableDim)
}

func (pos PiecePosition) isOutOfBounds() bool {
	return pos > 63
}

func NewPiece(pt PieceType, pp int) *Piece {
	p := Piece(pp<<8 | int(pt))
	return &p
}

func (piece *Piece) MoveTo(to int) {
	piece.setPosition(PiecePosition(to))
}

func AreSameColor(p1, p2 *Piece) bool {
	// white pieces are even numbers, and black pieces are odd numbers
	// to check if they are the same color, we just have to check if both are
	// either even or odd
	return (p1.getPieceType()%2 == p2.getPieceType()%2)
}

func (p *Piece) isIllegalMove(board *Board, newpos int) bool {
	// current and new piece positions
	ppos := PiecePosition(p.getPosition())
	npos := PiecePosition(newpos)

	px, py := ppos.getX(), ppos.getY()
	nx, ny := npos.getX(), npos.getY()

	xdiff := math.Abs(float64(px - nx))
	ydiff := math.Abs(float64(py - ny))

	switch p.getPieceType() {
	case WhitePawn, BlackPawn:
		// new position is out of bounds
		if npos.isOutOfBounds() {
			return true
		}

		// pawn in left side moves diagonal and goes to the right side
		if ((px == 0) && (nx == 7)) ||
			// pawn in right sido moves diagonal and goes to the left side
			((px == 7) && (nx == 0)) {
			return true
		}

		// pawn can only move diagonal if and only if another piece from the oposite color is
		// in a diagonal and the new move is a giagonal one
		if xdiff == 1 && ydiff == 1 {
			if !board.isThereAPieceAt(newpos) {
				return true
			} else {
				// we already checked if there's a piece there, so we do not have to check the error here
				np, _ := board.GetPieceAt(nx, ny)
				return AreSameColor(p, np)
			}
		}
		// a pawn can only move two steps if its in the orignal state
		if ydiff == 2 {
			// get y axis of the piece (row)
			return (py < 6) && (py > 1)
		}
		// and obviously, a pawn cannot move forward if any piece is in front of him
		if (ydiff == 1 || ydiff == 2) && board.isThereAPieceAt(newpos) {
			return true
		}
		// otherwise, its a good move
		return false

	case WhiteKnight, BlackKnight:

		// check out of bounds
		if npos.isOutOfBounds() {
			return true
		}

		// check if knight overflows from the bottom, top, left or right
		if (px <= 1 && nx >= 6) || (px >= 6 && nx <= 1) ||
			(py <= 1 && ny >= 6) || (py >= 6 && ny <= 1) {
			return true
		}

		// check if there's a piece in the new position
		if board.isThereAPieceAt(newpos) {
			np, _ := board.GetPieceAt(nx, ny)
			if AreSameColor(p, np) {
				return true
			}
		}

		return false

	case WhiteBishop, BlackBishop:

		// check out of bounds
		if npos.isOutOfBounds() {
			return true
		}

		return false
	default:
		log.Println("Not implemented")
		return false
	}
}

func (p *Piece) GetAvailableMovements(board *Board) (newPositions []int) {
	ppos := int(p.getPosition())
	for _, m := range Movements[p.getPieceType()] {
		for r := 1; r <= m.Limit; r++ {
			npos := ppos + m.Move(r)
			// check position validity
			if !p.isIllegalMove(board, npos) {
				// if it does not overflow, we append it
				newPositions = append(newPositions, npos)
			}
		}
	}
	return
}

func (p *Piece) GetLogicalPosition() (logX int, logY int) {
	pp := p.getPosition()
	return int(pp % globals.TableDim), int(pp / globals.TableDim)
}

func (p *Piece) getPosition() int {
	return int(uint16(*p) & positionMask >> 8)
}

func (p *Piece) setPosition(position PiecePosition) {
	newpos := (uint16(position) << 8) // xxxxxxxx11111111
	*p = Piece(newpos | uint16(p.getPieceType()))
}

func (p *Piece) getPieceType() PieceType {
	return PieceType(uint16(*p) & pieceMask)
}
