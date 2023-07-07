package ui

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

type FileServer struct {
	http.Handler
	port int
	srv  *http.Server
	ln   net.Listener
}

func NewFileServer(fs http.FileSystem, port int) *FileServer {

	ws := new(FileServer)
	ws.Handler = http.FileServer(fs)
	ws.port = port
	return ws
}

func (ws *FileServer) Run() (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("FileServer", "panic", r)
			err = fmt.Errorf("panic:%v", r)
		}
	}()
	return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", ws.port), ws.Handler)
}

func (ws *FileServer) Addr() string {
	return fmt.Sprintf("http://127.0.0.1:%d", ws.port)
}
