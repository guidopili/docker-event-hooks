package main

import (
	"flag"
	"fmt"
	"os"
)

var config Config
var version string

type Config struct {
	ConfigFilePath string
	Verbose        bool
	Detach         bool
}

func configureStartArgs() {
	command := flag.NewFlagSet("start", flag.ExitOnError)

	command.StringVar(&config.ConfigFilePath, "file", "docker-events-hooks.json", "File path (defaults to current dir)")
	command.StringVar(&config.ConfigFilePath, "f", "docker-events-hooks.json", "File path (defaults to current dir)")
	command.BoolVar(&config.Verbose, "verbose", false, "Verbose mode")
	command.BoolVar(&config.Verbose, "v", false, "Verbose mode")
	command.BoolVar(&config.Detach, "detach", false, "Detach from shell mode")
	command.BoolVar(&config.Detach, "d", false, "Detach from shell mode")

	_ = command.Parse(os.Args[2:])
}

func printHelp() {
	_, _ = fmt.Fprintf(os.Stderr, "Command required at first position. Options are \n - configureStartArgs\n - stop\n - version\n - help\n")
	os.Exit(1)
}

func Configure() {
	if len(os.Args) < 2 {
		printHelp()
	}

	switch os.Args[1] {
	case "start":
		configureStartArgs()
		break
	case "stop":
		stopBg()
		break
	case "version":
		fmt.Println(version)
		os.Exit(0)
	case "help":
	default:
		printHelp()
	}
}
