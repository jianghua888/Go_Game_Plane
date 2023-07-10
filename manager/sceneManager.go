package manager

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jianghua/game/utils"
	"image"
	"image/color"
	"math/rand"
	"strings"
	"sync"
)

const (
	GameStart = 1
	GameChooseSpirit = 2
	GameEnter = 3
)


var mutex sync.RWMutex

type SceneManager struct {
	Nodes        map[string]*ebiten.Image   //节点图片信息
	currNo       int8                       //当前场景
	NodePosition map[string]image.Rectangle //场景中所有节点的位置信息
	Enemy map[*ebiten.Image] image.Rectangle
	EnemyBullet map[*ebiten.Image] image.Rectangle
	EnemyHit map[*ebiten.Image] int

}

func NewSceneManager() *SceneManager {
	return &SceneManager {
		Nodes:        make(map[string]*ebiten.Image),
		currNo:       1,
		NodePosition: make(map[string]image.Rectangle),
		Enemy:make(map[*ebiten.Image]image.Rectangle),
		EnemyBullet:make(map[*ebiten.Image]image.Rectangle),
		EnemyHit:make(map[*ebiten.Image]int),
	}
}

func (scene *SceneManager) GetCurrNo() int8 {
	return scene.currNo
}

func (scene *SceneManager) SetCurrNo(currNo int8) {
	scene.currNo = currNo
}

func (scene *SceneManager) AddNodePosition(key string, r image.Rectangle) {
	scene.NodePosition[key] = r
}

func (scene *SceneManager) GetNodePostion(key string) image.Rectangle {
	return scene.NodePosition[key]
}

func (scene *SceneManager) InitStart(screen *ebiten.Image) {
	//加载开始背景
	startImg := utils.CreateImage("./images/start.png")
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(startImg, op)
	scene.Nodes["start"] = startImg

	//加载开始按钮
	startBtnImg := utils.CreateImage("./images/start_btn.png")
	btnOp := &ebiten.DrawImageOptions{}
	btnOp.GeoM.Translate(250, 450)
	screen.DrawImage(startBtnImg, btnOp)
	scene.Nodes["start_btn"] = startBtnImg

	imgW, imgH := startBtnImg.Size()
	rect := image.Rect(250, 450, 250+imgW, 450+imgH)
	scene.AddNodePosition("start_btn", rect)

}

func (scene *SceneManager) Enter() {
	for key, _ := range scene.Nodes {
		scene.Nodes[key].Clear()
	}

}

//选择我方战机
func (scene *SceneManager) EnterChooseSpirit(screen *ebiten.Image) {

	//加载选择战机图片
	img := utils.CreateImage("./images/select.png")
	selectOp := &ebiten.DrawImageOptions{}
	selectOp.GeoM.Translate(150, 0)
	screen.DrawImage(img, selectOp)
	scene.Nodes["select"] = img

	//确定按钮
	cfImg := utils.CreateImage("./images/cf.png")
	cfOp := &ebiten.DrawImageOptions{}
	cfOp.GeoM.Translate(200, 500)
	screen.DrawImage(cfImg, cfOp)
	scene.Nodes["confirm_btn"] = cfImg

	imgW, imgH := cfImg.Size()
	rect := image.Rect(200, 500, 200+imgW, 500+imgH)
	scene.AddNodePosition("confirm_btn", rect)


	spirit := [3]string{"s1.png", "s2.png", "s3.png"}

	defaultY := float64(150)

	//offsetx := 80
	offsetxMap := [3]int{100, 250, 400}
	for key, val := range spirit {

		nodeName := "spirit_" + string(key)
		//获取前面节点的位置，在基础上增加偏移

		offsetx := offsetxMap[key]
		img := utils.CreateImage("./images/" + val)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(offsetx), defaultY)
		screen.DrawImage(img, op)
		scene.Nodes[nodeName] = img

		//设置位置信息
		arrStr := strings.Split(val, ".")
		name := "spirit_" + arrStr[0]
		imgW, imgH := img.Size()
		rect := image.Rect(offsetx, int(defaultY), offsetx+imgW, int(defaultY)+imgH)
		scene.AddNodePosition(name, rect)
	}
}

func (scene *SceneManager) DrawChooseSpiritEffect(screen *ebiten.Image, target *ebiten.Image, nodeName string) {
	w, h := target.Size()
	effectImg := ebiten.NewImage(w, h)
	effectImg.Fill(color.RGBA{R: 200, G: 200, B: 200, A: 0})
	op := &ebiten.DrawImageOptions{}

	nodePos := scene.GetNodePostion(nodeName)
	op.GeoM.Translate(float64(nodePos.Min.X), float64(nodePos.Min.Y))
	screen.DrawImage(effectImg, op)
}

func (scene *SceneManager) DrawMySpirit(screen *ebiten.Image, spirit *MySpirit) {
	spiritName := strings.Split(spirit.Name, "_")[1]
	spiritImg := utils.CreateImage("./images/" + spiritName + ".png")
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(spirit.Re.Min.X), float64(spirit.Re.Min.Y))
	screen.DrawImage(spiritImg, op)

	w, h := spiritImg.Size()
	rect := image.Rect(spirit.Re.Min.X, spirit.Re.Min.Y, w+spirit.Re.Min.X, h+spirit.Re.Min.Y)
	spirit.Re = rect
}


func (scene *SceneManager) AddMySpiritBullet(spirit *MySpirit) {

	bulletImg := utils.CreateImage("./images/b1.png")
	w, h := bulletImg.Size()
	rect := image.Rect(spirit.Re.Min.X, spirit.Re.Min.Y, w+spirit.Re.Min.X, h+spirit.Re.Min.Y)
	spirit.Bullet[bulletImg] = rect
}


func (scene *SceneManager) DrawMySpiritBullet(screen *ebiten.Image, spirit *MySpirit) {

	if len(spirit.Bullet) > 0 {
		mutex.Lock()
		for img, val := range spirit.Bullet {
			bulletImg := utils.CreateImage("./images/b1.png")
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(val.Min.X), float64(val.Min.Y)-5)
			val.Min.Y = val.Min.Y - 5
			val.Max.Y = val.Max.Y - 5
			screen.DrawImage(bulletImg, op)
			spirit.Bullet[img] = val

			//超出距离，则移除
			if val.Min.Y <= 0 {
				delete(spirit.Bullet, img)
			}
		}
		mutex.Unlock()
	}
}




func (scene *SceneManager) AddEnemy() {

	img := utils.CreateImage("./images/e4.png")
	w, h := img.Size()
	xpos := rand.Intn(500)
	rect := image.Rect(xpos, 0, w+xpos, h)
	scene.Enemy[img] = rect
}


func (scene *SceneManager) DrawEnemy(screen *ebiten.Image) {
	if len(scene.Enemy) > 0 {
		mutex.Lock()
		for img, val := range scene.Enemy {

			enemyImg := utils.CreateImage("./images/e4.png")

			_,ok := scene.EnemyHit[img]
			if ok {
				enemyImg  = utils.CreateImage("./images/baozha.png")
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(val.Min.X), float64(val.Min.Y)+1)
			val.Min.Y = val.Min.Y + 1
			val.Max.Y = val.Max.Y + 1
			screen.DrawImage(enemyImg, op)

			if !ok {
				scene.Enemy[img] = val
			} else {
				delete(scene.Enemy, img)
			}

			//超出距离，则移除
			if val.Min.Y > 800 {
				delete(scene.Enemy, img)
			}
		}
		mutex.Unlock()
	}
}

func (scene *SceneManager) AddEnemyBullet() {

	if len(scene.Enemy) > 0 {
		for _, enemy := range scene.Enemy {
			img := utils.CreateImage("./images/e4.png")
			w,h := img.Size()
			rect := image.Rect(enemy.Min.X + w/2, enemy.Min.Y + h, w+enemy.Max.X, h+enemy.Max.Y)
			scene.EnemyBullet[img] = rect
		}
	}
}

func (scene *SceneManager) DrawEnemyBullet(screen *ebiten.Image) {
	if len(scene.EnemyBullet) > 0 {
		mutex.Lock()
		for img, val := range scene.EnemyBullet {
			enemyBulletImg := utils.CreateImage("./images/eb1.png")
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(val.Min.X), float64(val.Min.Y)+1)
			val.Min.Y = val.Min.Y + 2
			val.Max.Y = val.Max.Y + 2
			screen.DrawImage(enemyBulletImg, op)
			scene.EnemyBullet[img] = val

			//超出距离，则移除
			if val.Min.Y > 800 {
				delete(scene.EnemyBullet, img)
			}
		}
		mutex.Unlock()
	}
}








func (scene *SceneManager) Exit() {

}
