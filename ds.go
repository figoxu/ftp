package main

import (
	"github.com/astaxie/beego/orm"
	"github.com/quexer/utee"
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

func (p *DataSource) searchByAccount(Account string) *User {
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
