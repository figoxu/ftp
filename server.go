package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"strconv"
	"strings"
)

// serverOpts contains parameters for server.NewServer()
type ServerOpts struct {
	// The factory that will be used to create a new FTPDriver instance for
	// each client connection. This is a mandatory option.
	Factory DriverFactory

	//	Auth Auth

	// Server Name, Default is Go Ftp Server
	Name string

	// The hostname that the FTP server should listen on. Optional, defaults to
	// "::", which means all hostnames on ipv4 and ipv6.
	Hostname string

	// Public IP of the server
	PublicIp string

	// Passive ports
	PassivePorts string

	// The port that the FTP should listen on. Optional, defaults to 3000. In
	// a production environment you will probably want to change this to 21.
	Port int

	// use tls, default is false
	TLS bool

	// if tls used, cert file is required
	CertFile string

	// if tls used, key file is required
	KeyFile string

	// If ture TLS is used in RFC4217 mode
	ExplicitFTPS bool

	WelcomeMessage string
}

type Server struct {
	*ServerOpts
	name          string
	listenTo      string
	driverFactory DriverFactory
	logger        *Logger
	listener      net.Listener
	tlsConfig     *tls.Config
	publicIp      string
}

func serverOptsWithDefaults(opts *ServerOpts) *ServerOpts {
	var newOpts ServerOpts
	if opts == nil {
		opts = &ServerOpts{}
	}
	if opts.Hostname == "" {
		newOpts.Hostname = "::"
	} else {
		newOpts.Hostname = opts.Hostname
	}
	if opts.Port == 0 {
		newOpts.Port = 3000
	} else {
		newOpts.Port = opts.Port
	}
	newOpts.Factory = opts.Factory
	if opts.Name == "" {
		newOpts.Name = "Go FTP Server"
	} else {
		newOpts.Name = opts.Name
	}

	if opts.WelcomeMessage == "" {
		newOpts.WelcomeMessage = defaultWelcomeMessage
	} else {
		newOpts.WelcomeMessage = opts.WelcomeMessage
	}

	newOpts.TLS = opts.TLS
	newOpts.KeyFile = opts.KeyFile
	newOpts.CertFile = opts.CertFile
	newOpts.ExplicitFTPS = opts.ExplicitFTPS

	newOpts.PublicIp = opts.PublicIp
	newOpts.PassivePorts = opts.PassivePorts

	return &newOpts
}

func NewServer(opts *ServerOpts) *Server {
	opts = serverOptsWithDefaults(opts)
	s := new(Server)
	s.ServerOpts = opts
	s.listenTo = buildTCPString(opts.Hostname, opts.Port)
	s.name = opts.Name
	s.driverFactory = opts.Factory
	s.logger = newLogger("")
	return s
}

func (server *Server) newConn(tcpConn net.Conn, driver Driver) *Conn {
	c := new(Conn)
	c.namePrefix = "/"
	c.conn = tcpConn
	c.controlReader = bufio.NewReader(tcpConn)
	c.controlWriter = bufio.NewWriter(tcpConn)
	c.driver = driver
	c.server = server
	c.sessionID = newSessionID()
	c.logger = newLogger(c.sessionID)
	c.tlsConfig = server.tlsConfig
	driver.Init(c, nil)
	return c
}

func simpleTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	config := &tls.Config{}
	if config.NextProtos == nil {
		config.NextProtos = []string{"ftp"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (server *Server) ListenAndServe() error {
	var listener net.Listener
	var err error

	if server.ServerOpts.TLS {
		server.tlsConfig, err = simpleTLSConfig(server.CertFile, server.KeyFile)
		if err != nil {
			return err
		}

		if server.ServerOpts.ExplicitFTPS {
			listener, err = net.Listen("tcp", server.listenTo)
		} else {
			listener, err = tls.Listen("tcp", server.listenTo, server.tlsConfig)
		}
	} else {
		listener, err = net.Listen("tcp", server.listenTo)
	}
	if err != nil {
		return err
	}

	server.logger.Printf("%s listening on %d", server.Name, server.Port)

	server.listener = listener
	for {
		tcpConn, err := server.listener.Accept()
		if err != nil {
			server.logger.Printf("listening error: %v", err)
			break
		}
		driver, err := server.driverFactory.NewDriver()
		if err != nil {
			server.logger.Printf("Error creating driver, aborting client connection: %v", err)
			tcpConn.Close()
		} else {
			log.Println("===>>>>>> @tcpConn newConn")
			ftpConn := server.newConn(tcpConn, driver)
			go ftpConn.Serve()
		}
	}
	return nil
}

func (server *Server) Shutdown() error {
	if server.listener != nil {
		return server.listener.Close()
	}
	return nil
}

func buildTCPString(hostname string, port int) (result string) {
	if strings.Contains(hostname, ":") {
		// ipv6
		if port == 0 {
			result = "[" + hostname + "]"
		} else {
			result = "[" + hostname + "]:" + strconv.Itoa(port)
		}
	} else {
		// ipv4
		if port == 0 {
			result = hostname
		} else {
			result = hostname + ":" + strconv.Itoa(port)
		}
	}
	return
}
