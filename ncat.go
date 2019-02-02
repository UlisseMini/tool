package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/UlisseMini/tool/pkg/ncat"
)

var ncatCmd = command{
	usage: `simple ncat implementation`,
	exec:  ncatExec,
}

func ncatExec(args ...string) {
	set := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		compressConn = set.Bool("z", false,
			"compress data sent through the connection.")

		cmd = set.String("e", "",
			"execute command over connection")
	)

	set.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage of %s: [flags] <address>\n", args[0])
		set.PrintDefaults()
	}
	set.Parse(args[1:])

	if len(set.Args()) != 1 {
		log.Printf("wanted 2; got %d", len(set.Args()))
		set.Usage()
		return
	}

	rawConn, err := net.Dial("tcp", set.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	conn := io.ReadWriteCloser(rawConn)
	defer conn.Close()

	// compress the connection if requested
	if *compressConn {
		conn = compress(rawConn)
	}

	if *cmd != "" {
		a := strings.Split(*cmd, " ")
		cmd := exec.Command(a[0], a[1:]...)
		netConn, ok := conn.(net.Conn)
		if !ok {
			log.Println("Failed to assert conn as net.Conn")
			return
		}

		if err := ncat.Exec(cmd, netConn); err != nil {
			log.Println(err)
		}
		return
	}

	// Proxy file descriptors
	done := make(chan struct{})
	go func() {
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			log.Printf("copy stdin to conn: %v", err)
		}
		done <- struct{}{}
	}()

	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			log.Printf("copy conn to stdout: %v", err)
		}
		done <- struct{}{}
	}()

	<-done
}

// Compress wraps a stream with effective data compression
func compress(rw io.ReadWriteCloser) io.ReadWriteCloser {
	return nil
}
