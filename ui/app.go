package ui

import (
	"fmt"
	"github.com/gamelee/robot/invoker"
	"github.com/gamelee/robot/ui/lorca"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
)

type WebApp struct {
	*Config
	*invoker.FuncManager
	ui   lorca.UI
	stop <-chan string
	fs   *FileServer
}

func NewApp(opts ...Option) *WebApp {
	cfg := new(Config)
	for i := range opts {
		err := opts[i].Apply(cfg)
		if err != nil {
			log.Panicln(err)
		}
	}
	this := &WebApp{
		Config:      cfg,
		FuncManager: invoker.NewFuncManager(),
		fs:          NewFileServer(http.Dir(cfg.AssetPath), cfg.Port),
	}
	this.init()
	this.stop = make(chan string)
	return this
}

func (wa *WebApp) init() {

	args := make([]string, 0)
	log.Printf("配置信息: %#v", wa.Config)
	args = append(args, "--remote-allow-origins=*")
	args = append(args, "--disable-automation")
	go wa.fs.Run()

	var err error
	if wa.CachePath != "" {
		wa.CachePath, err = filepath.Abs(wa.CachePath)
		if err != nil {
			log.Fatal("init cache path", err.Error())
		}
	}
	wa.ui, err = lorca.New("", wa.CachePath, wa.Width, wa.Height, wa.Bind, args...)
	if err != nil {
		log.Fatal("init ui", err.Error())
	}
}

func (wa *WebApp) Run() (err error) {
	defer func() {
		if err != nil {
			log.Fatalf("运行出错:%s", err.Error())
		}
		if buf := recover(); buf != nil {
			err = fmt.Errorf("程序中☞：%#v", buf)
		}
		if wa.ui != nil {
			wa.ui.Close()
		}
	}()
	wa.ui.Load(wa.fs.Addr())
	defer func() { _ = wa.ui.Close() }()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	reason := "unknown"
	select {
	case <-c:
		reason = "system"
	case <-wa.ui.Done():
		reason = "ui"
	case reason = <-wa.stop:
	}
	log.Println("程序退出，原因：", reason)
	return nil
}

func (wa *WebApp) Bind() error {
	return wa.ui.Bind("GO", wa.invoke)
}

func (wa *WebApp) invoke(call *invoker.Call) *invoker.CallRst {
	log.Println("jscall", call)
	rst := wa.Call(call)
	if tmp, ok := rst.Data.([]interface{}); ok {
		if len(tmp) == 2 {
			rst.Data = tmp[0]
			err, ok := tmp[1].(error)
			if ok && err != nil {
				rst.Code = invoker.CodeFailed
				rst.Message = err.Error()
			}
		}
	}
	return rst
}

func (wa *WebApp) CallJS(call *invoker.Call) *invoker.CallRst {
	rst := new(invoker.CallRst)
	val := wa.ui.Eval(fmt.Sprintf("CallJS(%s)", call))
	err := val.To(rst)
	if err != nil {
		rst.Code = invoker.CodeFailed
		rst.Message = err.Error()
	}
	return rst
}

func (wa *WebApp) RunJS(code string) lorca.Value {
	return wa.ui.Eval(code)
}
