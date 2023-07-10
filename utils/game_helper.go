package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func CreateImage(imgfile string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(imgfile)
	if err != nil {
		return nil
	}
	return img
}
