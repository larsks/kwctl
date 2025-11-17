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
	BandsCommand struct {
		flags *flag.FlagSet
	}
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

func (c *BandsCommand) NeedsRadio() bool {
	return true
}

//nolint:errcheck
func (c *BandsCommand) Init() error {
	c.flags = flag.NewFlagSet("bands", flag.ContinueOnError)
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), helpers.Unindent(`
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

	var res string
	var err error

	if c.flags.NArg() == 0 {
		// Get current power
		res, err = r.SendCommand("DL")
	} else {
		num, exists := bandsNames[c.flags.Arg(0)]
		if !exists {
			return fmt.Errorf("unknown bands mode: %s", c.flags.Arg(0))
		}
		res, err = r.SendCommand("DL", num)
	}

	if err != nil {
		return fmt.Errorf("bands command failed: %w", err)
	}

	parts := strings.Split(res, ",")
	if len(parts) < 1 {
		return fmt.Errorf("invalid response: %s", res)
	}

	name, exists := bandsNumbers[parts[0]]
	if !exists {
		return fmt.Errorf("unknown bands mode: %s", res)
	}

	fmt.Printf("%s\n", name)
	return nil
}
