package commands

import (
	"fmt"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
)

type (
	PowerCommand struct{}
)

var powerNames map[string]string = map[string]string{
	"high":   "0",
	"medium": "1",
	"med":    "1",
	"low":    "2",
}

var powerNumbers map[string]string = map[string]string{
	"0": "high",
	"1": "medium",
	"2": "low",
}

func init() {
	Register("power", &PowerCommand{})
}

func (c PowerCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	flags := flag.NewFlagSet("power", flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	var res string
	var err error

	if flags.NArg() == 0 {
		// Get current power
		res, err = r.SendCommand("PC", ctx.Config.Vfo)
	} else {
		num, exists := powerNames[flags.Arg(0)]
		if !exists {
			return "", fmt.Errorf("unknown power: %s", flags.Arg(0))
		}
		res, err = r.SendCommand("PC", ctx.Config.Vfo, num)
	}

	if err != nil {
		return "", fmt.Errorf("power command failed: %w", err)
	}

	// Parse response: "vfo,power"
	parts := strings.Split(res, ",")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid response: %s", res)
	}

	// Map power value to human-readable string
	name, exists := powerNumbers[parts[1]]
	if !exists {
		return "", fmt.Errorf("unknown power: %s", res)
	}

	return name, nil
}
