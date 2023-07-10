package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jianghua/game/manager"
	"time"
	"log"
)

const (
	windowWidth  = 600
	windowHeight = 800
)


type Game struct {
	SceneM *manager.SceneManager
	StateM *manager.StateManager
	InputM *manager.InputManager
	MySpirit *manager.MySpirit
}

func (g *Game) Update() error {

	if g.SceneM.GetCurrNo() == manager.GameStart {
		if g.SceneM.Nodes["start_btn"] != nil {
			startBtnPos := g.SceneM.GetNodePostion("start_btn")
			if g.InputM.CheckMouseCursorOnBtn(startBtnPos) {
				//设置鼠标样式
				ebiten.SetCursorShape(ebiten.CursorShapePointer)
				if g.InputM.IsMouseClick() {
					g.StateM.Start()
					g.SceneM.SetCurrNo(manager.GameChooseSpirit)
				}
			} else {
				ebiten.SetCursorShape(ebiten.CursorShapeDefault)
			}
		}
	}

	//选择战机
	if g.SceneM.GetCurrNo() == manager.GameChooseSpirit {
		for k, v:= range [3]string{"s1", "s2", "s3"} {
			spiritKey := "spirit_" + v
			fmt.Println(spiritKey)
			spiritNode := g.SceneM.GetNodePostion(spiritKey)
			if g.InputM.CheckMouseCursorOnBtn(spiritNode) {
				ebiten.SetCursorShape(ebiten.CursorShapePointer)
				if g.InputM.IsMouseClick() {
					node := g.SceneM.Nodes["spirit_"+string(k)]
					g.MySpirit.Img = node
					g.MySpirit.Name = spiritKey
				}
				break
			} else {
				ebiten.SetCursorShape(ebiten.CursorShapeDefault)
			}
		}

		//确认按钮点击
		confirmBtnPos := g.SceneM.GetNodePostion("confirm_btn")
		if g.InputM.CheckMouseCursorOnBtn(confirmBtnPos) {
			ebiten.SetCursorShape(ebiten.CursorShapePointer)
			if g.InputM.IsMouseClick() {
				g.SceneM.SetCurrNo(manager.GameEnter)
			}
		} else {
			ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		}

	}

	//进入游戏
	if g.SceneM.GetCurrNo() == manager.GameEnter {

		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			g.MySpirit.MoveLeft()
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			g.MySpirit.MoveRight()
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			g.MySpirit.MoveUp()
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			g.MySpirit.MoveDown()
		}

		//射击
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.MySpirit.Shoot()
		}

		//碰撞检测
		g.CheckObjCollision()
	}

	return nil
}

func (g *Game) CheckObjCollision() {
	//监测我方子弹跟敌机是否碰撞
	if len(g.MySpirit.Bullet) > 0 {
		for key,val := range g.MySpirit.Bullet {
			if len(g.SceneM.Enemy) > 0 {
				for enemyKey,enemyObj := range g.SceneM.Enemy {
					if val.Overlaps(enemyObj) {
						g.SceneM.EnemyHit[enemyKey] = 1 //标记移除
						delete(g.MySpirit.Bullet, key)
					}
				}
			}

		}
	}

	//敌方子弹是否击中我方战机
	if len(g.SceneM.EnemyBullet) > 0 {
		for key,enemyBulletObj := range g.SceneM.EnemyBullet {
			if enemyBulletObj.Overlaps(g.MySpirit.Re) {
				delete(g.SceneM.EnemyBullet, key)
			}
		}
	}

	//敌机撞到我方敌机
	if len(g.SceneM.Enemy) > 0 {
		for key, enemyObj := range g.SceneM.Enemy {
			if enemyObj.Overlaps(g.MySpirit.Re) {
				delete(g.SceneM.Enemy, key)
			}
		}
	}

}





func (g *Game) Draw(screen *ebiten.Image) {
	//ebitenutil.DebugPrint(screen, g.kb.input.msg)
	//screen.Fill(color.RGBA{R: 200, G: 200, B: 200, A: 100})


	if g.SceneM.GetCurrNo() == manager.GameStart {  //游戏初始化状态
		g.SceneM.InitStart(screen)
	} else if g.SceneM.GetCurrNo() == manager.GameChooseSpirit { //选择战机场景
		g.SceneM.Enter()
		g.SceneM.EnterChooseSpirit(screen)
		if g.MySpirit.Img != nil {
			g.SceneM.DrawChooseSpiritEffect(screen, g.MySpirit.Img, g.MySpirit.Name)
		}
	} else if g.SceneM.GetCurrNo() == manager.GameEnter { //进入游戏场景
		g.SceneM.DrawMySpirit(screen, g.MySpirit)
		g.SceneM.DrawMySpiritBullet(screen, g.MySpirit)
		g.SceneM.DrawEnemy(screen)
		g.SceneM.DrawEnemyBullet(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}


//游戏主运行入口
func main() {

	g := Game{
		StateM: manager.NewStateManager(),
		SceneM: manager.NewSceneManager(),
		InputM: manager.NewInputManager(),
		MySpirit:manager.NewMySpirit(),
	}

	//监控我方战机是否射击
	go func() {
		for {
			isShoot := <-g.MySpirit.IsShoot
			if isShoot {
				g.SceneM.AddMySpiritBullet(g.MySpirit)
			}
		}

	}()

	//每3秒自动生成一架敌机
	go func() {
		for {
			time.Sleep(time.Second*3)
			if g.SceneM.GetCurrNo() == manager.GameEnter {
				g.SceneM.AddEnemy()
			}
		}
	}()

	//每2秒敌机自动发射子弹
	go func() {
		for {
			time.Sleep(time.Second*2)
			if g.SceneM.GetCurrNo() == manager.GameEnter {
				g.SceneM.AddEnemyBullet()
			}
		}
	}()


	g.MySpirit.SetXY(200,600)

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("全民射击")
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
