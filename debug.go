package sftpd

//var debug func(...interface{}) = log.Println
//var debugf func(string, ...interface{}) = log.Printf

var debug func(...interface{}) = func(...interface{}) {}
var debugf func(string, ...interface{}) = func(string, ...interface{}) {}
