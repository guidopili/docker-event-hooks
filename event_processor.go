package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"text/template"
)

func (hook Hooks) reactOnEvent(event DockerEvent) {
	if false == hook.EventList.supportsEvent(event) {
		return
	}

	cmd := exec.Command(hook.Command[0], hook.Command[1:]...)

	for _, args := range hook.Arguments {
		t, _ := template.New("argument").Parse(args)
		var tpl bytes.Buffer
		t.Execute(&tpl, event)
		cmd.Args = append(cmd.Args, tpl.String())
	}

	cmd.Run()
}

func (eventList Events) supportsEvent(event DockerEvent) bool {
	for _, singleEvent := range eventList {
		if singleEvent.Type != nil && *singleEvent.Type != event.Type {
			continue
		}
		if singleEvent.Action != nil && *singleEvent.Action != event.Action {
			continue
		}
		if singleEvent.Identifier != nil && *singleEvent.Identifier != event.ID {
			// handle name etc
			continue
		}

		return true
	}

	return 0 == len(eventList)
}

func ProcessEvent(jsonConfig JsonConfig, event DockerEvent) {
	log.Println(fmt.Sprintf("Received event with ID: %v", event.ID))

	if jsonConfig.shouldFilterEvent(event) {
		log.Println(fmt.Sprintf("Event %v filtered by options.filter", event.ID))
		return
	}

	for _, hook := range jsonConfig.Hooks {
		hook.reactOnEvent(event)
	}
}
