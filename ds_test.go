package main

import (
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
