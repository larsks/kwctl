package commands

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/helpers"
	"github.com/larsks/kwctl/pkg/radio"
)

type (
	ModeCommand struct {
		flags *flag.FlagSet
	}
)

func init() {
	Register("mode", &ModeCommand{})
}

var modeNames map[string]radio.VfoMode = map[string]radio.VfoMode{
	"vfo":    radio.VFO_MODE_VFO,
	"memory": radio.VFO_MODE_MEMORY,
	"call":   radio.VFO_MODE_CALL,
	"wx":     radio.VFO_MODE_WX,
}

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

func (c *ModeCommand) Run(r *radio.Radio, ctx config.Context, args []string) error {
	if err := c.flags.Parse(args); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	var err error
	var mode radio.VfoMode

	if c.flags.NArg() == 1 {
		var exists bool
		mode, exists = modeNames[c.flags.Arg(0)]
		if !exists {
			return fmt.Errorf("unknown mode: %s", c.flags.Arg(0))
		}

		err = r.SetVFOMode(ctx.Config.Vfo, mode)
		if err != nil {
			return fmt.Errorf("failed to set mode for vfo %s: %w", ctx.Config.Vfo, err)
		}
	}

	mode, err = r.GetVFOMode(ctx.Config.Vfo)
	if err != nil {
		return fmt.Errorf("failed to get mode for vfo %s: %w", ctx.Config.Vfo, err)
	}

	fmt.Printf("%s\n", mode)

	return nil
}
