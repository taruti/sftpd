package fuzz

import (
	"bytes"
	"io"

	"github.com/taruti/sftpd"
)

// Fuzz is the interface for the go-fuzz.
func Fuzz(data []byte) int {
	frd := &fakeRandChannel{bytes.NewReader(data), 0}
	err := sftpd.ServeChannel(frd, sftpd.EmptyFS{})
	if err != nil {
		return 0
	}
	if frd.nw >= 10 {
		return 2
	}
	return 1
}

type fakeRandChannel struct {
	*bytes.Reader
	nw int
}

func (f *fakeRandChannel) Write(bs []byte) (int, error) {
	f.nw += len(bs)
	return len(bs), nil
}
func (*fakeRandChannel) Close() error      { return nil }
func (*fakeRandChannel) CloseWrite() error { return nil }
func (*fakeRandChannel) SendRequest(name string, wantReply bool, payload []byte) (bool, error) {
	return true, nil
}
func (*fakeRandChannel) Stderr() io.ReadWriter { return nil }
