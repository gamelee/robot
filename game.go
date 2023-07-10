package robot

import (
	"github.com/gamelee/robot/connect"
	"github.com/gamelee/robot/invoker"
	"github.com/gamelee/robot/ui"
	"log"
	"strings"
)

type Game struct {
	ui     *ui.WebApp
	worker Worker
}

func NewGame(w Worker, opts ...ui.Option) *Game {
	game := &Game{ui: ui.NewApp(opts...), worker: w}
	w.setGame(game)
	return game
}

func snakeString(s string, sep byte) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, sep)
		}
		if d != sep {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToUpper(string(data[:]))
}

// initBind 初始化go方法绑定
func (app *Game) registerWorkerFunction() {

	if err := app.ui.RegStruct(app.worker, func(_, methodName string) (name string, use bool) {
		if !strings.HasPrefix(methodName, "JS") {
			return "", false
		}
		return strings.ToLower(snakeString(methodName[2:], '.')), true
	}); err != nil {
		log.Fatalf("注册函数失败:%v", err)
		return
	}
}

// Jsni 在客户端打印日志
func (app *Game) notifyMsg2JS(m *connect.Message) *invoker.CallRst {
	return app.ui.CallJS(invoker.NewCall("event", "打印日志", m))
}

func (app *Game) Run() error {
	app.registerWorkerFunction()
	return app.ui.Run()
}
