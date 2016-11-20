package main

type Auth interface {
	CheckPasswd(string, string) (bool, error)
}
