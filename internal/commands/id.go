package commands

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/gobot/tools"
	"github.com/larsks/kwctl/pkg/radio"
)

type (
	IDCommand struct {
		flags *flag.FlagSet
	}
)

func init() {
	Register("id", &IDCommand{})
}

func (c *IDCommand) NeedsRadio() bool {
	return true
}

//nolint:errcheck
func (c *IDCommand) Init() error {
	c.flags = flag.NewFlagSet("id", flag.ContinueOnError)
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), tools.Unindent(`
			Usage: kwctl id

			Display the radio ID response.
		`))
		c.flags.PrintDefaults()
	}
	return nil
}

func (c *IDCommand) Run(r *radio.Radio, _ config.Context, args []string) error {
	if err := c.flags.Parse(args); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	res, err := r.SendCommand("ID")
	if err != nil {
		return fmt.Errorf("failed to get radio id: %w", err)
	}

	fmt.Printf("%s\n", res)

	return nil
}
