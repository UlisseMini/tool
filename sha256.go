// sha.go contains a way to print the checksum of files.

package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var sha = command{
	usage: `print the sha256 sum of files.`,
	exec:  shaExec,
}

func shaExec(args ...string) {
	set := flag.NewFlagSet(args[0], flag.ExitOnError)
	set.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage of %s: <files...>", args[0])
		set.PrintDefaults()
	}
	set.Parse(args[1:])

	if len(set.Args()) < 1 {
		set.Usage()
		return
	}

	files := set.Args()
	h := sha256.New()

	for _, file := range files {
		h.Reset()
		f, err := os.Open(file)
		if err != nil {
			log.Printf("Failed to open %q: %v", file, err)
			continue
		}

		if _, err := io.Copy(h, f); err != nil {
			log.Printf("Failed to hash: %v", err)
		}

		// Output it
		sum := h.Sum(nil)
		fmt.Printf("%x  %s\n", sum, file)
	}
}
