package main

import (
	"errors"
	"fmt"

	"github.com/sejvlond/tarsier/plugins"

	"gopkg.in/yaml.v2"
)

var CommandNotRegistered = errors.New("Command not registered")

type Commander struct {
	Commands map[string]plugins.Command
}

func NewCommander() (*Commander, error) {
	cm := &Commander{
		Commands: make(map[string]plugins.Command),
	}
	return cm, nil
}

func (cm *Commander) Register(command plugins.Command) error {
	if _, exists := cm.Commands[command.Name()]; exists {
		return fmt.Errorf("Command '%v' was already registered", command.Name())
	}
	cm.Commands[command.Name()] = command
	return nil
}

func (cm *Commander) Execute(name string, data interface{}) (string, error) {
	command, exists := cm.Commands[name]
	if !exists {
		return "", CommandNotRegistered
	}
	if cmdWithCustomData, ok := command.(plugins.HasDataStruct); ok {
		raw, err := yaml.Marshal(data)
		if err != nil {
			return "Invalid data", fmt.Errorf("Unable to marshal data '%v'", err)
		}
		data = cmdWithCustomData.DataStruct()
		err = yaml.Unmarshal(raw, data)
		if err != nil {
			return "Invalid data", fmt.Errorf(
				"Unable to unmarshal data to custom struct '%v'", err)
		}
	}
	return command.Execute(data)
}
