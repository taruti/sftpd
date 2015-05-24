package sftpd

import "io"

type Attr struct {
	Flags        uint32
	Size         uint64
	Uid, Gid     uint32
	Mode         uint32
	ATime, MTime uint32
	Extended     []string
}

type NamedAttr struct {
	Name string
	Attr
}

const (
	ATTR_SIZE   = ssh_FILEXFER_ATTR_SIZE
	ATTR_UIDGID = ssh_FILEXFER_ATTR_UIDGID
	ATTR_MODE   = ssh_FILEXFER_ATTR_PERMISSIONS
	ATTR_TIME   = ssh_FILEXFER_ATTR_ACMODTIME
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
