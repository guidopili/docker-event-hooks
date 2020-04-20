package main

import (
	"flag"
	"fmt"
	"os"
)

var config Config

var version string
var binName string

type Config struct {
	ConfigFilePath string
	Verbose        bool
	Detach         bool
}

func initConfig() {
	config.ConfigFilePath = "config.json"
	config.Verbose = false
	config.Detach = false
}

func configureFlags(args []string) *flag.FlagSet {
	command := flag.NewFlagSet(binName, flag.ExitOnError)

	command.StringVar(&config.ConfigFilePath, "file", config.ConfigFilePath, "File path (defaults to current dir)")
	command.StringVar(&config.ConfigFilePath, "f", config.ConfigFilePath, "File path (defaults to current dir)")
	command.BoolVar(&config.Verbose, "verbose", config.Verbose, "Verbose mode")
	command.BoolVar(&config.Verbose, "v", config.Verbose, "Verbose mode")
	command.BoolVar(&config.Detach, "detach", config.Detach, "Detach from shell mode")
	command.BoolVar(&config.Detach, "d", config.Detach, "Detach from shell mode")

	command.Parse(args)

	return command
}

func printHelp() {
	_, _ = fmt.Fprintf(os.Stderr, "Command required at first position. Options are \n - start\n - stop\n - version\n - help\n")
	os.Exit(1)
}

func Configure() {
	initConfig()

	requestedCommand := ""
	if len(os.Args) > 1 {
		command := configureFlags(os.Args[1:])
		requestedCommand = command.Arg(0)
		if command.NArg() > 1 {
			_ = configureFlags(command.Args()[1:])
		}
	}

	switch requestedCommand {
	case "start":
	case "stop":
		stopBg()
	case "version":
		fmt.Println(version)
		os.Exit(0)
	case "help":
		printHelp()
	default:
		printHelp()
	}
}
