package board

import (
    
    "ChessEngine/utils"

    "os"
    "log"
    "fmt"

    "github.com/hajimehoshi/ebiten/v2"
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
	positionMask = 2016
	PieceMask    = 31
)

var (
    Movements = map[PieceType][]int {
        WhitePawn: {1,2}, BlackPawn: {-1,-2},
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

func NewPiece(pt PieceType, pp PiecePosition) Piece {
	return Piece(int(pp<<5) | int(pt))
}

type PiecePosition int

func (Piece Piece) getPosition() PiecePosition {
	return PiecePosition((int(Piece) & positionMask) >> 5)
}

func (Piece Piece) getPieceType() PieceType {
	return PieceType(int(Piece) & PieceMask)
}

func (Piece Piece) getImage() (img *ebiten.Image) {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch Piece.getPieceType() {
	case WhitePawn:
		return utils.NewImage(fmt.Sprintf("%s/%s/Chess_plt60.png", currDir, "images"))
	case BlackPawn:
		return utils.NewImage(fmt.Sprintf("%s/%s/Chess_pdt60.png", currDir, "images"))
	case WhiteBishop:
		return utils.NewImage(fmt.Sprintf("%s/%s/Chess_blt60.png", currDir, "images"))
	case BlackBishop:
		return utils.NewImage(fmt.Sprintf("%s/%s/Chess_bdt60.png", currDir, "images"))
	case WhiteKnight:
		return utils.NewImage(fmt.Sprintf("%s/%s/Chess_klt60.png", currDir, "images"))
	case BlackKnight:
		return utils.NewImage(fmt.Sprintf("%s/%s/Chess_kdt60.png", currDir, "images"))
	case WhiteRook:
		return utils.NewImage(fmt.Sprintf("%s/%s/Chess_rlt60.png", currDir, "images"))
	case BlackRook:
		return utils.NewImage(fmt.Sprintf("%s/%s/Chess_rdt60.png", currDir, "images"))
	case WhiteKing:
		return utils.NewImage(fmt.Sprintf("%s/%s/Chess_klt60.png", currDir, "images"))
	case BlackKing:
		return utils.NewImage(fmt.Sprintf("%s/%s/Chess_kdt60.png", currDir, "images"))
	case WhiteQueen:
		return utils.NewImage(fmt.Sprintf("%s/%s/Chess_qlt60.png", currDir, "images"))
	case BlackQueen:
		return utils.NewImage(fmt.Sprintf("%s/%s/Chess_qdt60.png", currDir, "images"))
	default:
		return nil
	}
}


