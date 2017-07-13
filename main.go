package main

import (
	"github.com/astaxie/beego/config"
	"github.com/figoxu/utee"
	"github.com/koofr/graval"
	"log"
)

func main() {
	cfg, err := config.NewConfig("ini", "conf.ini")
	utee.Chk(err)

	host := cfg.String("ftp::host")
	port,err := cfg.Int("ftp::port")
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
