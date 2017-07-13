package main

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/koofr/graval"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var G_CFG config.Configer

type DiskDriver struct {
	cfg  config.Configer
	user string
	path string
}

func (p *DiskDriver) Authenticate(username string, password string) bool {
	pass := p.cfg.String(fmt.Sprint(username, "::pass"))
	if password != pass {
		return false
	}
	p.user = username
	p.path = p.cfg.String(fmt.Sprint(username, "::path"))
	return true
}

func (p *DiskDriver) RealPath(releativePath string) string {
	return fmt.Sprint(p.path, releativePath)
}

func (p *DiskDriver) Bytes(path string) int64 {
	path = p.RealPath(path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

func (p *DiskDriver) ModifiedTime(path string) (time.Time, bool) {
	path = p.RealPath(path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Now(), false
	}
	return fileInfo.ModTime(), true
}

func (p *DiskDriver) ChangeDir(path string) bool {
	return true
}

func (p *DiskDriver) DirContents(path string) ([]os.FileInfo, bool) {
	fileInfos := make([]os.FileInfo, 0)
	realPath := p.RealPath(path)
	filepath.Walk(realPath, func(f string, info os.FileInfo, err error) error {
		if info == nil || path == fmt.Sprint("/", info.Name() ) {
			return nil
		}else if f== realPath ||fmt.Sprint(f,"/")== realPath {
			return nil
		}
		fileInfos = append(fileInfos, info)
		return nil
	})
	return fileInfos, true
}

func (p *DiskDriver) DeleteDir(path string) bool {
	path = p.RealPath(path)
	log.Printf("DeleteDir: %s", path)
	f, err := os.Lstat(path)
	if err != nil || !f.IsDir() {
		return false
	}
	os.Remove(path)
	return true
}

func (p *DiskDriver) DeleteFile(path string) bool {
	path = p.RealPath(path)
	log.Printf("DeleteFile: %s", path)
	f, err := os.Lstat(path)
	if err != nil || f.IsDir() {
		return false
	}
	os.Remove(path)
	return true
}

func (p *DiskDriver) Rename(fromPath string, toPath string) bool {
	oldPath := p.RealPath(fromPath)
	newPath := p.RealPath(toPath)
	return os.Rename(oldPath, newPath) != nil
}

func (p *DiskDriver) MakeDir(path string) bool {
	return os.Mkdir(p.RealPath(path), os.ModePerm) != nil
}

func (p *DiskDriver) GetFile(path string, position int64) (io.ReadCloser, bool) {
	path = p.RealPath(path)
	log.Println("GETFILE 00001 @path:", path)
	if f, err := os.Stat(path); err == nil && !f.IsDir() {
		log.Println("GETFILE 00002")
		b, _ := ioutil.ReadFile(path)
		return ioutil.NopCloser(bytes.NewReader(b[position:])), true
	} else {
		log.Println("GETFILE 00003")
		return nil, false
	}

}
func (p *DiskDriver) PutFile(path string, reader io.Reader) bool {
	destPath := p.RealPath(path)
	log.Printf("PutFile: %s", destPath)
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return false
	}
	var existFlag bool
	if fileInfo, err := os.Stat(destPath); err == nil {
		existFlag = true
		if fileInfo.IsDir() {
			return false
		}
	} else if os.IsNotExist(err) {
		existFlag = false
	}
	if existFlag {
		os.Remove(destPath)
	}
	ioutil.WriteFile(destPath, contents, os.ModePerm)

	return true
}

// graval requires a factory that will create a new driver instance for each
// client connection. Generally the factory will be fairly minimal. This is
// a good place to read any required config for your driver.
type DiskDriverFactory struct{}

func (factory *DiskDriverFactory) NewDriver() (graval.FTPDriver, error) {
	driver := &DiskDriver{
		cfg: G_CFG,
	}
	return driver, nil
}
