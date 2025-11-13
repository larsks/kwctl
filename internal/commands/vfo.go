package commands

import (
	"fmt"
	"strings"

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

	var res string
	var err error

	if flags.NArg() == 0 {
		res, err = r.SendCommand("BC")
	} else {
		res, err = r.SendCommand("BC", flags.Arg(0), flags.Arg(0))
	}

	if err != nil {
		return "", fmt.Errorf("failed to select vfo: %w", err)
	}

	parts := strings.Split(res, ",")
	return fmt.Sprintf("CONTROL: %s, PTT: %s", parts[0], parts[1]), nil
}
