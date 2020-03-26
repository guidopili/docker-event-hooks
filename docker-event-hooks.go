package main

import (
	"context"
	"fmt"
	"os"
)

func handleError(err error) {
	if err == nil {
		return
	}

	if config.Verbose {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, fmt.Sprintf("\033[1;31m%s\033[0m", err.Error()))
	os.Exit(1)
}

func eventsDaemon(yamlConfig YamlConfig) {
	ctx := context.Background()
	defer ctx.Done()

	cEvents, _ := eventsReader(ctx, yamlConfig)

	for {
		ProcessEvent(yamlConfig, <-cEvents)
	}
}

func main() {
	Configure()
	config := ParseConfigFile(config.ConfigFilePath)
	eventsDaemon(config)
}
