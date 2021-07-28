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
func GetLogicalPosition(absX, absY float64) (logX, logY int) {
	return (int(absX) / globals.CWidth), (int(absY) / globals.CHeight)
}

// Returns the absolute position from a logical position inside the board table
func GetAbsolutePosition(logX, logY int) (absX, absY float64) {
	return float64(logX * globals.CWidth), float64(logY * globals.CHeight)
}
