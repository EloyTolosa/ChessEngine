package board

import (
	"errors"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
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

	ErrNoPieceAtPos = errors.New("no piece at the given position")
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

type coordinate struct {
	x, y uint
}

type Board struct {
	// This variable saves the state of every individual piece. Its position, and its
	// type
	pieces []Piece
	// this saves the state of the pieces. If in between frames, the value of
	// the table is different, this means the board has changed
	tableCurrentFrame  table
	tablePreviousFrame table
	// this saves the state of the movements table. The movements table is a table
	// in which we store all possible movements from a piece.
	// This table is saved whenever the user clicks on a piece
	availableMovements []int
	// this saves a state of the clicked events
	// if in between frames, this two variables are equal, this means
	// the state of the board has not changed
	clickedCurrentFrame  bool
	clickedPreviousFrame bool
	// This saves the state of the location of the click in {x,y} coordinates
	// Being the coordinate {0,0} the top left corner of the table
	clickedAtCurrentFrame  coordinate
	clickedAtPreviousFrame coordinate
}

// Function that returns true if there is a piece at position p
func (board *Board) isThereAPieceAt(pos int) bool {
	// copy the table value
	b := board.tableCurrentFrame
	twoToThe63 := uint64(math.Pow(2, 63))
	// return true if there is a one at position 'pos' of the table
	// (bitmap implementation)
	return ((uint64(b) << pos) & twoToThe63) == twoToThe63
}

func (board *Board) UpdateState() {
	// Update table code here ...
	board.tablePreviousFrame = board.tableCurrentFrame
	board.clickedPreviousFrame = board.clickedCurrentFrame
	board.clickedAtPreviousFrame = board.clickedAtCurrentFrame
}

func (board *Board) SetPieceMovements(xpos, ypos int) (err error) {
	// Get piece in position xpos,ypos
	xLog, yLog := utils.GetLogicalPosition(float64(xpos), float64(ypos))
	p := board.pieces[yLog*tDimensions+xLog]
	// no piece at the given position
	if p == 0 {
		return ErrNoPieceAtPos
	}
	board.availableMovements = p.GetAvailableMovements(board)
	return
}

func (board *Board) IsNilTable() bool {
	return board.tableCurrentFrame == tNilValue
}

func (board *Board) HasChanged() bool {
	return ((board.clickedCurrentFrame != board.clickedPreviousFrame) ||
		(board.tableCurrentFrame != board.tablePreviousFrame) ||
		(board.clickedAtCurrentFrame != board.clickedAtPreviousFrame))
}

func (board *Board) IsClicked() bool {
	return board.clickedCurrentFrame
}

func (board *Board) SetClicked(cl bool) {
	board.clickedCurrentFrame = cl
}

func (board *Board) SetClickedAt(x, y int) {
	board.clickedAtCurrentFrame.x = uint(x)
	board.clickedAtCurrentFrame.y = uint(y)
}

func (board *Board) ResetMovements() {
	board.availableMovements = make([]int, 0)
}

func (board *Board) LoadImages() {

	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	images = make(map[PieceType]*ebiten.Image)

	// Append textures so we don't have to search them after this
	images[WhitePawn] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_plt60.png", currDir, "assets/images"))
	images[BlackPawn] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_pdt60.png", currDir, "assets/images"))
	images[WhiteBishop] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_blt60.png", currDir, "assets/images"))
	images[BlackBishop] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_bdt60.png", currDir, "assets/images"))
	images[WhiteKnight] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_nlt60.png", currDir, "assets/images"))
	images[BlackKnight] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_ndt60.png", currDir, "assets/images"))
	images[WhiteRook] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_rlt60.png", currDir, "assets/images"))
	images[BlackRook] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_rdt60.png", currDir, "assets/images"))
	images[WhiteKing] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_klt60.png", currDir, "assets/images"))
	images[BlackKing] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_kdt60.png", currDir, "assets/images"))
	images[WhiteQueen] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_qlt60.png", currDir, "assets/images"))
	images[BlackQueen] = utils.NewImage(fmt.Sprintf("%s/%s/Chess_qdt60.png", currDir, "assets/images"))

}

func (board *Board) InitBoard() {

	// Initial table values
	board.tableCurrentFrame = tInitValue
	board.tablePreviousFrame = tNilValue

	// initial clicked at values
	board.clickedAtCurrentFrame = coordinate{0, 0}
	board.clickedAtPreviousFrame = coordinate{0, 0}

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
	// log.Printf("painting screen...")
	board.paintCells(screen)
	board.paintPieces(screen)
	board.paintAvailableMovements(screen)
}

func (board *Board) paintCells(screen *ebiten.Image) {
	cWidth := int(globals.WindowWidth / tDimensions)
	cHeight := int(globals.WindowHeight / tDimensions)
	for i := 0; i < tDimensions; i++ {
		for j := 0; j < tDimensions; j++ {
			if (i+j)%2 == 0 {
				ebitenutil.DrawRect(screen, float64(i*cWidth), float64(j*cHeight), float64(cWidth), float64(cHeight), color.White)
			} else {
				ebitenutil.DrawRect(screen, float64(i*cWidth), float64(j*cHeight), float64(cWidth), float64(cHeight), color.Black)
			}
		}
	}
}

func (board *Board) paintPieces(screen *ebiten.Image) {
	// TODO: refactor
	for _, p := range board.pieces {
		if p != Piece(0) {
			// get x and y coordinates
			// pPosition := p.getPosition()
			logX, logY := p.GetLogicalPosition()
			x, y := utils.GetAbsolutePosition(logX, logY)
			// Center the images
			x += (float64(globals.CWidth) - 60) / 2
			y += (float64(globals.CHeight) - 60) / 2
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

func (board *Board) paintAvailableMovements(screen *ebiten.Image) {
	if !board.IsClicked() {
		return
	}
	// get image
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	imgPath := fmt.Sprintf("%s/%s/%s", currDir, "assets/images", "movement.png")
	movementImage := utils.NewImage(imgPath)
	// draw red circles at the given positions
	for _, p := range board.availableMovements {
		// get x and y coordinates
		xLogic := int(p % 8)
		yLogic := int(p / 8)
		// center the dots
		x := float64(xLogic * globals.CWidth)
		y := float64(yLogic * globals.CHeight)
		// declare geom struct
		geom := &ebiten.GeoM{}
		// scale (first)
		xf := float64(globals.CWidth) / float64(movementImage.Bounds().Dx())
		yf := float64(globals.CHeight) / float64(movementImage.Bounds().Dy())
		geom.Scale(xf, yf)
		// translate (then translate)
		geom.Translate(x, y)
		// draw image
		screen.DrawImage(movementImage, &ebiten.DrawImageOptions{
			GeoM: *geom,
		})
	}
}
