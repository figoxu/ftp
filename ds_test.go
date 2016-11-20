package main

import (
	"github.com/astaxie/beego/orm"
	"github.com/figoxu/Figo"
	"log"
	"testing"
)

func getTestDs() *DataSource {
	conf := Figo.MysqlConf{
		User:       "root",
		Pwd:        "xxaaaeee",
		Host:       "127.0.0.1",
		Port:       "3306",
		Name:       "figo_research",
		ConnIdle:   1,
		ConnActive: 1,
	}
	conf.Conf(new(User))
	orm.RunSyncdb("default", false, true)
	ds := &DataSource{
		MYSQL: orm.NewOrm(),
	}
	return ds
}

func TestSave(t *testing.T) {
	ds := getTestDs()
	user := &User{
		Account:  "figo",
		Password: "123456",
		Basepath: "/home/figo",
	}
	ds.save(user)
}

func TestSearchByAccount(t *testing.T) {
	ds := getTestDs()
	user := ds.searchByAccount("figo")
	log.Println(user)
}
