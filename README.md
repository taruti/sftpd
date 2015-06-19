# sftpd - SFTP server library in Go

# License: MIT

# [![GoDoc](https://godoc.org/github.com/taruti/sftpd?status.png)](http://godoc.org/github.com/taruti/sftpd)

# FAQ

## ssh: no common algorithms

The client and the server cannot agree on algorithms. Typically this
is caused by an ECDSA host key. If using sshutil add the
``sshutil.RSA2048`` flag.

# TODO
+ [ ] Renames
+ [ ] Symlink creation
