package main

import (
	"bytes"
	"github.com/docker/docker/api/types/events"
	"os/exec"
	"text/template"
)

func (hook Hooks) reactOnEvent(event events.Message) {
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

func (eventList Events) supportsEvent(event events.Message) bool {
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

func ProcessEvent(yamlConfig YamlConfig, event events.Message) {
	for _, hook := range yamlConfig.Hooks {
		hook.reactOnEvent(event)
	}
}
