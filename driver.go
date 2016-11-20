package main

import "io"

type DriverFactory interface {
	NewDriver() (Driver, error)
}

type Driver interface {
	Init(*Conn)
	Stat(string) (IFileInfo, error)
	ChangeDir(string) error
	ListDir(string, func(IFileInfo) error) error
	DeleteDir(string) error
	DeleteFile(string) error
	Rename(string, string) error
	MakeDir(string) error
	GetFile(string, int64) (int64, io.ReadCloser, error)
	PutFile(string, io.Reader, bool) (int64, error)
}
