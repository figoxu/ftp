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
	orm.RunSyncdb("default", false, true)
	ds := &DataSource{
		MYSQL: orm.NewOrm(),
	}
	return ds
}

func (p *DataSource) searchByAccount(Account string) *User {
	defer Figo.Catch()
	user := &User{}
	err := p.MYSQL.Raw("SELECT id,account,password,basepath from user where Account=?", Account).QueryRow(user)
	utee.Chk(err)
	return user
}

func (p *DataSource) save(user *User) *User {
	_, err := p.MYSQL.Insert(user)
	utee.Chk(err)
	return user
}

func (p *DataSource) update(user *User) *User {
	_, err := p.MYSQL.Update(user)
	utee.Chk(err)
	return user
}

func (p *DataSource) delete(user *User) {
	_, err := p.MYSQL.Delete(user)
	utee.Chk(err)
}

func (p *DataSource) search() []User {
	var users []User
	_, err := p.MYSQL.Raw("SELECT id,account,password,basepath from user").QueryRows(&users)
	utee.Chk(err)
	return users
}
