package main

// Synthetic file system example

import (
	"bytes"
	"errors"
	"io"
	"log"

	"github.com/taruti/sftpd"
)

type synthetic struct {
	sftpd.EmptyFS
}

type synthDir struct {
	m map[string]*synthFile
}

type synthFile struct {
	sftpd.EmptyFile
	bs []byte
}

var synthRoot = &synthDir{map[string]*synthFile{
	"foo": &synthFile{bs: []byte("foo contents")},
	"bar": &synthFile{bs: []byte("bar contents")},
}}

func (fs synthetic) OpenDir(path string) (sftpd.Dir, error) {
	d := synthDir{m: map[string]*synthFile{}}
	for k, v := range synthRoot.m {
		d.m[k] = v
	}
	return &d, nil
}

func (fs synthetic) OpenFile(path string, mode uint32, attr *sftpd.Attr) (sftpd.File, error) {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	f, ok := synthRoot.m[path]
	if ok {
		return f, nil
	}
	return nil, errors.New("Not found!")
}

func (fs synthetic) Stat(path string, islstat bool) (*sftpd.Attr, error) {
	var a sftpd.Attr
	log.Println("STAT", path)
	if path == "" || path == "/" || path == "." {
		a.Flags = sftpd.ATTR_MODE
		a.Mode = sftpd.MODE_DIR | 0755
		return &a, nil
	}
	if path[0] == '/' {
		path = path[1:]
	}
	f, ok := synthRoot.m[path]
	if ok {
		f.fillAttr(&a)
		log.Println("STAT FILE", f, a)
		return &a, nil
	}
	return nil, errors.New("not found")
}

func (d synthDir) Readdir(count int) ([]sftpd.NamedAttr, error) {
	var rs []sftpd.NamedAttr
	for k, v := range d.m {
		na := sftpd.NamedAttr{Name: k}
		v.fillAttr(&na.Attr)
		rs = append(rs, na)
		delete(d.m, k)
	}
	if len(rs) == 0 {
		return nil, io.EOF
	}
	return rs, nil
}

func (d synthDir) Close() error {
	return nil
}

func (f *synthFile) ReadAt(bs []byte, offset int64) (int, error) {
	return bytes.NewReader(f.bs).ReadAt(bs, offset)
}

func (f *synthFile) fillAttr(attr *sftpd.Attr) {
	attr.Flags = sftpd.ATTR_SIZE | sftpd.ATTR_MODE
	attr.Size = uint64(len(f.bs))
	attr.Mode = sftpd.MODE_REGULAR | 0644
}
