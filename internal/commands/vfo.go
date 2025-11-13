package commands

import (
	"fmt"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
)

type (
	VFOCommand struct{}
)

func init() {
	Register("vfo", &VFOCommand{})
}

func (c VFOCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	flags := flag.NewFlagSet("vfo", flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	var response string
	var err error

	if flags.NArg() == 0 {
		response, err = r.SendCommand("BC")
	} else {
		response, err = r.SendCommand("BC", flags.Arg(0), flags.Arg(0))
	}

	if err != nil {
		return "", fmt.Errorf("failed to select vfo: %w", err)
	}

	return response, nil
}
