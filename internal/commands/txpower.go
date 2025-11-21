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
	TxPowerCommand struct {
		flags *flag.FlagSet
	}
)

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
		//nolint:errcheck
		fmt.Fprint(c.flags.Output(), tools.Unindent(`
			Usage: kwctl txpower [high|medium|low]

			Get or set the transmit power for the selected VFO.
		`))
		c.flags.PrintDefaults()
	}
	return nil
}

func (c *TxPowerCommand) Run(r *radio.Radio, ctx config.Context, args []string) error {
	if err := c.flags.Parse(args); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	if c.flags.NArg() == 1 {
		tx, err := types.ParseTxPower(c.flags.Arg(0))
		if err != nil {
			return fmt.Errorf("failed to parse tx power: %w", err)
		}
		if err := r.SetTxPower(ctx.Config.Vfo, tx); err != nil {
			return fmt.Errorf("failed to set txpower: %w", err)
		}
	}

	txpower, err := r.GetTxPower(ctx.Config.Vfo)
	if err != nil {
		return fmt.Errorf("failed to get tx power: %w", err)
	}

	fmt.Printf("%s\n", txpower)

	return nil
}
