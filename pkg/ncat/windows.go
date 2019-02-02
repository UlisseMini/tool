// +build windows

package ncat

import (
	"io"
	"net"
	"os/exec"
)

// Exec execeutes a command over a stream (usually a net.Conn),
// It will execute cmd inside a pty if possible.
func Exec(cmd *exec.Cmd, conn net.Conn) error {
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Stdin = conn

	err := cmd.Run()
	if err != nil {
		io.WriteString(conn, err.Error())
	}

	return err
}
