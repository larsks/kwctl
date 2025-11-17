package commands

import (
	"log/slog"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
)

type (
	Command interface {
		Init() error
		NeedsRadio() bool
		Run(r *radio.Radio, ctx config.Context, args []string) error
	}
)

var commands map[string]Command = make(map[string]Command)

func Register(name string, command Command, aliases ...string) {
	if _, exists := commands[name]; exists {
		slog.Error("ignoring duplicate registration", "command", name)
		return
	}

	if err := command.Init(); err != nil {
		slog.Error("failed to initialize command", "command", name)
	}

	commands[name] = command

	for _, alias := range aliases {
		if _, exists := commands[alias]; exists {
			slog.Error("ignoring duplicate registration", "alias", alias)
			continue
		}
		commands[alias] = command
	}
}

func Lookup(name string) Command {
	if handler, ok := commands[name]; ok {
		return handler
	}

	return nil
}

func List() []string {
	var names []string
	for name := range commands {
		names = append(names, name)
	}

	return names
}
