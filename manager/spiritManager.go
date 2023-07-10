package manager

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type MySpirit struct {
	Name string
	Img *ebiten.Image
	Re image.Rectangle
	Speed int
	Bullet map[*ebiten.Image] image.Rectangle //子弹
	BulletPos map[*ebiten.Image] *ebiten.DrawImageOptions
	IsShoot chan bool
}

func NewMySpirit() *MySpirit {
	return &MySpirit{
		Speed:8,
		Bullet: make(map[*ebiten.Image]image.Rectangle),
		BulletPos : make(map[*ebiten.Image] *ebiten.DrawImageOptions),
		IsShoot : make(chan bool),
	}
}

func (sp *MySpirit) SetXY(x,y int) {
	sp.Re.Min.X = x
	sp.Re.Min.Y = y
}

func (sp *MySpirit) GetX() int {
	return sp.Re.Min.X
}

func (sp *MySpirit) GetY() int {
	return sp.Re.Min.Y
}

func (sp *MySpirit) MoveLeft() {
	if sp.CheckBoundary(1) {
		return
	}
	sp.Re.Min.X -= sp.Speed
}

func (sp *MySpirit) MoveRight() {
	if sp.CheckBoundary(2) {
		return
	}
	sp.Re.Min.X += sp.Speed
}

func (sp *MySpirit) MoveUp() {
	if sp.CheckBoundary(3) {
		return
	}
	sp.Re.Min.Y -= sp.Speed
}

func (sp *MySpirit) MoveDown() {
	if sp.CheckBoundary(4) {
		return
	}
	sp.Re.Min.Y += sp.Speed
}


func (sp *MySpirit) Shoot()  {
	fmt.Println("shoot one")

	sp.IsShoot <- true

}


func (sp *MySpirit) CheckBoundary(dir int) bool {

	if dir == 1 && sp.GetX() - sp.Speed < 0 {
		return true
	}

	if dir == 2 && sp.Re.Max.X + sp.Speed  >= 600 {
		return true
	}

	if dir == 3 && sp.GetY() - sp.Speed < 0 {
		return true
	}

	if dir == 4 && sp.Re.Max.Y + sp.Speed > 800 {
		return true
	}

	return false
}




