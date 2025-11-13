package commands

import (
	"log/slog"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
)

type (
	Command interface {
		Run(r *radio.Radio, ctx config.Context, args []string) (string, error)
	}
)

var commands map[string]Command = make(map[string]Command)

func Register(name string, command Command) {
	if _, exists := commands[name]; exists {
		slog.Error("ignoring duplicate registration", "command", name)
	}

	commands[name] = command
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
