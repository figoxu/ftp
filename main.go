package main

import (
	"github.com/figoxu/ftp/server"
	"github.com/figoxu/ftp/driver"
	"github.com/astaxie/beego/config"
	"github.com/figoxu/utee"
)

func main() {
	cfg, err := config.NewConfig("ini", "conf.ini")
	utee.Chk(err)
	go server.Monitor()
	fm := driver.NewDefaultFileSystem()
	am := driver.NewDefaultAuthSystem(cfg)
	server.Start(fm, am, false)
}
