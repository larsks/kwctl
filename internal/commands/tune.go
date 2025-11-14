package commands

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/helpers"
	"github.com/larsks/kwctl/internal/radio"
	"github.com/larsks/kwctl/internal/types"
)

type (
	TuneCommand struct {
		flags        *flag.FlagSet
		forceVfoMode bool
		radioFlags   types.RadioFlagValues
	}
)

func init() {
	Register("tune", &TuneCommand{})
}

//nolint:errcheck
func (c *TuneCommand) Init() error {
	c.flags = flag.NewFlagSet("id", flag.ContinueOnError)
	c.flags.BoolVarP(&c.forceVfoMode, "force", "f", false, "change to vfo mode before tuning")

	// Add common radio setting flags
	types.AddRadioSettingFlags(c.flags, &c.radioFlags)

	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), helpers.Unindent(`
			Usage: kwctl tune [options]

			Tune the selected VFO.

			Options:
		`))
		c.flags.PrintDefaults()
	}
	return nil
}

func (c *TuneCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	if err := c.flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	vfo, err := r.GetVFO(ctx.Config.Vfo)
	if err != nil {
		return "", fmt.Errorf("failed to read vfo %s: %w", ctx.Config.Vfo, err)
	}
	oldVfo := vfo

	// Apply common radio settings
	types.ApplyRadioSettingFlags(c.flags, &c.radioFlags, &vfo)

	if vfo != oldVfo {
		if c.forceVfoMode {
			err := r.SetVFOMode(ctx.Config.Vfo, radio.VFO_MODE_VFO)
			if err != nil {
				return "", fmt.Errorf("failed to change to vfo mode: %w", err)
			}
		}
		err := r.SetVFO(ctx.Config.Vfo, vfo)
		if err != nil {
			return "", fmt.Errorf("failed to tune vfo %s: %w", ctx.Config.Vfo, err)
		}
		vfo, err = r.GetVFO(ctx.Config.Vfo)
		if err != nil {
			return "", fmt.Errorf("failed to read vfo %s: %w", ctx.Config.Vfo, err)
		}
	}

	return fmt.Sprintf("%s\n", vfo), nil
}
