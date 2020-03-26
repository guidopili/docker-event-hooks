package main

import (
	"flag"
	"os"
)

var config Config

type Config struct {
	ConfigFilePath string
	Verbose        bool
}

func Configure() {
	command := flag.NewFlagSet("list", flag.ExitOnError)

	command.StringVar(&config.ConfigFilePath, "file", "docker-events-hooks.json", "File path (defaults to current dir)")
	command.StringVar(&config.ConfigFilePath, "f", "docker-events-hooks.json", "File path (defaults to current dir)")
	command.BoolVar(&config.Verbose, "verbose", false, "Verbose mode")
	command.BoolVar(&config.Verbose, "v", false, "Verbose mode")

	_ = command.Parse(os.Args[1:])
}
