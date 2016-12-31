package main

import (
	"github.com/figoxu/Figo"
	"log"
	"testing"
)

func getTestDs() *DataSource {
	return DS
}

func TestSave(t *testing.T) {
	ds := getTestDs()
	user := &User{
		Account:  "figo2",
		Password: "123456",
		Basepath: "/home/figo",
	}
	ds.save(user)
}

func TestSearchByAccount(t *testing.T) {
	ds := getTestDs()
	user := ds.searchByAccount("figo")
	log.Println(user)

	user.Password = "changePwd"
	ds.update(user)
}

func TestDelete(t *testing.T) {
	ds := getTestDs()
	user := ds.searchByAccount("figo2")
	ds.delete(user)

}

func TestSearch(t *testing.T) {
	ds := getTestDs()
	data := ds.search()
	log.Println(data)
}

func TestUserDAO(t *testing.T) {
	userManager := &Figo.Manager{
		Dao: &UserDAO{
			orm: getTestDs().MYSQL,
		},
	}
	log.Println("CountALl is :", userManager.CountAll())
	log.Println("QueryALl is :", userManager.QueryAll())
	log.Println("CountFilter is :", userManager.CountFilter("Account='andy'"))
	log.Println("QueryFilter is :", userManager.QueryFilter("Account='andy'"))
	log.Println("QueryFilterPaging is :", userManager.QueryFilterPaging("Account='andy'", 0, 5))

}
