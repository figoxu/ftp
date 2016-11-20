package main

import (
	"fmt"
	"github.com/quexer/utee"
	"os"
)

func main() {
	port, rootPath, ftpName := 21, "/", "Figo's FTP"
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
