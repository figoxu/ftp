package driver

import (
	"encoding/json"
	"github.com/figoxu/utee"
	"log"
	"github.com/astaxie/beego/config"
	"fmt"
)

type AuthSystem interface {
	CheckUser(user, pass string, userInfo *map[string]string) bool
}

type AuthManager struct {
	AuthSystem
}

type DefaultAuthSystem struct {
	cfg config.Configer
}

func (p DefaultAuthSystem) CheckUser(user, pass string, userInfo *map[string]string) bool {
	b,e:=json.Marshal(userInfo)
	utee.Chk(e)
	log.Println(string(b))
	password := p.cfg.String(fmt.Sprint(user, "::pass"))
	if pass != password || password=="" {
		return false;
	}
	info := *userInfo
	info["user"] = user
	info["basepath"]= p.cfg.String(fmt.Sprint(user, "::path"))
	return true
}

func NewDefaultAuthSystem(cfg config.Configer) *AuthManager {
	am := AuthManager{}
	am.AuthSystem = DefaultAuthSystem{
		cfg:cfg,
	}
	return &am
}
