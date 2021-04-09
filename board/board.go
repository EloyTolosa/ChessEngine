package board

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"ChessEngine/globals"
	"ChessEngine/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	// Table dimensions (squared)
	tDimensions = 8
	// Table inital value
	tInitValue table = 18446462598732906495
	// NillValue is all 1's, so its the max value for an uint64
	tNilValue table = 18446744073709551615
)

var (
	images map[PieceType]*ebiten.Image
)

// A table is a number which represents the cells that have Pieces in it. For
// instance, the initial table would be represented by the number
// 18446462598732906495, which is the decimal representation for
//
// 1 1 1 1 1 1 1 1
// 1 1 1 1 1 1 1 1
// 0 0 0 0 0 0 0 0
// 0 0 0 0 0 0 0 0
// 0 0 0 0 0 0 0 0
// 0 0 0 0 0 0 0 0
// 1 1 1 1 1 1 1 1
// 1 1 1 1 1 1 1 1
//
// Starting from the upper left corner down until the lower right one
type table uint64

type Board struct {
	pieces        []Piece
	table         table
	lastTable     table
	movementTable table
	// TODO: move images to globals package
	changed bool
}

func (board *Board) UpdateTable() {
	// Update table code here ...
	board.table = board.lastTable
}

func (board *Board) SetPieceMovements(xpos, ypos int) {
	// Get piece in position xpos,ypos
	cWidth := int(globals.WindowWidth / tDimensions)
	cHeight := int(globals.WindowHeight / tDimensions)
	xLog := xpos / cWidth
	yLog := ypos / cHeight
	board.movementTable = table(board.pieces[yLog*tDimensions+xLog].GetAvailableMovements())
	log.Printf("Available movements for piece are: %d\n", board.movementTable)
}

func (board *Board) IsNilTable() bool {
	return board.lastTable == tNilValue
}

func (board *Board) HasChanged() bool {
	return board.changed
}

func (board *Board) SetChanged(ch bool) {
	board.changed = ch
}

func (board *Board) LoadImages() {

	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	images = make(map[PieceType]*ebiten.Image)

	// Append textures so we don't have to search them after this
	images[WhitePawn] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_plt60.png", currDir, "images"))
	images[BlackPawn] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_pdt60.png", currDir, "images"))
	images[WhiteBishop] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_blt60.png", currDir, "images"))
	images[BlackBishop] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_bdt60.png", currDir, "images"))
	images[WhiteKnight] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_nlt60.png", currDir, "images"))
	images[BlackKnight] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_ndt60.png", currDir, "images"))
	images[WhiteRook] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_rlt60.png", currDir, "images"))
	images[BlackRook] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_rdt60.png", currDir, "images"))
	images[WhiteKing] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_klt60.png", currDir, "images"))
	images[BlackKing] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_kdt60.png", currDir, "images"))
	images[WhiteQueen] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_qlt60.png", currDir, "images"))
	images[BlackQueen] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_qdt60.png", currDir, "images"))

}

func (board *Board) InitBoard() {

	// Initial table values
	board.table = tInitValue
	board.lastTable = tNilValue

	// Changed must be true in fisrt iteration
	board.changed = true

	// Set initial pieces values
	board.pieces = make([]Piece, 64)
	// White and black's pawns
	for i := 0; i < 8; i++ {
		board.pieces[i+48] = NewPiece(WhitePawn, PiecePosition(48+i))
		board.pieces[i+8] = NewPiece(BlackPawn, PiecePosition(8+i))
	}
	// White player's pieces
	board.pieces[58] = NewPiece(WhiteBishop, PiecePosition(58))
	board.pieces[61] = NewPiece(WhiteBishop, PiecePosition(61))
	board.pieces[57] = NewPiece(WhiteKnight, PiecePosition(57))
	board.pieces[62] = NewPiece(WhiteKnight, PiecePosition(62))
	board.pieces[56] = NewPiece(WhiteRook, PiecePosition(56))
	board.pieces[63] = NewPiece(WhiteRook, PiecePosition(63))
	board.pieces[60] = NewPiece(WhiteKing, PiecePosition(60))
	board.pieces[59] = NewPiece(WhiteQueen, PiecePosition(59))
	// Black player's pieces
	board.pieces[2] = NewPiece(BlackBishop, PiecePosition(2))
	board.pieces[5] = NewPiece(BlackBishop, PiecePosition(5))
	board.pieces[1] = NewPiece(BlackKnight, PiecePosition(1))
	board.pieces[6] = NewPiece(BlackKnight, PiecePosition(6))
	board.pieces[0] = NewPiece(BlackRook, PiecePosition(0))
	board.pieces[7] = NewPiece(BlackRook, PiecePosition(7))
	board.pieces[4] = NewPiece(BlackKing, PiecePosition(4))
	board.pieces[3] = NewPiece(BlackQueen, PiecePosition(3))

}

func (board *Board) Paint(screen *ebiten.Image) {
	// Table cell dimensions
	cWidth := int(globals.WindowWidth / tDimensions)
	cHeight := int(globals.WindowHeight / tDimensions)
	for i := 0; i < tDimensions; i++ {
		for j := 0; j < tDimensions; j++ {
			if (i+j)%2 == 0 {
				ebitenutil.DrawRect(screen, float64(i*cWidth), float64(j*cHeight), float64(cWidth), float64(cHeight), color.White)
			}
		}
	}
	// Paint textures (Pieces) on the board
	for _, p := range board.pieces {
		if p != Piece(0) {
			pPosition := p.getPosition()
			xLogic := int(pPosition % 8)
			yLogic := int(pPosition / 8)
			x := float64(xLogic * cWidth)
			y := float64((yLogic) * cHeight)
			// Center the images
			x += (float64(cWidth) - 60) / 2
			y += (float64(cHeight) - 60) / 2
			// Apply transformations
			geom := &ebiten.GeoM{}
			geom.Translate(x, y)
			// TODO read images in SVG format and scale them
			screen.DrawImage(
				images[p.getPieceType()],
				&ebiten.DrawImageOptions{
					GeoM: *geom,
				},
			)
		}

	}

}
