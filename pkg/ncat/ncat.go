// Package ncat provides methods for more flexible networking
package ncat

import (
	"io"
)

// Compress compresses a data stream (usually a net.Conn) and tries to
// optimize it for speed.
func Compress(rw io.ReadWriteCloser) io.ReadWriteCloser {
	return nil
}
