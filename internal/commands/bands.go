package commands

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/gobot/tools"
	"github.com/larsks/kwctl/pkg/radio"
	"github.com/larsks/kwctl/pkg/radio/types"
)

type (
	BandsCommand struct {
		flags *flag.FlagSet
	}
)

func init() {
	Register("bands", &BandsCommand{})
}

func (c *BandsCommand) NeedsRadio() bool {
	return true
}

//nolint:errcheck
func (c *BandsCommand) Init() error {
	c.flags = flag.NewFlagSet("bands", flag.ContinueOnError)
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), tools.Unindent(`
			Usage: kwctl bands [dual|single]

			Get or set dual/single band mode.
		`))
	}
	c.flags.PrintDefaults()
	return nil
}

func (c *BandsCommand) Run(r *radio.Radio, ctx config.Context, args []string) error {
	if err := c.flags.Parse(args); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	if c.flags.NArg() != 0 {
		mode, err := types.ParseBandMode(c.flags.Arg(0))
		if err != nil {
			return fmt.Errorf("unknown bands mode: %s", c.flags.Arg(0))
		}
		if err := r.SetBandMode(mode); err != nil {
			return fmt.Errorf("failed to set bands mode: %w", err)
		}
	}

	mode, err := r.GetBandMode()
	if err != nil {
		return fmt.Errorf("failed to read bands mode: %w", err)
	}

	fmt.Printf("%s\n", mode)
	return nil
}
