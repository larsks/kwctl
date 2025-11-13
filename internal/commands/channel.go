package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
	"github.com/larsks/kwctl/internal/types"
)

type (
	ChannelCommand     struct{}
	ChannelListCommand struct{}
	ChannelEditCommand struct{}
)

func init() {
	Register("channel", &ChannelCommand{})
	Register("channel-list", &ChannelListCommand{}, "list")
	Register("channel-edit", &ChannelEditCommand{}, "edit")
}

func (c ChannelListCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	for channel := range 999 {
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
	return "", nil
}

func (c ChannelEditCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	return "", nil
}

func (c ChannelCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	flags := flag.NewFlagSet("channel", flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	var res string
	var err error
	var channelNum int

	if flags.NArg() == 0 {
		// Get current channel
		res, err = r.SendCommand("MR", ctx.Config.Vfo)
	} else {
		channel := flags.Arg(0)

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
