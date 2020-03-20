package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Event struct {
	Type       string `yaml:"type"`
	Event      string `yaml:"event"`
	Identifier string `yaml:"identifier,omitempty"`
}

type Options struct {
	ComposeFilePath string  `yaml:"compose-file-path,omitempty"`
	Filters         []Event `yaml:"filters"`
}

type Command string // TODO: this should be accepted as string or slice

type Hooks struct {
	On      []Event `yaml:"on"`
	Command Command `yaml:"command"`
}

type YamlConfig struct {
	Version string  `yaml:"version"`
	Options Options `yaml:"options"`
	Hooks   []Hooks `yaml:"hooks"`
}

func ParseConfigFile(filename string) YamlConfig {
	dat, err := ioutil.ReadFile(filename)
	handleError(err)

	var y YamlConfig

	err = yaml.Unmarshal(dat, &y)
	handleError(err)

	fmt.Println("Version: ", y.Version)
	fmt.Println("Options: ", y.Options.Filters)
	fmt.Println("Hooks: ", y.Hooks)

	return y
}
