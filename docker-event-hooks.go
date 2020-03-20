package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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

func eventsDaemon() {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	opts := types.EventsOptions{}
	cEvents, _ := cli.Events(ctx, opts)

	defer ctx.Done()

	for {
		ProcessEvent(<-cEvents)
	}
}

func main() {
	Configure()
	ParseConfigFile(config.ConfigFilePath)
	eventsDaemon()
}
