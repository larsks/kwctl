package commands

import (
	"fmt"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
)

type (
	BandsCommand struct{}
)

var bandsNames map[string]string = map[string]string{
	"dual":   "0",
	"single": "1",
}

var bandsNumbers map[string]string = map[string]string{
	"0": "dual",
	"1": "single",
}

func init() {
	Register("bands", &BandsCommand{})
}

func (c BandsCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	flags := flag.NewFlagSet("bands", flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	var res string
	var err error

	if flags.NArg() == 0 {
		// Get current power
		res, err = r.SendCommand("DL")
	} else {
		num, exists := bandsNames[flags.Arg(0)]
		if !exists {
			return "", fmt.Errorf("unknown bands mode: %s", flags.Arg(0))
		}
		res, err = r.SendCommand("DL", num)
	}

	if err != nil {
		return "", fmt.Errorf("bands command failed: %w", err)
	}

	parts := strings.Split(res, ",")
	if len(parts) < 1 {
		return "", fmt.Errorf("invalid response: %s", res)
	}

	name, exists := bandsNumbers[parts[0]]
	if !exists {
		return "", fmt.Errorf("unknown bands mode: %s", res)
	}

	return name, nil
}
