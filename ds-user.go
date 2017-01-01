package main

import (
	"github.com/astaxie/beego/orm"
	"github.com/figoxu/Figo"
	"github.com/quexer/utee"
)

func (p *DataSource) searchByAccount(Account string) *User {
	defer Figo.Catch()
	user := &User{}
	err := p.MYSQL.Raw("SELECT id,account,password,basepath from user where Account=?", Account).QueryRow(user)
	utee.Chk(err)
	return user
}

func (p *DataSource) save(user *User) *User {
	defer Figo.Catch()
	_, err := p.MYSQL.Insert(user)
	utee.Chk(err)
	return user
}

func (p *DataSource) update(user *User) *User {
	defer Figo.Catch()
	_, err := p.MYSQL.Update(user)
	utee.Chk(err)
	return user
}

func (p *DataSource) delete(user *User) {
	defer Figo.Catch()
	_, err := p.MYSQL.Delete(user)
	utee.Chk(err)
}

func (p *DataSource) search() []User {
	defer Figo.Catch()
	var users []User
	_, err := p.MYSQL.Raw("SELECT id,account,password,basepath from user").QueryRows(&users)
	utee.Chk(err)
	return users
}

type DAO struct {
	orm   orm.Ormer
	Item  interface{}
	Items interface{}
}

func (p *DAO) GetORM() orm.Ormer {
	return p.orm
}
func (p *DAO) GetItemContainer() interface{} {
	return Figo.Clone(p.Item)
}
func (p *DAO) GetItemsContainer() interface{} {
	return Figo.Clone(p.Items)
}
