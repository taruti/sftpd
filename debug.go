// +build debug

package sftpd

import "log"

var debug func(...interface{}) = log.Println
var debugf func(string, ...interface{}) = log.Printf

func (b ssh_fxp) String() string {
	return ssh_fxp_map[b]
}

var ssh_fxp_map = map[ssh_fxp]string{
	ssh_FXP_VERSION:        `ssh_FXP_VERSION`,
	ssh_FXP_OPEN:           `ssh_FXP_OPEN`,
	ssh_FXP_CLOSE:          `ssh_FXP_CLOSE`,
	ssh_FXP_READ:           `ssh_FXP_READ`,
	ssh_FXP_WRITE:          `ssh_FXP_WRITE`,
	ssh_FXP_LSTAT:          `ssh_FXP_LSTAT`,
	ssh_FXP_FSTAT:          `ssh_FXP_FSTAT`,
	ssh_FXP_SETSTAT:        `ssh_FXP_SETSTAT`,
	ssh_FXP_FSETSTAT:       `ssh_FXP_FSETSTAT`,
	ssh_FXP_OPENDIR:        `ssh_FXP_OPENDIR`,
	ssh_FXP_READDIR:        `ssh_FXP_READDIR`,
	ssh_FXP_REMOVE:         `ssh_FXP_REMOVE`,
	ssh_FXP_MKDIR:          `ssh_FXP_MKDIR`,
	ssh_FXP_RMDIR:          `ssh_FXP_RMDIR`,
	ssh_FXP_REALPATH:       `ssh_FXP_REALPATH`,
	ssh_FXP_STAT:           `ssh_FXP_STAT`,
	ssh_FXP_RENAME:         `ssh_FXP_RENAME`,
	ssh_FXP_READLINK:       `ssh_FXP_READLINK`,
	ssh_FXP_SYMLINK:        `ssh_FXP_SYMLINK`,
	ssh_FXP_STATUS:         `ssh_FXP_STATUS`,
	ssh_FXP_HANDLE:         `ssh_FXP_HANDLE`,
	ssh_FXP_DATA:           `ssh_FXP_DATA`,
	ssh_FXP_NAME:           `ssh_FXP_NAME`,
	ssh_FXP_ATTRS:          `ssh_FXP_ATTRS`,
	ssh_FXP_EXTENDED:       `ssh_FXP_EXTENDED`,
	ssh_FXP_EXTENDED_REPLY: `ssh_FXP_EXTENDED_REPLY`,
}

func (b ssh_fx) String() string {
	return ssh_fx_map[b]
}

var ssh_fx_map = map[ssh_fx]string{
	ssh_FX_OK:                `ssh_FX_OK`,
	ssh_FX_EOF:               `ssh_FX_EOF`,
	ssh_FX_NO_SUCH_FILE:      `ssh_FX_NO_SUCH_FILE`,
	ssh_FX_PERMISSION_DENIED: `ssh_FX_PERMISSION_DENIED`,
	ssh_FX_FAILURE:           `ssh_FX_FAILURE`,
	ssh_FX_BAD_MESSAGE:       `ssh_FX_BAD_MESSAGE`,
	ssh_FX_NO_CONNECTION:     `ssh_FX_NO_CONNECTION`,
	ssh_FX_CONNECTION_LOST:   `ssh_FX_CONNECTION_LOST`,
	ssh_FX_OP_UNSUPPORTED:    `ssh_FX_OP_UNSUPPORTED`,
}
