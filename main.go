package main

import (
	"github.com/figoxu/ftp/server"
	"github.com/figoxu/ftp/driver"
)

func main() {
	go server.Monitor()
	fm := driver.NewDefaultFileSystem()
	am := driver.NewDefaultAuthSystem()
	server.Start(fm, am, false)
}
