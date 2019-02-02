package ncat

import (
	"io/ioutil"
	"log"
)

var logger = log.New(ioutil.Discard, "",
	log.Lshortfile)

// SetLogger sets the standard package logger.
func SetLogger(l *log.Logger) {
	logger = l
}
