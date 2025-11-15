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
	ModeCommand struct {
		flags *flag.FlagSet
	}
)

func init() {
	Register("mode", &ModeCommand{})
}

var modeNames map[string]string = map[string]string{
	"vfo":    "0",
	"memory": "1",
	"call":   "2",
	"wx":     "3",
}

var modeNumbers map[string]string = helpers.ReverseMap(modeNames)

func (c *ModeCommand) NeedsRadio() bool {
	return true
}

func (c *ModeCommand) Init() error {
	c.flags = flag.NewFlagSet("mode", flag.ContinueOnError)
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		//nolint:errcheck
		fmt.Fprint(c.flags.Output(), helpers.Unindent(`
			Usage: kwctl mode [vfo|memory|call|wx]

			Get or set the operating mode for the selected VFO.
		`))
		c.flags.PrintDefaults()
	}
	return nil
}

func (c *ModeCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	if err := c.flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	var res string
	var err error

	if c.flags.NArg() == 0 {
		res, err = r.SendCommand("VM", ctx.Config.Vfo)
	} else {
		num, exists := modeNames[c.flags.Arg(0)]
		if !exists {
			return "", fmt.Errorf("unknown mode: %s", c.flags.Arg(0))
		}

		res, err = r.SendCommand("VM", ctx.Config.Vfo, num)
	}

	if err != nil {
		return "", fmt.Errorf("mode command failed: %w", err)
	}

	parts := strings.Split(res, ",")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid response: %s", res)
	}

	name, exists := modeNumbers[parts[1]]
	if !exists {
		return "", fmt.Errorf("unknown mode: %s", res)
	}

	return name, nil
}
