package main

import (
	"github.com/astaxie/beego/orm"
	"github.com/figoxu/Figo"
	"github.com/quexer/utee"
	"log"
	"strconv"
)

type User struct {
	Id       int
	Account  string
	Password string
	Basepath string
}

type DataSource struct {
	MYSQL orm.Ormer
}

var DS = NewDataSource()

func NewDataSource() *DataSource {
	idleStr := utee.Env("MYSQL_IDLE", false, true)
	activeStr := utee.Env("MYSQL_ACTIVE", false, true)

	log.Println("@idle:", idleStr, "  @active:", activeStr)
	idle, err := strconv.Atoi(idleStr)
	utee.Chk(err)
	active, err := strconv.Atoi(activeStr)
	utee.Chk(err)
	conf := Figo.MysqlConf{
		User:       utee.Env("MYSQL_USER", false),
		Pwd:        utee.Env("MYSQL_PWD", false),
		Host:       utee.Env("MYSQL_HOST", false),
		Port:       utee.Env("MYSQL_PORT", false),
		Name:       utee.Env("MYSQL_NAME", false),
		ConnIdle:   idle,
		ConnActive: active,
	}
	conf.Conf(new(User))
	DataReg("user", User{}, []User{})
	orm.RunSyncdb("default", false, true)
	ds := &DataSource{
		MYSQL: orm.NewOrm(),
	}
	return ds
}

type DataInstance struct {
	Item  interface{}
	Items interface{}
}

var DataInstanceMap = make(map[string]*DataInstance)

func DataReg(name string, item interface{}, items interface{}) {
	DataInstanceMap[name] = &DataInstance{
		Item:  item,
		Items: items,
	}
}

func GetDao(name string, orm orm.Ormer) *DAO {
	v := DataInstanceMap[name]
	if v == nil {
		return nil
	}
	return &DAO{
		orm:   orm,
		Item:  &v.Item,
		Items: &v.Items,
	}
}
