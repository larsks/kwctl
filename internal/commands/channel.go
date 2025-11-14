package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/helpers"
	"github.com/larsks/kwctl/internal/radio"
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
