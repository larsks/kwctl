package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

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
	ChannelListCommand struct {
		flags *flag.FlagSet
	}
	ChannelEditCommand struct {
		flags *flag.FlagSet
	}
)

func init() {
	Register("channel", &ChannelCommand{})
	Register("channel-list", &ChannelListCommand{}, "list")
	Register("channel-edit", &ChannelEditCommand{}, "edit")
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

	for _, arg := range c.flags.Args() {
		for channel, err := range helpers.RangeIterator(arg) {
			if err != nil {
				return "", fmt.Errorf("invalid range: %w", err)
			}
			if channel < 0 || channel > 999 {
				return "", fmt.Errorf("invalid range (channels must be between 0 and 999)")
			}
			var parsedChannel types.Channel
			var channelName string

			channelNum := fmt.Sprintf("%03d", channel)
			res, err := r.SendCommand("ME", channelNum)
			if err != nil {
				if errors.Is(err, radio.ErrUnavailableCommand) {
					parsedChannel = types.Channel{}
				} else {
					return "", fmt.Errorf("failed to list channels: %w", err)
				}
			} else {
				parsedChannel, err = types.ParseChannel(res)
				if err != nil {
					return "", fmt.Errorf("invalid channel data: %w", err)
				}

				res, err := r.SendCommand("MN", channelNum)
				if err != nil {
					if !errors.Is(err, radio.ErrUnavailableCommand) {
						return "", fmt.Errorf("failed to get channel name: %w", err)
					}
				} else {
					parts := strings.SplitN(res, ",", 2)
					if len(parts) == 2 {
						channelName = parts[1]
					}
				}
			}

			if parsedChannel.RxFreq == 0 {
				fmt.Printf("[%-6s] %03d\n", channelName, channel)
			} else {
				fmt.Printf("[%-6s] %s\n", channelName, parsedChannel)
			}
		}
	}
	return "", nil
}

//nolint:errcheck
func (c *ChannelEditCommand) Init() error {
	c.flags = flag.NewFlagSet("channel-edit", flag.ContinueOnError)
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), helpers.Unindent(`
			Usage: kwctl channel-edit [options] <channel>

			Edit channel configuration.

			Arguments:
				channel    Channel number to edit (0-999)

			Options:
			`))
		c.flags.PrintDefaults()
	}
	return nil
}

func (c *ChannelEditCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	if err := c.flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	return "", nil
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

	var res string
	var err error
	var channelNum int

	if c.flags.NArg() == 0 {
		// Get current channel
		res, err = r.SendCommand("MR", ctx.Config.Vfo)
	} else {
		channel := c.flags.Arg(0)

		if channel == "up" || channel == "down" {
			res, err := r.SendCommand("MR", ctx.Config.Vfo)
			if err != nil {
				return "", fmt.Errorf("failed to get current channel: %w", err)
			}

			parts := strings.Split(res, ",")
			if len(parts) != 2 {
				return "", fmt.Errorf("unable to determine current channel: %w", err)
			}

			channelNum, err = strconv.Atoi(parts[1])
			if err != nil {
				return "", fmt.Errorf("invalid response from radio: %w", err)
			}

			if channel == "up" {
				channelNum = min(channelNum+1, 999)
			} else {
				channelNum = max(channelNum-1, 0)
			}
		} else {
			channelNum, err = strconv.Atoi(channel)
			if err != nil {
				return "", fmt.Errorf("invalid channel number: %s", channel)
			}
		}

		// Set channel - zero-pad to 3 digits
		paddedChannel := fmt.Sprintf("%03d", channelNum)
		res, err = r.SendCommand("MR", ctx.Config.Vfo, paddedChannel)
	}

	if err != nil {
		return "", err
	}

	// Parse response: "vfo,channel"
	parts := strings.Split(res, ",")
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected response format: %s", res)
	}

	return parts[1], nil
}
