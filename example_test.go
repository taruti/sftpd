package sftpd

import (
	"log"

	"github.com/taruti/sshutil"
)

func ExampleConfig(fs FileSystem) {
	cfg := Config{HostPort: ":2022", FileSystem: fs, LogFunc: log.Println}
	cfg.Init()
	cfg.PasswordCallback = sshutil.CreatePasswordCheck(testUser, testPass)

	// This creates a new host key for each run of the test.
	// Add the sshutil.RSA2048 and sshutil.Save flags if wanted.
	hkey, e := sshutil.KeyLoader{Flags: sshutil.Create}.Load()
	if e != nil {
		log.Println(e)
		return
	}
	cfg.AddHostKey(hkey)

	cfg.RunServer()
}
