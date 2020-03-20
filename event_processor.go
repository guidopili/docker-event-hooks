package main

import (
	"fmt"
	"github.com/docker/docker/api/types/events"
)

func ProcessEvent(event events.Message) {
	fmt.Print(event.TimeNano, " ", event.Actor, " ", event.Status, "\n")
}
