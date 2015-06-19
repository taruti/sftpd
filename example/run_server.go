package main

import (
	"crypto/rand"
	"io"
	"log"
	"net"

	"github.com/taruti/sftpd"
	"github.com/taruti/sshutil"
	"golang.org/x/crypto/ssh"
)

// RunServer runs a server serving a given filesystem
func RunServer(hostport string, fs sftpd.FileSystem) {
	e := runServer(hostport, fs)
	if e != nil {
		log.Println("running server errored:", e)
	}
}

func runServer(hostport string, fs sftpd.FileSystem) error {
	config := &ssh.ServerConfig{
		PasswordCallback: sshutil.CreatePasswordCheck(testUser, testPass),
	}

	hkey, e := sshutil.KeyLoader{Flags: sshutil.Create}.Load()
	if e != nil {
		return e
	}

	config.AddHostKey(hkey)

	listener, e := net.Listen("tcp", hostport)
	if e != nil {
		return e
	}

	log.Printf("Listening on %s user %s pass %s\n", hostport, testUser, testPass)

	for {
		conn, e := listener.Accept()
		if e != nil {
			return e
		}
		go HandleConn(conn, config, fs)
	}
}

func HandleConn(conn net.Conn, config *ssh.ServerConfig, fs sftpd.FileSystem) {
	defer conn.Close()
	e := handleConn(conn, config, fs)
	if e != nil {
		log.Println("sftpd connection errored:", e)
	}
}
func handleConn(conn net.Conn, config *ssh.ServerConfig, fs sftpd.FileSystem) error {
	sc, chans, reqs, e := ssh.NewServerConn(conn, config)
	if e != nil {
		return e
	}
	defer sc.Close()

	// The incoming Request channel must be serviced.
	go PrintDiscardRequests(reqs)

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
				case sftpd.IsSftpRequest(req):
					ok = true
					go func() {
						e := sftpd.ServeChannel(channel, fs)
						if e != nil {
							log.Println("sftpd servechannel failed:", e)
						}
					}()
				}
				req.Reply(ok, nil)
			}
		}(requests)

	}
	return nil
}

func PrintDiscardRequests(in <-chan *ssh.Request) {
	for req := range in {
		log.Println("Discarding ssh request", req.Type, *req)
		if req.WantReply {
			req.Reply(false, nil)
		}
	}
}

var testUser = "test"
var testPass = prandAlphaNumeric(16)

const alnum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Pseudo random alpha numeric password generation for this example
func prandAlphaNumeric(n int) []byte {
	bs := make([]byte, n)
	_, e := io.ReadFull(rand.Reader, bs)
	if e != nil {
		panic(e)
	}
	for i, b := range bs {
		bs[i] = alnum[int(b)%len(alnum)]
	}
	return bs
}
