package commands

import (
	"fmt"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
)

type (
	ChannelCommand struct{}
)

func init() {
	Register("channel", &ChannelCommand{})
}

func (c ChannelCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	flags := flag.NewFlagSet("channel", flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	var response string
	var err error

	if flags.NArg() == 0 {
		// Get current channel
		response, err = r.SendCommand("MR", ctx.Config.Vfo)
	} else {
		// Set channel - zero-pad to 3 digits
		channelNum, parseErr := strconv.Atoi(flags.Arg(0))
		if parseErr != nil {
			return "", fmt.Errorf("invalid channel number: %s", flags.Arg(0))
		}
		paddedChannel := fmt.Sprintf("%03d", channelNum)
		response, err = r.SendCommand("MR", ctx.Config.Vfo, paddedChannel)
	}

	if err != nil {
		return "", err
	}

	// Parse response: "vfo,channel"
	parts := strings.Split(response, ",")
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected response format: %s", response)
	}

	return parts[1], nil
}
