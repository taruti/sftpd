package main

import (
	"log"

	"github.com/taruti/sftpd"
	"github.com/taruti/sshutil"
)

// RunServerHighLevel is an example how to use the low level API
func RunServerHighLevel(hostport string, fs sftpd.FileSystem) {
	cfg := sftpd.Config{HostPort: hostport, FileSystem: fs, LogFunc: log.Println}
	cfg.PasswordCallback = sshutil.CreatePasswordCheck(testUser, testPass)

	// Add the sshutil.RSA2048 and sshutil.Save flags if needed for the server in question...
	hkey, e := sshutil.KeyLoader{Flags: sshutil.Create}.Load()
	if e != nil {
		log.Println(e)
		return
	}
	cfg.AddHostKey(hkey)

	log.Printf("Listening on %s user %s pass %s\n", hostport, testUser, testPass)
	cfg.RunServer()
}
