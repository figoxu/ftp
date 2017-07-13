package server

import (
	"bufio"
	"sync"
	"time"
	"fmt"
	"strings"
	"net"
	"io"
	"github.com/figoxu/ftp/driver"
)

var CommandMap map[string]func(*Paradise)
var ConnectionMap map[string]*Paradise
var UpSince int64
var FileManager *driver.FileManager
var AuthManager *driver.AuthManager

type Paradise struct {
	writer        *bufio.Writer
	reader        *bufio.Reader
	theConnection net.Conn
	waiter        sync.WaitGroup
	user          string
	homeDir       string
	path          string
	ip            string
	command       string
	param         string
	total         int64
	buffer        []byte
	cid           string
	connectedAt   int64
	passive       *Passive
	userInfo      map[string]string
	tls           bool
}

func init() {
	UpSince = time.Now().Unix()

	CommandMap = make(map[string]func(*Paradise))

	CommandMap["USER"] = (*Paradise).HandleUser
	CommandMap["PASS"] = (*Paradise).HandlePass
	CommandMap["STOR"] = (*Paradise).HandleStore
	CommandMap["APPE"] = (*Paradise).HandleStore
	CommandMap["STAT"] = (*Paradise).HandleStat

	CommandMap["SYST"] = (*Paradise).HandleSyst
	CommandMap["PWD"] = (*Paradise).HandlePwd
	CommandMap["TYPE"] = (*Paradise).HandleType
	CommandMap["PASV"] = (*Paradise).HandlePassive
	CommandMap["EPSV"] = (*Paradise).HandlePassive
	CommandMap["NLST"] = (*Paradise).HandleList
	CommandMap["LIST"] = (*Paradise).HandleList
	CommandMap["QUIT"] = (*Paradise).HandleQuit
	CommandMap["CWD"] = (*Paradise).HandleCwd
	CommandMap["SIZE"] = (*Paradise).HandleSize
	CommandMap["RETR"] = (*Paradise).HandleRetr
	CommandMap["AUTH"] = (*Paradise).HandleAuth
	CommandMap["PROT"] = (*Paradise).HandleProt
	CommandMap["PBSZ"] = (*Paradise).HandlePbsz

	ConnectionMap = make(map[string]*Paradise)
}

func NewParadise(connection net.Conn, cid string, now int64) *Paradise {
	p := Paradise{}

	p.writer = bufio.NewWriter(connection)
	p.reader = bufio.NewReader(connection)
	p.path = "/"
	p.theConnection = connection
	p.ip = connection.RemoteAddr().String()
	p.cid = cid
	p.connectedAt = now
	p.userInfo = make(map[string]string)
	p.userInfo["path"] = "/"
	return &p
}

func (p *Paradise) lastPassive() *Passive {
	if p.passive == nil {
		return nil
	}
	p.passive.command = p.command
	p.passive.param = p.param
	return p.passive
}

func (p *Paradise) HandleCommands() {
	//fmt.Println(p.id, " Got client on: ", p.ip)
	p.writeMessage(220, "Welcome to Paradise")
	for {
		line, err := p.reader.ReadString('\n')
		if err != nil {
			delete(ConnectionMap, p.cid)
			//fmt.Println(p.id, " end ", len(ConnectionMap))
			if err == io.EOF {
				//continue
			}
			break
		}
		command, param := parseLine(line)
		p.command = command
		p.param = param

		fn := CommandMap[command]
		if fn == nil {
			p.writeMessage(550, "not allowed")
		} else {
			fn(p)
		}
	}
}

func (p *Paradise) writeMessage(code int, message string) {
	line := fmt.Sprintf("%d %s\r\n", code, message)
	p.writer.WriteString(line)
	p.writer.Flush()
}

func parseLine(line string) (string, string) {
	params := strings.SplitN(strings.Trim(line, "\r\n"), " ", 2)
	if len(params) == 1 {
		return params[0], ""
	}
	return params[0], strings.TrimSpace(params[1])
}
