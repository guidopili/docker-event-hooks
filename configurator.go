package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

type Command   []string // TODO: this should be accepted as string or slice

type Hooks struct {
	EventList  Events    `json:"on"`
	Commands   []Command   `json:"commands"`
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

	if config.Verbose {
		log.Printf("Config Version: %v", y.Version)
		log.Printf("Options.ComposeFilePath: %v", y.Options.ComposeFilePath)
		log.Printf("Options.Filters:")
		for _, f := range y.Options.Filters {
			j, _ := json.Marshal(f)
			log.Printf("	%v", string(j))
		}
		log.Printf("Hooks:")
		for _, h := range y.Hooks {
			j, _ := json.Marshal(h)
			log.Printf("	%v", string(j))
		}
	}

	return y
}
