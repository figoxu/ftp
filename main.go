package main

import (
	"github.com/astaxie/beego/config"
	"path/filepath"
	"os"
	"time"
	"io"
	"io/ioutil"
	"log"
	"fmt"
	"github.com/figoxu/utee"
)

var G_CFG config.Configer

type DiskDriver struct {
	cfg  config.Configer
	user string
	path string
}

func (p *DiskDriver) Authenticate(user string, pass string) bool {
	password := p.cfg.String(fmt.Sprint(user, "::pass"))
	if pass != password {
		return false;
	}
	p.user = user
	p.path = p.cfg.String(fmt.Sprint(user, "::path"))
	return true
}

func (p *DiskDriver) RealPath(releativePath string) string {
	return fmt.Sprint(p.path, releativePath)
}

func (p *DiskDriver) Bytes(path string) (bytes int) {
	path = p.RealPath(path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return int(fileInfo.Size())
}

func (p *DiskDriver) ModifiedTime(path string) (time.Time, error) {
	path = p.RealPath(path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Now(), err
	}
	return fileInfo.ModTime(), nil
}

func (p *DiskDriver) ChangeDir(path string) bool {
	return true
}

func (p *DiskDriver) DirContents(path string) ([]os.FileInfo) {
	files := make([]os.FileInfo, 0)
	files = append(files, NewDirItem("."))
	files = append(files, NewDirItem(".."))
	filepath.Walk(p.RealPath(path), func(f string, info os.FileInfo, err error) error {
		if info == nil ||  path== fmt.Sprint("/",info.Name()){
			return nil
		}else if info.IsDir() {
			files = append(files, NewDirItem(info.Name()))
		} else {
			files = append(files, NewFileItem(info.Name(), int(info.Size())))
		}
		return nil
	})
	return files
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

func (p *DiskDriver) GetFile(path string) (string, error) {
	path = p.RealPath(path)
	b, err := ioutil.ReadFile(path)
	return string(b), err
}
func (p *DiskDriver) PutFile(destPath string, data io.Reader) bool {
	destPath = p.RealPath(destPath)
	log.Printf("PutFile: %s", destPath)
	contents, err := ioutil.ReadAll(data)
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

func (factory *DiskDriverFactory) NewDriver() (FTPDriver, error) {
	driver := &DiskDriver{
		cfg:G_CFG,
	}
	return driver, nil
}

// it's alive!
func main() {
	cfg, err := config.NewConfig("ini", "conf.ini")
	utee.Chk(err)
	G_CFG = cfg
	port,err := cfg.Int("ftp::port")
	utee.Chk(err)
	opts := &FTPServerOpts{
		Hostname:cfg.String("ftp::host"),
		Port:port,
		Factory: &DiskDriverFactory{},
	}
	ftpServer := NewFTPServer(opts)
	utee.Chk(ftpServer.ListenAndServe())
}

