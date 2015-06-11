package main

// Read only file system using the file system

import (
	"errors"
	"os"
	"strings"

	"github.com/taruti/sftpd"
)

type rfs struct {
	sftpd.EmptyFS
}

type rdir struct {
	d *os.File
}

type rfile struct {
	sftpd.EmptyFile
	f *os.File
}

func (rf rfile) Close() error                             { return rf.f.Close() }
func (rf rfile) ReadAt(bs []byte, pos int64) (int, error) { return rf.f.ReadAt(bs, pos) }

func (d rdir) Readdir(count int) ([]sftpd.NamedAttr, error) {
	fis, e := d.d.Readdir(count)
	if e != nil {
		return nil, e
	}
	rs := make([]sftpd.NamedAttr, len(fis))
	for i, fi := range fis {
		rs[i].Name = fi.Name()
		rs[i].FillFrom(fi)
	}
	return rs, nil
}
func (d rdir) Close() error {
	return d.d.Close()
}

// Warning:
// Use your own path mangling functionality in production code.
// This can be quite non-trivial depending on the operating system.
// The code below is not sufficient for production servers.
func rfsMangle(path string) (string, error) {
	if strings.Contains(path, "..") {
		return "<invalid>", errors.New("Invalid path")
	}
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	path = "/tmp/test-sftpd/" + path
	return path, nil
}

func (fs rfs) OpenDir(path string) (sftpd.Dir, error) {
	p, e := rfsMangle(path)
	if e != nil {
		return nil, e
	}
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	return rdir{f}, nil
}

func (fs rfs) OpenFile(path string, mode uint32, a *sftpd.Attr) (sftpd.File, error) {
	p, e := rfsMangle(path)
	if e != nil {
		return nil, e
	}
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	return rfile{f: f}, nil
}

func (fs rfs) Stat(name string, islstat bool) (*sftpd.Attr, error) {
	p, e := rfsMangle(name)
	if e != nil {
		return nil, e
	}
	var fi os.FileInfo
	if islstat {
		fi, e = os.Lstat(p)
	} else {
		fi, e = os.Stat(p)
	}
	if e != nil {
		return nil, e
	}
	var a sftpd.Attr
	e = a.FillFrom(fi)

	return &a, e
}
