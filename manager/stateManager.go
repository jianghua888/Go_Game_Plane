package manager

const INIT_STATE = 0 //未开始
const RUNNING = 1 //运行中
const STOP = 2 //暂停
const END = 3 //结束

type StateManager struct {
	state int8 //游戏状态
}

func NewStateManager() *StateManager {
	return &StateManager{ state: INIT_STATE}
}

func (gm StateManager) Start() {
	gm.state = RUNNING
}

func (gm StateManager) Stop() {
	gm.state = STOP
}

func (gm StateManager) End() {
	gm.state = END
}

func (gm StateManager) GetGameState() int8{
	return gm.state
}






