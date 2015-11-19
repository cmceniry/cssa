package main

import (
	"github.com/cmceniry/cssa/command"
	"flag"
	"fmt"
	"os"
)

func main() {

	generalFlag := flag.NewFlagSet("general", flag.ExitOnError)
	configfile := generalFlag.String("conf", "/etc/cassandra/cssa.conf", "cssa configuration file")
	generalFlag.Parse(os.Args[1:])

	opts, err := command.ParseConfigFile(*configfile)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(-2)
	}

	if len(generalFlag.Args()) < 1 {
		fmt.Printf("No subcommand given\n")
		os.Exit(-2)
	}

	if cmd, ok := command.Commands[generalFlag.Arg(0)]; ok {
		cmd(opts, generalFlag.Args()[1:])
	} else {
		fmt.Printf("Unrecognized subcommand\n")
		os.Exit(-2)
	}

}
