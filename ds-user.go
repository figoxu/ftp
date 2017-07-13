package main

import (
	"github.com/astaxie/beego/orm"
	"github.com/figoxu/Figo"
	"github.com/quexer/utee"
	"log"
)

func (p *DataSource) searchByAccount(Account string) *User {
	defer Figo.Catch()
	user := &User{}
	err := p.MYSQL.Raw("SELECT id,account,password,basepath from user where account=?", Account).QueryRow(user)
	if err!= nil {
		return nil
	}
	log.Println("@user : ",user.Password," basePath:",user.Basepath)
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

type UserDAO struct {
	orm orm.Ormer
}

func (p *UserDAO) GetORM() orm.Ormer {
	return p.orm
}
func (p *UserDAO) GetItemContainer() interface{} {
	return &User{}
}
func (p *UserDAO) GetItemsContainer() interface{} {
	return &[]User{}
}
