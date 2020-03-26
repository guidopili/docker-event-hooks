package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Event struct {
	Type       *string `json:"type"`
	Action     *string `json:"action"`
	Identifier *string `json:"identifier,omitempty"`
}

type Events []Event

type Options struct {
	ComposeFilePath string `json:"compose-file-path,omitempty"`
	Filters         Events `json:"filters"`
}

type Command []string // TODO: this should be accepted as string or slice
type Arguments []string

type Hooks struct {
	EventList Events    `json:"on"`
	Command   Command   `json:"command"`
	Arguments Arguments `json:"arguments"`
}

type JsonConfig struct {
	Version string  `json:"version"`
	Options Options `json:"options"`
	Hooks   []Hooks `json:"hooks"`
}

func (jsonConfig JsonConfig) shouldFilterEvent(event DockerEvent) bool {
	return !jsonConfig.Options.Filters.supportsEvent(event)
}

func ParseConfigFile(filename string) JsonConfig {
	dat, err := ioutil.ReadFile(filename)
	handleError(err)

	var y JsonConfig

	err = json.Unmarshal(dat, &y)
	handleError(err)

	fmt.Println("Version: ", y.Version)
	fmt.Println("Options: ", y.Options.Filters)
	fmt.Println("Hooks: ", y.Hooks)

	return y
}
