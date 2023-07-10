package main

import (
	"github.com/gamelee/robot"
	"github.com/gamelee/robot/app/nodes"
	"github.com/gamelee/robot/ui"
	"log"
	"os"
)

const (
	GameName  = "多多自走棋"
	Height    = 720
	Width     = 1280
	Port      = 52341
	CachePath = "./cache/browser"
	TreeFile  = "./cache/behavior_tree.json"
	LogFile   = "./cache/app.log"
	AssetPath = "./static"
)

func main() {
	nodes.InitNodes()
	panicLogError(makeDirIfNotExists(CachePath))

	log.SetFlags(log.Lshortfile | log.Ltime)
	fd, err := os.OpenFile(LogFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	panicLogError(err)
	log.SetOutput(fd)
	worker := NewDDChessWorker()
	game := robot.NewGame(
		worker,
		ui.WithName(GameName),
		ui.WithCachePath(CachePath),
		ui.WithPort(Port),
		ui.WithSize{Width, Height},
		ui.WithAssetPath(AssetPath),
	)
	go worker.msgReceiver()
	panicLogError(game.Run())
}

func panicLogError(err error) {
	if err == nil {
		return
	}
	log.Fatalln(err)
}

func makeDirIfNotExists(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return err
}
