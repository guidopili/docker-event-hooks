package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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

func eventsHandler(jsonConfig JsonConfig) {
	ctx := context.Background()
	defer ctx.Done()

	cEvents, err := eventsReader(ctx)

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

	if config.Detach == true {
		startBg(exec.Command(os.Args[0], append(os.Args[1:], "-d=false")...))
	}

	eventsHandler(ParseConfigFile(config.ConfigFilePath))
}
