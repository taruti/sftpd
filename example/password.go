package main

import (
	"crypto/rand"
	"io"
)

var testUser = "test"
var testPass = prandAlphaNumeric(16)

const alnum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Pseudo random alpha numeric password generation for this example
func prandAlphaNumeric(n int) []byte {
	bs := make([]byte, n)
	_, e := io.ReadFull(rand.Reader, bs)
	if e != nil {
		panic(e)
	}
	for i, b := range bs {
		bs[i] = alnum[int(b)%len(alnum)]
	}
	return bs
}
