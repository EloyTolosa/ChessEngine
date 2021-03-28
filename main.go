package main

import (
	"fmt"
	raylib "github.com/gen2brain/raylib-go/raylib"
	"log"
	"os"
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

	// Window constants
	wWidth  = 800
	wHeight = 800

	// App constants
	maxFPS = 60
)

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

func (piece piece) getImage() (img *raylib.Image) {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch piece.getPieceType() {
	case WhitePawn:
		return raylib.LoadImage(fmt.Sprintf("%s/%s/Chess_plt60.png", currDir, "images"))
	case BlackPawn:
		return raylib.LoadImage(fmt.Sprintf("%s/%s/Chess_pdt60.png", currDir, "images"))
	case WhiteBishop:
		return raylib.LoadImage(fmt.Sprintf("%s/%s/Chess_blt60.png", currDir, "images"))
	case BlackBishop:
		return raylib.LoadImage(fmt.Sprintf("%s/%s/Chess_bdt60.png", currDir, "images"))
	case WhiteKnight:
		return raylib.LoadImage(fmt.Sprintf("%s/%s/Chess_klt60.png", currDir, "images"))
	case BlackKnight:
		return raylib.LoadImage(fmt.Sprintf("%s/%s/Chess_kdt60.png", currDir, "images"))
	case WhiteRook:
		return raylib.LoadImage(fmt.Sprintf("%s/%s/Chess_rlt60.png", currDir, "images"))
	case BlackRook:
		return raylib.LoadImage(fmt.Sprintf("%s/%s/Chess_rdt60.png", currDir, "images"))
	case WhiteKing:
		return raylib.LoadImage(fmt.Sprintf("%s/%s/Chess_klt60.png", currDir, "images"))
	case BlackKing:
		return raylib.LoadImage(fmt.Sprintf("%s/%s/Chess_kdt60.png", currDir, "images"))
	case WhiteQueen:
		return raylib.LoadImage(fmt.Sprintf("%s/%s/Chess_qlt60.png", currDir, "images"))
	case BlackQueen:
		return raylib.LoadImage(fmt.Sprintf("%s/%s/Chess_qdt60.png", currDir, "images"))
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
		player.pieces[12] = NewPiece(WhiteRook, piecePosition(57))
		player.pieces[13] = NewPiece(WhiteRook, piecePosition(63))
		player.pieces[14] = NewPiece(WhiteKing, piecePosition(60))
		player.pieces[15] = NewPiece(WhiteQueen, piecePosition(59))
	} else if playerType == "black" {
		player.pieces[8] = NewPiece(BlackBishop, piecePosition(8))
		player.pieces[9] = NewPiece(BlackBishop, piecePosition(9))
		player.pieces[10] = NewPiece(BlackKnight, piecePosition(10))
		player.pieces[11] = NewPiece(BlackKnight, piecePosition(11))
		player.pieces[12] = NewPiece(BlackRook, piecePosition(12))
		player.pieces[13] = NewPiece(BlackRook, piecePosition(13))
		player.pieces[14] = NewPiece(BlackKing, piecePosition(14))
		player.pieces[15] = NewPiece(BlackQueen, piecePosition(15))
	}
}

type Board struct {
	table       table
	WhitePlayer *Player
	BlackPlayer *Player
}

func (board *Board) Paint() {
	for i := 0; i < tDimensions; i++ {
		for j := 0; j < tDimensions; j++ {
			if (i+j)%2 == 0 {
				raylib.DrawRectangle(int32(i*cWidth), int32(j*cHeight), cWidth, cHeight, raylib.Black)
			}
		}
	}
}

type App struct {
	Board *Board
}

func (app *App) initApp() {

	// Initializes app struct and prepares everything just to be painted
	app.Board = &Board{}
	app.Board.table = tInitValue
	app.Board.WhitePlayer = &Player{}
	app.Board.WhitePlayer.initPlayer("white")
	app.Board.BlackPlayer = &Player{}
	app.Board.BlackPlayer.initPlayer("black")

	// Test, where is the first white pawn in the table
	log.Printf("White's first white pawn is in position %v", app.Board.WhitePlayer.pieces[0].getPosition())

}

func (app *App) paintBoard() {
	app.Board.Paint()
}

func (app *App) start() {

	raylib.InitWindow(wWidth, wHeight, "GOLANG GAME ENGINE (GOGEN)")

	raylib.SetTargetFPS(maxFPS)

	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.RayWhite)

		app.paintBoard()

		raylib.EndDrawing()
	}

}

func main() {

	app := &App{}
	app.initApp()
	app.start()

	raylib.CloseWindow()

	log.Println("This is your first golang graphic project")

}
