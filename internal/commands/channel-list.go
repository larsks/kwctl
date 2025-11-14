package commands

import (
	"errors"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/helpers"
	"github.com/larsks/kwctl/internal/radio"
	"github.com/larsks/kwctl/internal/types"
)

type (
	ChannelListCommand struct {
		flags *flag.FlagSet
	}
)

func init() {
	Register("channel-list", &ChannelListCommand{}, "list")
}

//nolint:errcheck
func (c *ChannelListCommand) Init() error {
	c.flags = flag.NewFlagSet("channel-list", flag.ContinueOnError)
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), helpers.Unindent(`
			Usage: kwctl channel-list [options] <range> [<range> [...]]

			List a range of channels.

			Arguments:
				range      A range specification (e.g. "1", "1-10", "1,5,10,15,20")

			Options:
			`))
		c.flags.PrintDefaults()
	}

	return nil
}

func (c *ChannelListCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	if err := c.flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	var ranges []string
	if c.flags.NArg() == 0 {
		ranges = []string{"0-999"}
	} else {
		ranges = c.flags.Args()
	}

	for _, arg := range ranges {
		for channelNumber, err := range helpers.RangeIterator(arg) {
			if err != nil {
				return "", fmt.Errorf("invalid range: %w", err)
			}
			if channelNumber < 0 || channelNumber > 999 {
				return "", fmt.Errorf("invalid range (channels must be between 0 and 999)")
			}

			channel, err := r.GetMemoryChannel(channelNumber)
			if err != nil {
				if errors.Is(err, radio.ErrUnavailableCommand) {
					channel = types.EmptyChannel
				} else {
					return "", fmt.Errorf("failed to list channels: %w", err)
				}
			}

			if channel.RxFreq == 0 {
				fmt.Printf("[      ] %03d\n", channelNumber)
			} else {
				fmt.Printf("%s\n", channel)
			}
		}
	}
	return "", nil
}
