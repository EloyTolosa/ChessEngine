package main

import (
	"fmt"
	"log"
	"os"

    "image/color"
    _ "image/png"

    "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	// Piece constants they're a piece type so they don't get confused
	WhitePawn   pieceType = 1
	BlackPawn   pieceType = 2
	WhiteKnight pieceType = 3
	BlackKnight pieceType = 6
	WhiteBishop pieceType = 4
	BlackBishop pieceType = 8
	WhiteRook   pieceType = 5
	BlackRook   pieceType = 10
	WhiteKing   pieceType = 7
	BlackKing   pieceType = 14
	WhiteQueen  pieceType = 9
	BlackQueen  pieceType = 18

	positionMask = 2016
	pieceMask    = 31

	// Table dimensions (squared)
	tDimensions = 8
	// Table cell dimensions
	cWidth  = wWidth / tDimensions
	cHeight = wHeight / tDimensions
	// Table inital value
	tInitValue table = 18446462598732906495
    // NillValue is all 1's, so its the max value for an uint64
    tNilValue table = 18446744073709551615 

	// Window constants
	wWidth  = 800
	wHeight = 800

	// App constants
	maxFPS = 60
)

// #################################
// UTILS (to move to /lib/utils.go)
// #################################

// THis call fatals when error
func NewImage(path string) (img *ebiten.Image) {

    img, _, err := ebitenutil.NewImageFromFile(path)
    if err != nil {
        log.Fatalln(err.Error())
    }
    // log.Printf("Loading image %s\n", path)
    return img

}

// A table is a number which represents the cells that have pieces in it. For
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

// A piece is represented with a 10 bit number. From this 10 bits, the 6 leftmost
// of them represent the position, and the 5 rightmost of them represent the piece
// itself.
//
// For example, a White rook at E4 (which is the cell number 36, if we say that
// cell 0 is the column A8), would be represented like this
// ############################################################################
// 100100 => cell number 36, E4
// 1010   => White rook, represented by the number 5
// piece  => (cellNumber << 5)|pieceNumber
// piece  => (100100 << 5)|00101 => 10010000101 => 1157
//
// To get the piece or the position, we just need to use a mask and perform the
// bitwise AND operation
//
// piece mask    => 00000011111 => 31
// position mask => 11111100000 => 2016
type piece int

type pieceType int

func NewPiece(pt pieceType, pp piecePosition) piece {
	return piece(int(pp<<5) | int(pt))
}

type piecePosition int

func (piece piece) getPosition() piecePosition {
	return piecePosition((int(piece) & positionMask) >> 5)
}

func (piece piece) getPieceType() pieceType {
	return pieceType(int(piece) & pieceMask)
}

func (piece piece) getImage() (img *ebiten.Image) {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch piece.getPieceType() {
	case WhitePawn:
		return NewImage(fmt.Sprintf("%s/%s/Chess_plt60.png", currDir, "images"))
	case BlackPawn:
		return NewImage(fmt.Sprintf("%s/%s/Chess_pdt60.png", currDir, "images"))
	case WhiteBishop:
		return NewImage(fmt.Sprintf("%s/%s/Chess_blt60.png", currDir, "images"))
	case BlackBishop:
		return NewImage(fmt.Sprintf("%s/%s/Chess_bdt60.png", currDir, "images"))
	case WhiteKnight:
		return NewImage(fmt.Sprintf("%s/%s/Chess_klt60.png", currDir, "images"))
	case BlackKnight:
		return NewImage(fmt.Sprintf("%s/%s/Chess_kdt60.png", currDir, "images"))
	case WhiteRook:
		return NewImage(fmt.Sprintf("%s/%s/Chess_rlt60.png", currDir, "images"))
	case BlackRook:
		return NewImage(fmt.Sprintf("%s/%s/Chess_rdt60.png", currDir, "images"))
	case WhiteKing:
		return NewImage(fmt.Sprintf("%s/%s/Chess_klt60.png", currDir, "images"))
	case BlackKing:
		return NewImage(fmt.Sprintf("%s/%s/Chess_kdt60.png", currDir, "images"))
	case WhiteQueen:
		return NewImage(fmt.Sprintf("%s/%s/Chess_qlt60.png", currDir, "images"))
	case BlackQueen:
		return NewImage(fmt.Sprintf("%s/%s/Chess_qdt60.png", currDir, "images"))
	default:
		return nil
	}
}

type Player struct {
	pieces [16]piece
}

func (player *Player) initPlayer(playerType string) {
	// Pawns
	for i := 0; i < 8; i++ {
		if playerType == "white" {
			player.pieces[i] = NewPiece(WhitePawn, piecePosition(48+i))
		} else if playerType == "black" {
			player.pieces[i] = NewPiece(BlackPawn, piecePosition(8+i))
		}
	}
	// Rest of the pieces
	if playerType == "white" {
		player.pieces[8] = NewPiece(WhiteBishop, piecePosition(58))
		player.pieces[9] = NewPiece(WhiteBishop, piecePosition(61))
		player.pieces[10] = NewPiece(WhiteKnight, piecePosition(57))
		player.pieces[11] = NewPiece(WhiteKnight, piecePosition(62))
		player.pieces[12] = NewPiece(WhiteRook, piecePosition(56))
		player.pieces[13] = NewPiece(WhiteRook, piecePosition(63))
		player.pieces[14] = NewPiece(WhiteKing, piecePosition(60))
		player.pieces[15] = NewPiece(WhiteQueen, piecePosition(59))
	} else if playerType == "black" {
		player.pieces[8] = NewPiece(BlackBishop, piecePosition(2))
		player.pieces[9] = NewPiece(BlackBishop, piecePosition(5))
		player.pieces[10] = NewPiece(BlackKnight, piecePosition(1))
		player.pieces[11] = NewPiece(BlackKnight, piecePosition(6))
		player.pieces[12] = NewPiece(BlackRook, piecePosition(0))
		player.pieces[13] = NewPiece(BlackRook, piecePosition(7))
		player.pieces[14] = NewPiece(BlackKing, piecePosition(4))
		player.pieces[15] = NewPiece(BlackQueen, piecePosition(3))
	}
}

type Board struct {
	WhitePlayer *Player
	BlackPlayer *Player
	table       table
    lastTable   table
    images      map[pieceType]*ebiten.Image
    changed     bool
}

func (board *Board) hasChanged() bool {
   return board.table != board.lastTable 
}

func (board *Board) loadImages() {

	// Initialize textures first
	board.images = make(map[pieceType]*ebiten.Image)

	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Append textures so we don't have to search them after this
	board.images[WhitePawn] = NewImage(fmt.Sprintf("%s/%s/Chess_plt60.png", currDir, "images"))
	board.images[BlackPawn] = NewImage(fmt.Sprintf("%s/%s/Chess_pdt60.png", currDir, "images"))
	board.images[WhiteBishop] = NewImage(fmt.Sprintf("%s/%s/Chess_blt60.png", currDir, "images"))
	board.images[BlackBishop] = NewImage(fmt.Sprintf("%s/%s/Chess_bdt60.png", currDir, "images"))
	board.images[WhiteKnight] = NewImage(fmt.Sprintf("%s/%s/Chess_nlt60.png", currDir, "images"))
	board.images[BlackKnight] = NewImage(fmt.Sprintf("%s/%s/Chess_ndt60.png", currDir, "images"))
	board.images[WhiteRook] = NewImage(fmt.Sprintf("%s/%s/Chess_rlt60.png", currDir, "images"))
	board.images[BlackRook] = NewImage(fmt.Sprintf("%s/%s/Chess_rdt60.png", currDir, "images"))
	board.images[WhiteKing] = NewImage(fmt.Sprintf("%s/%s/Chess_klt60.png", currDir, "images"))
	board.images[BlackKing] = NewImage(fmt.Sprintf("%s/%s/Chess_kdt60.png", currDir, "images"))
	board.images[WhiteQueen] = NewImage(fmt.Sprintf("%s/%s/Chess_qlt60.png", currDir, "images"))
	board.images[BlackQueen] = NewImage(fmt.Sprintf("%s/%s/Chess_qdt60.png", currDir, "images"))

}

func (board *Board) initBoard() {
	board.table = tInitValue
    board.lastTable = tNilValue
    // True forced initial value
    board.changed = true
	board.WhitePlayer = &Player{}
	board.WhitePlayer.initPlayer("white")
	board.BlackPlayer = &Player{}
	board.BlackPlayer.initPlayer("black")
}

func (board *Board) Paint(screen *ebiten.Image) {
	// Paint black and white board
	for i := 0; i < tDimensions; i++ {
		for j := 0; j < tDimensions; j++ {
			if (i+j)%2 == 0 {
                ebitenutil.DrawRect(screen, float64(i*cWidth), float64(j*cHeight), float64(cWidth), float64(cHeight), color.White)
			}
		}
	}
	// Paint textures (pieces) on the board
	for _, p := range board.WhitePlayer.pieces {
		pPosition := p.getPosition()
		xLogic := pPosition % 8
		yLogic := pPosition / 8
		x := float64(xLogic * cWidth)
		y := float64((yLogic) * cHeight)
        geom := &ebiten.GeoM{}
        geom.Translate(x,y)
        log.Printf("Drawing %d in position (%f,%f)",p.getPieceType(),x,y)
        screen.DrawImage(
            board.images[p.getPieceType()], 
            &ebiten.DrawImageOptions{
                GeoM: *geom,
            },
        )
	}
	for _, p := range board.BlackPlayer.pieces {
		pPosition := p.getPosition()
		xLogic := pPosition % 8
		yLogic := pPosition / 8
    	x := float64(xLogic * cWidth)
		y := float64((yLogic) * cHeight)
        geom := &ebiten.GeoM{}
        geom.Translate(x,y)
        log.Printf("Drawing %d in position (%f,%f)",p.getPieceType(),x,y)
        screen.DrawImage(
            board.images[p.getPieceType()], 
            &ebiten.DrawImageOptions{
                GeoM: *geom,   
            },   
        )
	}
}

type App struct {
	Board *Board
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (app *App) Update() error {
    // Write your game's logical update.
    if app.Board.lastTable == tNilValue {
       // do nothing 
    } else {
        app.Board.changed = app.Board.hasChanged()
    }
    app.Board.lastTable = app.Board.table

    // Update actual table here...
    return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (app *App) Draw(screen *ebiten.Image) {
    // Write your game's rendering.
    if app.Board.changed {
        app.Board.Paint(screen)
    }
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (app *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return wWidth, wHeight
}

func (app *App) initApp() {

    // Window size and title
    ebiten.SetWindowSize(wWidth, wHeight)
    ebiten.SetWindowTitle("Golang Chess Engine (GOCHEN)")

	// Initializes app struct and prepares everything just to be painted
	app.Board = &Board{}
	app.Board.initBoard()
	app.Board.loadImages()

}

func main() {

	app := &App{}
	app.initApp()
    if err := ebiten.RunGame(app); err != nil {
        log.Fatalln(err.Error())
    }


}
