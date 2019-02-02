// Package main provides useful some devops and cli tools.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
)

type command struct {
	usage string               // short string explaining the usage
	exec  func(args ...string) // execute the command
}

// commands will be stored in here, the command name must not contain spaces.
// TODO: Automatically find the correct name via the function name.
var commands = map[string]command{
	"sha256": sha,
	"ncat":   ncatCmd,
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage of %s: [flags] <command> [args...]\n", os.Args[0])
		flag.PrintDefaults()
	}

	list := flag.Bool("l", false, "list all commands")
	flag.Parse()
	if *list {
		listCommands()
		return
	}

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	cmd := commands[os.Args[1]]
	if reflect.DeepEqual(cmd, nil) {
		fmt.Printf("Unknown command %q\n", os.Args[1])
		return
	}

	// Run the command
	cmd.exec(os.Args[1:]...)
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("[*] ")
}

// listCommands lists all commands.
func listCommands() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"List of all commands:\n")

	for name, cmd := range commands {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s - %s\n", name, cmd.usage)
	}
}
