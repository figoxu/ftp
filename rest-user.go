package main

import (
	"fmt"
	"github.com/figoxu/Figo"
	"github.com/go-martini/martini"
	"github.com/monoculum/formam"
	"github.com/quexer/utee"
	"log"
	"net/http"
)

//curl -i -X GET http://127.0.0.1:8080/admin/user
func HandlerUsers(userManager *Figo.Manager, web *utee.Web) (int, string) {
	users := userManager.QueryAll()
	return web.Json(200, users)
}

//curl -i -H "Accept: application/json" http://127.0.0.1:8080/admin/user/1
func HandlerUsersByKey(userManager *Figo.Manager, param martini.Params, web *utee.Web) (int, string) {
	filter := fmt.Sprint("id=", param["id"])
	log.Println("@filter:", filter)
	user := userManager.QueryOne(filter)
	log.Println("@user:", user)
	return web.Json(200, user)
}

//curl -i -X POST http://127.0.0.1:8080/admin/user -d "Account=sample"
func HandlerUsersAdd(userManager *Figo.Manager, r *http.Request, web *utee.Web) (int, string) {
	r.ParseForm()
	user := User{}
	decoder := formam.NewDecoder(&formam.DecoderOptions{})
	err := decoder.Decode(r.Form, &user)
	utee.Chk(err)
	DS.save(&user)
	return web.Json(200, map[string]string{"success": "true"})
}

//curl -i -X PUT http://127.0.0.1:8080/admin/user -d "Id=6&Account=changeName&Basepath=/test"
func HandlerUsersModify(userManager *Figo.Manager, r *http.Request, web *utee.Web) (int, string) {
	r.ParseForm()
	user := User{}
	decoder := formam.NewDecoder(&formam.DecoderOptions{})
	err := decoder.Decode(r.Form, &user)
	utee.Chk(err)
	DS.update(&user)
	return web.Json(200, map[string]string{"success": "true"})
}

//curl -i -X DELETE http://127.0.0.1:8080/admin/user/4
func HandlerUsersDelete(userManager *Figo.Manager, param martini.Params, web *utee.Web) (int, string) {
	user := userManager.QueryOne(fmt.Sprint("id='", param["id"], "'"))
	DS.delete(user.(*User))
	return web.Json(200, map[string]string{"success": "true"})
}

func midUserMng(w http.ResponseWriter, r *http.Request, c martini.Context) {
	log.Println("  midUserMng call")
	userManager := &Figo.Manager{
		Dao: &UserDAO{
			orm: DS.MYSQL,
		},
	}
	c.Map(userManager)
}
