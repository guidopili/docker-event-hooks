package main

import (
	"bytes"
	"github.com/docker/docker/api/types/events"
	"os/exec"
	"text/template"
)

func (hook Hooks) reactOnEvent(event events.Message) {
	if false == hook.supportsEvent(event) {
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

func (hook Hooks) supportsEvent(event events.Message) bool {
	for _, on := range hook.On {
		if on.Type != nil && *on.Type != event.Type {
			continue
		}
		if on.Action != nil && *on.Action != event.Action {
			continue
		}
		if on.Identifier != nil && *on.Identifier != event.ID {
			// handle name etc
			continue
		}

		return true
	}

	return 0 == len(hook.On)
}

func ProcessEvent(yamlConfig YamlConfig, event events.Message) {
	for _, hook := range yamlConfig.Hooks {
		hook.reactOnEvent(event)
	}
}
