package sftpd

import (
	"io"
	"os"
	"time"
)

type Attr struct {
	Flags        uint32
	Size         uint64
	Uid, Gid     uint32
	User, Group  string
	Mode         os.FileMode
	ATime, MTime time.Time
	Extended     []string
}

type NamedAttr struct {
	Name string
	Attr
}

const (
	ATTR_SIZE    = ssh_FILEXFER_ATTR_SIZE
	ATTR_UIDGID  = ssh_FILEXFER_ATTR_UIDGID
	ATTR_MODE    = ssh_FILEXFER_ATTR_PERMISSIONS
	ATTR_TIME    = ssh_FILEXFER_ATTR_ACMODTIME
	MODE_REGULAR = os.FileMode(0)
	MODE_DIR     = os.ModeDir
)

type Dir interface {
	io.Closer
	Readdir(count int) ([]NamedAttr, error)
}

type File interface {
	io.Closer
	io.ReaderAt
	io.WriterAt
	FStat() (*Attr, error)
	FSetStat(*Attr) error
}

type FileSystem interface {
	OpenFile(name string, flags uint32, attr *Attr) (File, error)
	OpenDir(name string) (Dir, error)
	Remove(name string) error
	Rename(old string, new string, flags uint32) error
	Mkdir(name string, attr *Attr) error
	Rmdir(name string) error
	Stat(name string, islstat bool) (*Attr, error)
	SetStat(name string, attr *Attr) error
	ReadLink(path string) (string, error)
	CreateLink(path string, target string, flags uint32) error
	RealPath(path string) (string, error)
}

// FillFrom fills an Attr from a os.FileInfo
func (a *Attr) FillFrom(fi os.FileInfo) error {
	*a = Attr{}
	a.Flags = ATTR_SIZE | ATTR_MODE
	a.Size = uint64(fi.Size())
	a.Mode = fi.Mode()
	a.MTime = fi.ModTime()
	return nil
}

func fileModeToSftp(m os.FileMode) uint32 {
	var raw = uint32(m.Perm())
	switch {
	case m.IsDir():
		raw |= 0040000
	case m.IsRegular():
		raw |= 0100000
	}
	return raw
}

func sftpToFileMode(raw uint32) os.FileMode {
	var m = os.FileMode(raw & 0777)
	switch {
	case raw&0040000 != 0:
		m |= os.ModeDir
	case raw&0100000 != 0:
		// regular
	}
	return m
}
