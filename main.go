package main

import (
	"ChessEngine/board"
	"ChessEngine/globals"
	"ChessEngine/utils"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const ()

type App struct {
	Board *board.Board
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (app *App) Update() (err error) {
	// Check for mouse pressed events
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		var x, y int
		x, y = ebiten.CursorPosition()
		if err = app.Board.SetPieceMovements(x, y); err != nil {
			if err == board.ErrNoPieceAtPos {
				app.Board.SetClicked(false)
				app.Board.SetClickedAt(0, 0)
				app.Board.ResetMovements()
			} else {
				return err
			}
		} else {
			app.Board.SetClicked(true)
			app.Board.SetClickedAt(utils.GetLogicalPosition(float64(x), float64(y)))
		}
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (app *App) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	app.Board.Paint(screen)
	// update board with last frame value
	app.Board.UpdateState()
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (app *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return globals.WindowWidth, globals.WindowHeight
}

func (app *App) initApp() {

	// Window size and title
	ebiten.SetWindowSize(globals.WindowWidth, globals.WindowHeight)
	ebiten.SetWindowTitle("Golang Chess Engine (GOCHEN)")
	// If the board has not changed in between frames the window should not be cleared and
	// should stay the same
	ebiten.SetScreenClearedEveryFrame(false)

	// Initializes app struct and prepares everything just to be painted
	app.Board = &board.Board{}
	app.Board.InitBoard()
	app.Board.LoadImages()

}

func main() {

	app := &App{}
	app.initApp()
	if err := ebiten.RunGame(app); err != nil {
		log.Fatalln(err.Error())
	}
}
