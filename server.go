package sftpd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"

	"github.com/taruti/binp"
	"github.com/taruti/bytepool"
	"golang.org/x/crypto/ssh"
)

var sftpSubSystem = []byte{0, 0, 0, 4, 115, 102, 116, 112}

// IsSftpRequest checks whether a given ssh.Request is for sftp.
func IsSftpRequest(req *ssh.Request) bool {
	return req.Type == "subsystem" && bytes.Equal(sftpSubSystem, req.Payload)
}

var initReply = []byte{0, 0, 0, 5, ssh_FXP_VERSION, 0, 0, 0, 3}

// ServeChannel serves a ssh.Channel with the given FileSystem.
func ServeChannel(c ssh.Channel, fs FileSystem) error {
	defer c.Close()
	var h handles
	h.init()
	defer h.closeAll()
	brd := bufio.NewReader(c)
	var e error
	var plen int
	var op byte
	var bs []byte
	var id uint32
	for {
		if e != nil {
			debug("Sending errror", e)
			e = writeErr(c, id, e)
			if e != nil {
				return e
			}
		}
		discard(brd, plen)
		plen, op, e = readPacketHeader(brd)
		if e != nil {
			return e
		}
		plen--
		debugf("CR op=%v data len=%d\n", ssh_fxp(op), plen)
		if plen < 2 {
			return errors.New("Packet too short")
		}
		// Feeding too large values to peek is ok, it just errors.
		bs, e = brd.Peek(plen)
		if e != nil {
			return e
		}
		debugf("Data %X\n", bs)
		p := binp.NewParser(bs)
		switch op {
		case ssh_FXP_INIT:
			e = wrc(c, initReply)
		case ssh_FXP_OPEN:
			var path string
			var flags uint32
			var a Attr
			e = parseAttr(p.B32(&id).B32String(&path).B32(&flags), &a).End()
			if e != nil {
				return e
			}
			if h.nfiles() >= maxFiles {
				e = etoomany
				continue
			}
			var f File
			f, e = fs.OpenFile(path, flags, &a)
			if e != nil {
				continue
			}
			e = writeHandle(c, id, h.newFile(f))
		case ssh_FXP_CLOSE:
			var handle string
			e = p.B32(&id).B32String(&handle).End()
			if e != nil {
				return e
			}
			h.closeHandle(handle)
			e = writeErr(c, id, nil)
		case ssh_FXP_READ:
			var handle string
			var offset uint64
			var length uint32
			var n int
			e = p.B32(&id).B32String(&handle).B64(&offset).B32(&length).End()
			if e != nil {
				return e
			}
			f := h.getFile(handle)
			if f == nil {
				return einvhandle
			}
			if length > 64*1024 {
				length = 64 * 1024
			}
			bs := bytepool.Alloc(int(length))
			n, e = f.ReadAt(bs, int64(offset))
			// Handle go readers that return io.EOF and bytes at the same time.
			if e == io.EOF && n > 0 {
				e = nil
			}
			if e != nil {
				bytepool.Free(bs)
				continue
			}
			bs = bs[0:n]
			e = wrc(c, binp.Out().B32(1+4+4+uint32(len(bs))).Byte(ssh_FXP_DATA).B32(id).B32(uint32(len(bs))).Out())
			if e == nil {
				e = wrc(c, bs)
			}
			bytepool.Free(bs)
		case ssh_FXP_WRITE:
			var handle string
			var offset uint64
			var length uint32
			p.B32(&id).B32String(&handle).B64(&offset).B32(&length)
			f := h.getFile(handle)
			if f == nil {
				return einvhandle
			}
			var bs []byte
			e = p.NBytesPeek(int(length), &bs).End()
			if e != nil {
				return e
			}
			_, e = f.WriteAt(bs, int64(offset))
			e = writeErr(c, id, e)
		case ssh_FXP_LSTAT, ssh_FXP_STAT:
			var path string
			var a *Attr
			e = p.B32(&id).B32String(&path).End()
			if e != nil {
				return e
			}
			a, e = fs.Stat(path, op == ssh_FXP_LSTAT)
			debug("stat/lstat", path, "=>", a, e)
			e = writeAttr(c, id, a, e)
		case ssh_FXP_FSTAT:
			var handle string
			var a *Attr
			e = p.B32(&id).B32String(&handle).End()
			if e != nil {
				return e
			}
			f := h.getFile(handle)
			if f == nil {
				return einvhandle
			}
			a, e = f.FStat()
			e = writeAttr(c, id, a, e)
		case ssh_FXP_SETSTAT:
			var path string
			var a Attr
			e = parseAttr(p.B32(&id).B32String(&path), &a).End()
			if e != nil {
				return e
			}
			e = writeErr(c, id, fs.SetStat(path, &a))
		case ssh_FXP_FSETSTAT:
			var handle string
			var a Attr
			e = parseAttr(p.B32(&id).B32String(&handle), &a).End()
			if e != nil {
				return e
			}
			f := h.getFile(handle)
			if f == nil {
				return einvhandle
			}
			e = writeErr(c, id, f.FSetStat(&a))
		case ssh_FXP_OPENDIR:
			var path string
			var dh Dir
			e = p.B32(&id).B32String(&path).End()
			if e != nil {
				return e
			}
			dh, e = fs.OpenDir(path)
			debug("opendir", id, path, "=>", dh, e)
			if e != nil {
				continue
			}
			e = writeHandle(c, id, h.newDir(dh))

		case ssh_FXP_READDIR:
			var handle string
			e = p.B32(&id).B32String(&handle).End()
			if e != nil {
				return e
			}
			f := h.getDir(handle)
			if f == nil {
				return einvhandle
			}
			var fis []NamedAttr
			fis, e = f.Readdir(1024)
			debug("readdir", id, handle, fis, e)
			if e != nil {
				continue
			}
			var l binp.Len
			o := binp.Out().LenB32(&l).LenStart(&l).Byte(ssh_FXP_NAME).B32(id).B32(uint32(len(fis)))
			for _, fi := range fis {
				// FIXME should we do special handling for long names?
				n := fi.Name
				o.B32String(n).B32String(n).B32(fi.Flags)
				if fi.Flags&ATTR_SIZE != 0 {
					o.B64(uint64(fi.Size))
				}
				if fi.Flags&ATTR_UIDGID != 0 {
					o.B32(fi.Uid).B32(fi.Gid)
				}
				if fi.Flags&ATTR_MODE != 0 {
					o.B32(fi.Mode)
				}
				if fi.Flags&ATTR_TIME != 0 {
					o.B32(fi.ATime).B32(fi.MTime)
				}
			}
			o.LenDone(&l)
			e = wrc(c, o.Out())

		case ssh_FXP_REMOVE:
			var path string
			e = p.B32(&id).B32String(&path).End()
			if e != nil {
				return e
			}
			e = writeErr(c, id, fs.Remove(path))
		case ssh_FXP_MKDIR:
			var path string
			var a Attr
			p = p.B32(&id).B32String(&path)
			e = parseAttr(p, &a).End()
			if e != nil {
				return e
			}
			e = writeErr(c, id, fs.Mkdir(path, &a))
		case ssh_FXP_RMDIR:
			var path string
			e = p.B32(&id).B32String(&path).End()
			if e != nil {
				return e
			}
			e = writeErr(c, id, fs.Rmdir(path))
		case ssh_FXP_REALPATH:
			var path, newpath string
			p.B32(&id).B32String(&path).End()
			newpath, e = fs.RealPath(path)
			debug("realpath: mapping", path, "=>", newpath, e)
			e = writeNameOnly(c, id, newpath, e)
		case ssh_FXP_RENAME:
			debug("FIXME RENAME NOT SUPPORTED")
			e = writeFail(c, id) // FIXME
		case ssh_FXP_READLINK:
			var path string
			e = p.B32(&id).B32String(&path).End()
			path, e = fs.ReadLink(path)
			e = writeNameOnly(c, id, path, e)
		case ssh_FXP_SYMLINK:
			debug("FIXME SYMLINK NOT SUPPORTED")
			e = writeFail(c, id) // FIXME
		}
		if e != nil {
			return e
		}
	}
}

var einvhandle = errors.New("Client supplied an invalid handle")
var etoomany = errors.New("Too many files")

const maxFiles = 0x100

func readPacketHeader(rd *bufio.Reader) (int, byte, error) {
	bs := make([]byte, 5)
	_, e := io.ReadFull(rd, bs)
	if e != nil {
		return 0, 0, e
	}
	return int(binary.BigEndian.Uint32(bs)), bs[4], nil
}

func parseAttr(p *binp.Parser, a *Attr) *binp.Parser {
	p = p.B32(&a.Flags)
	if a.Flags&ssh_FILEXFER_ATTR_SIZE != 0 {
		p = p.B64(&a.Size)
	}
	if a.Flags&ssh_FILEXFER_ATTR_UIDGID != 0 {
		p = p.B32(&a.Uid).B32(&a.Gid)
	}
	if a.Flags&ssh_FILEXFER_ATTR_PERMISSIONS != 0 {
		p = p.B32(&a.Mode)
	}
	if a.Flags&ssh_FILEXFER_ATTR_ACMODTIME != 0 {
		p = p.B32(&a.ATime).B32(&a.MTime)
	}
	if a.Flags&ssh_FILEXFER_ATTR_EXTENDED != 0 {
		var count uint32
		p = p.B32(&count)
		if count > 0xFF {
			return nil
		}
		ss := make([]string, 2*int(count))
		for i := 0; i < int(count); i++ {
			var k, v string
			p = p.B32String(&k).B32String(&v)
			ss[2*i+0] = k
			ss[2*i+1] = v
		}
		a.Extended = ss
	}
	return p
}

func writeAttr(c ssh.Channel, id uint32, a *Attr, e error) error {
	if e != nil {
		return writeErr(c, id, e)
	}
	var l binp.Len
	o := binp.Out().LenB32(&l).LenStart(&l).Byte(ssh_FXP_ATTRS).B32(id).B32(a.Flags)
	if a.Flags&ssh_FILEXFER_ATTR_SIZE != 0 {
		o = o.B64(a.Size)
	}
	if a.Flags&ssh_FILEXFER_ATTR_UIDGID != 0 {
		o = o.B32(a.Uid).B32(a.Gid)
	}
	if a.Flags&ssh_FILEXFER_ATTR_PERMISSIONS != 0 {
		o = o.B32(a.Mode)
	}
	if a.Flags&ssh_FILEXFER_ATTR_ACMODTIME != 0 {
		o = o.B32(a.ATime).B32(a.MTime)
	}
	if a.Flags&ssh_FILEXFER_ATTR_EXTENDED != 0 {
		count := uint32(len(a.Extended) / 2)
		o = o.B32(count)
		for _, s := range a.Extended {
			o = o.B32String(s)
		}
	}
	o.LenDone(&l)
	return wrc(c, o.Out())
}

func writeNameOnly(c ssh.Channel, id uint32, path string, e error) error {
	if e != nil {
		return writeErr(c, id, e)
	}
	var l binp.Len
	o := binp.Out().LenB32(&l).LenStart(&l).Byte(ssh_FXP_NAME).B32(id).B32(1)
	o.B32String(path).B32String(path).B32(0)
	o.LenDone(&l)
	return wrc(c, o.Out())
}

var failTmpl = []byte{0, 0, 0, 1 + 4 + 4 + 4 + 4, ssh_FXP_STATUS, 0, 0, 0, 0, 0, 0, 0, ssh_FX_FAILURE, 0, 0, 0, 0, 0, 0, 0, 0}

func writeFail(c ssh.Channel, id uint32) error {
	bs := make([]byte, len(failTmpl))
	copy(bs, failTmpl)
	binary.BigEndian.PutUint32(bs[5:], id)
	return wrc(c, bs)
}

func writeErr(c ssh.Channel, id uint32, err error) error {
	bs := make([]byte, len(failTmpl))
	copy(bs, failTmpl)
	binary.BigEndian.PutUint32(bs[5:], id)
	var code ssh_fx = ssh_FX_FAILURE
	switch err {
	case nil:
		code = ssh_FX_OK
	case io.EOF:
		code = ssh_FX_EOF
	}
	debug("Sending sftp error code", code)
	bs[12] = byte(code)
	return wrc(c, bs)
}

func writeHandle(c ssh.Channel, id uint32, handle string) error {
	return wrc(c, binp.OutCap(4+9+len(handle)).B32(uint32(9+len(handle))).B8(ssh_FXP_HANDLE).B32(id).B32String(handle).Out())
}

func wrc(c ssh.Channel, bs []byte) error {
	_, e := c.Write(bs)
	return e
}

func discard(brd *bufio.Reader, n int) error {
	if n == 0 {
		return nil
	}
	m, e := io.Copy(ioutil.Discard, &io.LimitedReader{R: brd, N: int64(n)})
	if int(m) == n && e == io.EOF {
		e = nil
	}
	return e
}
