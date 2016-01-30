package sftpd

import (
	"net"

	"golang.org/x/crypto/ssh"
)

// Config is the configuration struct for the high level API.
type Config struct {
	// ServerConfig should be initialized properly with
	// e.g. PasswordCallback and AddHostKey
	ssh.ServerConfig
	// HostPort specifies specifies [host]:port to listen on.
	// e.g. ":2022" or "127.0.0.1:2023".
	HostPort string
	// LogFunction is used to log errors.
	// e.g. log.Println has the right type.
	LogFunc func(v ...interface{})
	// FileSystem contains the FileSystem used for this server.
	FileSystem FileSystem

	readyChan chan error
	connChan chan net.Listener
}

// Init inits a Config.
func (c *Config) Init() {
	c.readyChan = make(chan error, 1)
	c.connChan = make(chan net.Listener, 1)
}

// RunServer runs the server using the high level API.
func (c *Config) RunServer() error {
	if c.LogFunc == nil {
		c.LogFunc = func(...interface{}) {}
	}
	e := runServer(c)
	if e != nil {
		c.LogFunc("sftpd server failed:", e)
	}
	return e
}

func runServer(c *Config) error {
	listener, e := net.Listen("tcp", c.HostPort)
	c.readyChan <- e
	close(c.readyChan)
	c.connChan <- listener
	close(c.connChan)
	if e != nil {
		return e
	}

	for {
		conn, e := listener.Accept()
		if e != nil {
			return e
		}
		go handleConn(conn, c)
	}
}

func handleConn(conn net.Conn, config *Config) {
	defer conn.Close()
	e := doHandleConn(conn, config)
	if e != nil {
		config.LogFunc("sftpd connection error:", e)
	}
}

func doHandleConn(conn net.Conn, config *Config) error {
	sc, chans, reqs, e := ssh.NewServerConn(conn, &config.ServerConfig)
	if e != nil {
		return e
	}
	defer sc.Close()

	// The incoming Request channel must be serviced.
	go printDiscardRequests(config, reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			return err
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok := false
				switch {
				case IsSftpRequest(req):
					ok = true
					go func() {
						e := ServeChannel(channel, config.FileSystem)
						if e != nil {
							config.LogFunc("sftpd servechannel failed:", e)
						}
					}()
				}
				req.Reply(ok, nil)
			}
		}(requests)

	}
	return nil
}

func printDiscardRequests(c *Config, in <-chan *ssh.Request) {
	for req := range in {
		c.LogFunc("sftpd discarding ssh request", req.Type, *req)
		if req.WantReply {
			req.Reply(false, nil)
		}
	}
}

// BlockTillReady will block till the Config is ready to accept connections.
// Returns an error if listening failed. Can be called in a concurrent fashion.
// This is new API - make sure Init is called on the Config before using it.
func (c *Config) BlockTillReady() error {
	err, _ := <-c.readyChan
	return err
}

// Close closes the server assosiated with this config. Can be called in a concurrent
// fashion.
// This is new API - make sure Init is called on the Config before using it.
func (c *Config) Close() error {
	for c := range c.connChan {
		c.Close()
	}
	return nil
}
