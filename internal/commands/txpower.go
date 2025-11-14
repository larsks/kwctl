package commands

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/helpers"
	"github.com/larsks/kwctl/internal/radio"
)

type (
	TxPowerCommand struct {
		flags *flag.FlagSet
	}
)

var txpowerNames map[string]string = map[string]string{
	"high":   "0",
	"medium": "1",
	"low":    "2",
}

var txpowerNumbers map[string]string = helpers.ReverseMap(txpowerNames)

func init() {
	Register("txpower", &TxPowerCommand{})
}

func (c *TxPowerCommand) NeedsRadio() bool {
	return true
}

func (c *TxPowerCommand) Init() error {
	c.flags = flag.NewFlagSet("txpower", flag.ContinueOnError)
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), helpers.Unindent(`
			Usage: kwctl txpower [high|medium|low]

			Set the transmit power for the selected VFO.

			Options:
		`))
		c.flags.PrintDefaults()
	}
	return nil
}

func (c *TxPowerCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	if err := c.flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	var res string
	var err error

	if c.flags.NArg() == 0 {
		// Get current txpower
		res, err = r.SendCommand("PC", ctx.Config.Vfo)
	} else {
		num, exists := txpowerNames[c.flags.Arg(0)]
		if !exists {
			return "", fmt.Errorf("unknown txpower: %s", c.flags.Arg(0))
		}
		res, err = r.SendCommand("PC", ctx.Config.Vfo, num)
	}

	if err != nil {
		return "", fmt.Errorf("txpower command failed: %w", err)
	}

	// Parse response: "vfo,txpower"
	parts := strings.Split(res, ",")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid response: %s", res)
	}

	// Map txpower value to human-readable string
	name, exists := txpowerNumbers[parts[1]]
	if !exists {
		return "", fmt.Errorf("unknown txpower: %s", res)
	}

	return name, nil
}
