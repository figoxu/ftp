package driver

import (
	"encoding/json"
	"github.com/figoxu/utee"
	"log"
	"fmt"
	"time"
	"path/filepath"
	"os"
)

type FileSystem interface {
	GetFiles(userInfo *map[string]string) ([]map[string]string, error)
}

type FileManager struct {
	FileSystem
}

type DefaultFileSystem struct {
}

func (dfs DefaultFileSystem) GetFiles(userInfo *map[string]string) ([]map[string]string, error) {
	files := make([]map[string]string, 0)

	b,e:=json.Marshal(userInfo)
	utee.Chk(e)
	log.Println(string(b))


	appendFile:=func(name string,size int64,isdir bool,modTime int64){
		file := make(map[string]string)
		file["size"] = fmt.Sprint(size)
		file["name"] = name
		file["modTime"] = fmt.Sprint(modTime)
		if isdir {
			file["isDir"] = "true"
		}
		files = append(files, file)
	}
	appendFile(".",0,true,time.Now().Unix())
	appendFile("..",0,true,time.Now().Unix())
	u:=*userInfo
	path:=fmt.Sprint(u["basepath"],u["path"])
	log.Println("@path:",path)
	filepath.Walk(path, func(f string, info os.FileInfo, err error) error {
		if info == nil ||  u["path"]== fmt.Sprint("/",info.Name()){
			return nil
		} else {
			appendFile(info.Name(),info.Size(),info.IsDir(),info.ModTime().Unix())
		}
		return nil
	})
	return files, nil
}

func NewDefaultFileSystem() *FileManager {
	fm := FileManager{}
	fm.FileSystem = DefaultFileSystem{}
	return &fm
}
