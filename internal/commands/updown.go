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
	UpCommand struct {
		flags *flag.FlagSet
	}
	DownCommand struct {
		flags *flag.FlagSet
	}
)

func init() {
	Register("up", &UpCommand{})
	Register("down", &DownCommand{})
}

func (c *UpCommand) NeedsRadio() bool {
	return true
}

//nolint:errcheck
func (c *UpCommand) Init() error {
	c.flags = flag.NewFlagSet("up", flag.ContinueOnError)
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), helpers.Unindent(`
			Usage: kwctl up

			Emulate the microphone Up key
		`))
		c.flags.PrintDefaults()
	}
	return nil
}

func (c *UpCommand) Run(r *radio.Radio, _ config.Context, args []string) error {
	if err := c.flags.Parse(args); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	if err := r.MicUp(); err != nil {
		return fmt.Errorf("failed to up: %w", err)
	}

	return nil
}

func (c *DownCommand) NeedsRadio() bool {
	return true
}

//nolint:errcheck
func (c *DownCommand) Init() error {
	c.flags = flag.NewFlagSet("down", flag.ContinueOnError)
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), helpers.Unindent(`
			Usage: kwctl down

			Emulate the microphone Down key
		`))
		c.flags.PrintDefaults()
	}
	return nil
}

func (c *DownCommand) Run(r *radio.Radio, _ config.Context, args []string) error {
	if err := c.flags.Parse(args); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	if err := r.MicDown(); err != nil {
		return fmt.Errorf("failed to down: %w", err)
	}

	return nil
}
