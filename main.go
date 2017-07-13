package main

import (
	"fmt"
	"github.com/figoxu/Figo"
	"github.com/go-martini/martini"
	"github.com/quexer/utee"
	"log"
	"net/http"
	"os"
)

func main() {
	go startHttp()
	port, rootPath, ftpName := 2121, "", "Figo's FTP"
	_, err := os.Lstat(rootPath)
	if os.IsNotExist(err) {
		os.MkdirAll(rootPath, os.ModePerm)
	} else if err != nil {
		fmt.Println(err)
		return
	}
	factory := &FileDriverFactory{
		RootPath: rootPath,
	}
	opt := &ServerOpts{
		Name:    ftpName,
		Factory: factory,
		Port:    port,
		//		Auth:    auth,
	}
	utee.Chk(NewServer(opt).ListenAndServe())
}

func startHttp() {
	m := Figo.NewMartini(1, 100, "")
	m.Use(func(w http.ResponseWriter, c martini.Context) {
		web := &utee.Web{W: w}
		c.Map(web)
	})
	m.Group("/admin", func(r martini.Router) {
		r.Get("/user", HandlerUsers)
		r.Get("/user/:id", HandlerUsersByKey)
		r.Post("/user", HandlerUsersAdd)
		r.Put("/user", HandlerUsersModify)
		r.Delete("/user/:id", HandlerUsersDelete)
	}, midUserMng)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
