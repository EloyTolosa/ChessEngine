package utils

import (
	"ChessEngine/globals"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Returns a ebiten.Image object from a file path.
// Fatals if error
func NewImage(path string) (img *ebiten.Image) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return img
}

// Returns the logical position (x,y) from a pixel coordinate
func GetLogicalPosition(absX, absY int) (logX, logY int) {
	cWidth := int(globals.WindowWidth / globals.TableDim)
	cHeight := int(globals.WindowHeight / globals.TableDim)
	return (absX / cWidth), (absY / cHeight)
}
