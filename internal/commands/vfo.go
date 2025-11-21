package commands

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/gobot/tools"
	"github.com/larsks/kwctl/pkg/radio"
)

type (
	VFOCommand struct {
		flags *flag.FlagSet
	}
)

func init() {
	Register("vfo", &VFOCommand{})
}

func (c *VFOCommand) NeedsRadio() bool {
	return true
}

//nolint:errcheck
func (c *VFOCommand) Init() error {
	c.flags = flag.NewFlagSet("vfo", flag.ContinueOnError)
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), tools.Unindent(`
			Usage: kwctl vfo [0|1]

			Get or set ptt/control VFO.
		`))
		c.flags.PrintDefaults()
	}
	return nil
}

func (c *VFOCommand) Run(r *radio.Radio, ctx config.Context, args []string) error {
	if err := c.flags.Parse(args); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	var res string
	var err error

	if c.flags.NArg() == 0 {
		res, err = r.SendCommand("BC")
	} else {
		res, err = r.SendCommand("BC", c.flags.Arg(0), c.flags.Arg(0))
	}

	if err != nil {
		return fmt.Errorf("failed to select vfo: %w", err)
	}

	parts := strings.Split(res, ",")
	fmt.Printf("CONTROL: %s, PTT: %s\n", parts[0], parts[1])

	return nil
}
