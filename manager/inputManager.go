package manager

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/png"
)

type InputManager struct {
}

func NewInputManager() *InputManager{
	return &InputManager{}
}

func (in *InputManager) Keyboard() {

}

func (in *InputManager) Mouse() {

}

func (in *InputManager) IsMouseClick() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return false
}

func (in *InputManager) CheckMouseCursorOnBtn(r image.Rectangle) bool {
	x, y := ebiten.CursorPosition()
	point := image.Pt(x, y)
	if point.In(r) {
		return true
	}
	return false
}
