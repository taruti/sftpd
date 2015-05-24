package sftpd

import "strconv"

type handles struct {
	f map[string]File
	d map[string]Dir
	c int64
}

func (h *handles) init() {
	h.f = map[string]File{}
	h.d = map[string]Dir{}
}

func (h *handles) closeAll() {
	for _, x := range h.f {
		x.Close()
	}
	for _, x := range h.d {
		x.Close()
	}
}
func (h *handles) closeHandle(k string) {
	if k == "" {
		return
	}
	if k[0] == 'f' {
		x, ok := h.f[k]
		if ok {
			x.Close()
		}
		delete(h.f, k)
	} else if k[0] == 'd' {
		x, ok := h.d[k]
		if ok {
			x.Close()
		}
		delete(h.d, k)
	}
}
func (h *handles) nfiles() int { return len(h.f) }
func (h *handles) ndir() int   { return len(h.d) }

func (h *handles) newFile(f File) string {
	h.c++
	k := "f" + strconv.FormatInt(h.c, 16)
	h.f[k] = f
	return k
}
func (h *handles) newDir(f Dir) string {
	h.c++
	k := "d" + strconv.FormatInt(h.c, 16)
	h.d[k] = f
	return k
}
func (h *handles) getFile(n string) File {
	return h.f[n]
}
func (h *handles) getDir(n string) Dir {
	return h.d[n]
}
