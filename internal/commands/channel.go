package commands

import (
	"fmt"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/helpers"
	"github.com/larsks/kwctl/internal/radio"
	"github.com/larsks/kwctl/internal/types"
)

type (
	ChannelCommand struct {
		flags *flag.FlagSet
	}
)

func init() {
	Register("channel", &ChannelCommand{})
}

//nolint:errcheck
func (c *ChannelCommand) Init() error {
	c.flags = flag.NewFlagSet("channel", flag.ContinueOnError)
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), helpers.Unindent(`
			Usage: kwctl channel [options] [<channel>|up|down]

			Get or set the current channel.

			Arguments:
				channel    Channel number (0-999) or 'up'/'down' to increment/decrement

			Options:
		`))
		c.flags.PrintDefaults()
	}
	return nil
}

func (c ChannelCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	if err := c.flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	var err error
	var channelNum int
	var channel types.Channel

	if c.flags.NArg() == 1 {
		selected := c.flags.Arg(0)

		if selected == "up" || selected == "down" {
			channelNum, err = r.GetCurrentChannelNumber(ctx.Config.Vfo)
			if err != nil {
				return "", fmt.Errorf("failed to get channel: %w", err)
			}

			if selected == "up" {
				channelNum = min(channelNum+1, 999)
			} else {
				channelNum = max(channelNum-1, 0)
			}
		} else {
			channelNum, err = strconv.Atoi(selected)
			if err != nil {
				return "", fmt.Errorf("invalid channel number: %s", selected)
			}
		}

		// Set channel - zero-pad to 3 digits
		if err = r.SetCurrentChannel(ctx.Config.Vfo, channelNum); err != nil {
			return "", fmt.Errorf("failed to set channel: %w", err)
		}
	}

	channel, err = r.GetCurrentChannel(ctx.Config.Vfo)
	if err != nil {
		return "", fmt.Errorf("failed to get current channel: %w", err)
	}

	return channel.String(), nil
}
