package sftpd

import (
	"errors"
)

var Failure = errors.New("Failure")

type EmptyFile struct{}

func (EmptyFile) Close() error                       { return nil }
func (EmptyFile) ReadAt([]byte, int64) (int, error)  { return 0, Failure }
func (EmptyFile) WriteAt([]byte, int64) (int, error) { return 0, Failure }
func (EmptyFile) FStat() (*Attr, error)              { return nil, Failure }
func (EmptyFile) FSetStat(*Attr) error               { return Failure }

type EmptyFS struct{}

func (EmptyFS) OpenFile(string, uint32, *Attr) (File, error)  { return nil, Failure }
func (EmptyFS) OpenDir(string) (Dir, error)                   { return nil, Failure }
func (EmptyFS) Remove(string) error                           { return Failure }
func (EmptyFS) Rename(string, string, uint32) error           { return Failure }
func (EmptyFS) Mkdir(string, *Attr) error                     { return Failure }
func (EmptyFS) Rmdir(string) error                            { return Failure }
func (EmptyFS) Stat(string, bool) (*Attr, error)              { return nil, Failure }
func (EmptyFS) SetStat(string, *Attr) error                   { return Failure }
func (EmptyFS) ReadLink(p string) (string, error)             { return "", Failure }
func (EmptyFS) CreateLink(p string, t string, f uint32) error { return Failure }
func (EmptyFS) RealPath(p string) (string, error)             { return simpleRealPath(p), nil }

func simpleRealPath(path string) string {
	switch path {
	case "", ".":
		path = "/"
	}
	return path
}
