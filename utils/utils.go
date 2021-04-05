package utils

import (
    "log"
    
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2"
)

// THis call fatals when error
func NewImage(path string) (img *ebiten.Image) {

    img, _, err := ebitenutil.NewImageFromFile(path)
    if err != nil {
        log.Fatalln(err.Error())
    }
    // log.Printf("Loading image %s\n", path)
    return img

}


