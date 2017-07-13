package main

import (
	"github.com/astaxie/beego/config"
	"github.com/figoxu/utee"
	"github.com/koofr/graval"
	"log"
)

func main() {
	host := "0.0.0.0"
	port := 8021

	cfg, err := config.NewConfig("ini", "conf.ini")
	utee.Chk(err)
	G_CFG = cfg
	server := graval.NewFTPServer(&graval.FTPServerOpts{
		ServerName: "Example FTP server",
		Factory:    &DiskDriverFactory{},
		Hostname:   host,
		Port:       port,
		PassiveOpts: &graval.PassiveOpts{
			ListenAddress: host,
			NatAddress:    host,
			PassivePorts: &graval.PassivePorts{
				Low:  42000,
				High: 45000,
			},
		},
	})
	log.Fatal(server.ListenAndServe())
}
