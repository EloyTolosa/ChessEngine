package board

import (
	"ChessEngine/globals"
	"log"
	"math"
)

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

var (
	Movements = map[PieceType][]Movement{
		WhitePawn: {
			Movement{
				UP, 2,
			},
			Movement{
				func(i int) int {
					return UP(i) + RIGHT(i)
				}, 1,
			},
			Movement{
				func(i int) int {
					return UP(i) + LEFT(i)
				}, 1,
			},
		}, BlackPawn: {
			Movement{
				DOWN, 2,
			},
			Movement{
				func(i int) int {
					return DOWN(i) + RIGHT(i)
				}, 1,
			},
			Movement{
				func(i int) int {
					return DOWN(i) + LEFT(i)
				}, 1,
			},
		}, WhiteKnight: {
			Movement{
				func(i int) int {
					return UP(2*i) + RIGHT(i)
				}, 1,
			},
			Movement{
				func(i int) int {
					return UP(2*i) + LEFT(i)
				}, 1,
			},
			Movement{
				func(i int) int {
					return DOWN(2*i) + RIGHT(i)
				}, 1,
			},
			Movement{
				func(i int) int {
					return DOWN(2*i) + LEFT(i)
				}, 1,
			},
			Movement{
				func(i int) int {
					return UP(i) + RIGHT(i*2)
				}, 1,
			},
			Movement{
				func(i int) int {
					return UP(i) + LEFT(i*2)
				}, 1,
			},
			Movement{
				func(i int) int {
					return DOWN(i) + RIGHT(i*2)
				}, 1,
			},
			Movement{
				func(i int) int {
					return DOWN(i) + LEFT(i*2)
				}, 1,
			},
		},
	}
)

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

func NewPiece(pt PieceType, pp int) *Piece {
	p := Piece(pp<<8 | int(pt))
	return &p
}

func (piece *Piece) MoveTo(to int) {
	piece.setPosition(PiecePosition(to))
}

// TODO: REFACTOR
//
// we can get ppos from p.getPosition, there's no need to pass it as a parameter to the function
func (p *Piece) isIllegalMove(board *Board, ppos, npos int) bool {
	switch p.getPieceType() {
	case WhitePawn, BlackPawn:
		// new position is out of bounds
		if npos > 63 || npos < 0 {
			return true
		}
		// pawn in left side moves diagonal and goes to the right side

		// TODO: ADD FUNCTIONS
		//
		// - piece.getX()
		// - piece.getY()
		if ((ppos%globals.TableDim == 0) && (npos%globals.TableDim == 7)) ||
			// pawn in right sido moves diagonal and goes to the left side
			((ppos%globals.TableDim == 7) && (npos%globals.TableDim == 0)) {
			return true
		}
		// pawn can only move diagonal if and only if another piece from the oposite color is
		// in a diagonal and the new move is a giagonal one
		if (math.Abs(float64(ppos-npos)) == 7.0) || (math.Abs(float64(ppos-npos)) == 9.0) ||
			(math.Abs(float64(ppos+npos)) == 7.0) || (math.Abs(float64(ppos+npos)) == 9.0) {
			if !board.isThereAPieceAt(npos) {
				return true
			} else {
				xLog, yLog := int(npos%globals.TableDim), int(npos/globals.TableDim)
				// we already checked if there's a piece there, so we do not have to check the error here
				np, _ := board.GetPieceAt(xLog, yLog)
				// white pieces are the same as black pieces, but divided by 2, so to know if a piece is from
				// the opposite color, we have to perform the division and, in case they are the from different color,
				// the divission has to be either 2 or 1/2

				// TODO: ADD FUNCTION
				//
				// AreSameColor(piece1, piece2)
				d := float64(np.getPieceType() / p.getPieceType())
				return !(d == 2.0 || d == 1/2)
			}
		}
		// a pawn can only move two steps if its in the orignal state
		if math.Abs(float64(ppos-npos)) == 16.0 {
			// get y axis of the piece (row)
			y := int(ppos / globals.TableDim)
			return (y < 6) && (y > 1)
		}
		// and obviously, a pawn cannot move forward if any piece is in front of him
		if ((math.Abs(float64(ppos-npos)) == 8.0) || (math.Abs(float64(ppos-npos)) == 16.0)) &&
			board.isThereAPieceAt(npos) {
			return true
		}
		// otherwise, its a good move
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
			if !p.isIllegalMove(board, ppos, npos) {
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
