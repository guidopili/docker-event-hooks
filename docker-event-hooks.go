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

func eventsDaemon(jsonConfig JsonConfig) {
	ctx := context.Background()
	defer ctx.Done()

	cEvents, err := eventsReader(ctx, jsonConfig)

	for {
		select {
			case e := <-err:
				handleError(e)
			default:
				ProcessEvent(jsonConfig, <-cEvents)
		}
	}
}

func main() {
	Configure()
	config := ParseConfigFile(config.ConfigFilePath)
	eventsDaemon(config)
}
