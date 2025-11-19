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
	TxPowerCommand struct {
		flags *flag.FlagSet
	}
)

var txpowerNames map[string]radio.TxPower = map[string]radio.TxPower{
	"high":   radio.TX_POWER_HIGH,
	"medium": radio.TX_POWER_MEDIUM,
	"low":    radio.TX_POWER_LOW,
}

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
		fmt.Fprint(c.flags.Output(), helpers.Unindent(`
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
		num, exists := txpowerNames[c.flags.Arg(0)]
		if !exists {
			return fmt.Errorf("unknown txpower: %s", c.flags.Arg(0))
		}
		err := r.SetTxPower(ctx.Config.Vfo, num)
		if err != nil {
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
