// +build !windows

package ncat

import (
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

var DefaultWinSize = &pty.Winsize{
	Rows: 28,
	Cols: 94,
}

// Exec execeutes a command over a stream (usually a net.Conn),
// It will execute cmd inside a pty if on a unix based system.
func Exec(cmd *exec.Cmd, conn net.Conn) error {
	// Start the command with a pty.
	ptmx, err := pty.StartWithSize(cmd, DefaultWinSize)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer ptmx.Close()

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)

	// Copy pty file descriptors until one finishes
	done := make(chan struct{})
	go func() {
		io.Copy(ptmx, conn)
		done <- struct{}{}
	}()

	go func() {
		io.Copy(conn, ptmx)
		done <- struct{}{}
	}()

	<-done
	return nil
}
