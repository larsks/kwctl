package commands

import (
	"fmt"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
)

type (
	IDCommand struct{}
)

func init() {
	Register("id", &IDCommand{})
}

func (c IDCommand) Run(r *radio.Radio, _ config.Context, args []string) (string, error) {
	flags := flag.NewFlagSet("id", flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}
	return r.SendCommand("ID")
}
