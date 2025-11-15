package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/formatters"
	"github.com/larsks/kwctl/internal/helpers"
	"github.com/larsks/kwctl/internal/radio"
	"github.com/larsks/kwctl/internal/types"
)

type (
	ChannelEditCommand struct {
		flags       *flag.FlagSet
		radioFlags  types.RadioFlagValues
		channelName string
		txFreq      int
		txStep      int
		clear       bool
		srcChannel  int
	}
)

func init() {
	Register("edit", &ChannelEditCommand{})
}

func (c *ChannelEditCommand) NeedsRadio() bool {
	return true
}

//nolint:errcheck
func (c *ChannelEditCommand) Init() error {
	c.flags = flag.NewFlagSet("channel-edit", flag.ContinueOnError)

	// Add common radio setting flags
	types.AddRadioSettingFlags(c.flags, &c.radioFlags)

	// Add channel-specific flags
	c.flags.StringVarP(&c.channelName, "name", "n", "", "set channel name")
	c.flags.VarP(types.NewFrequencyMHz(&c.txFreq), "txfreq", "", "frequency in MHz (e.g., 144.39)")
	c.flags.VarP(types.NewStepSize(&c.txStep), "txstep", "", "step size in hz (e.g., 5)")
	c.flags.Bool("lockout", false, "skip channel during scan")
	c.flags.Bool("no-lockout", false, "don't skip channel during scan")
	c.flags.BoolVarP(&c.clear, "clear", "", false, "clear channel")
	c.flags.IntVarP(&c.srcChannel, "copy", "", -1, "copy data from another channel")

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

	if c.flags.NArg() != 1 {
		return "", fmt.Errorf("missing channel number")
	}

	channelNumber, err := strconv.Atoi(c.flags.Arg(0))
	if err != nil {
		return "", fmt.Errorf("invalid channel number")
	}

	if c.clear {
		if err := r.ClearMemoryChannel(channelNumber); err != nil {
			return "", fmt.Errorf("failed to clear channel %d: %w", channelNumber, err)
		}
		return "", nil
	}

	var channel types.Channel
	var oldChannel types.Channel

	if c.srcChannel >= 0 {
		oldChannel = types.Channel{Number: channelNumber}
		channel, err = r.GetMemoryChannel(c.srcChannel)
		if err != nil {
			return "", fmt.Errorf("failed to read channel %03d: %w", c.srcChannel, err)
		}
		channel.Number = channelNumber
	} else {
		channel, err = r.GetMemoryChannel(channelNumber)
		if err != nil {
			if errors.Is(err, radio.ErrUnavailableCommand) {
				channel = types.Channel{Number: channelNumber}
			} else {
				return "", fmt.Errorf("failed to read channel %03d: %w", channelNumber, err)
			}
		}
		oldChannel = channel
	}

	// Apply common radio settings
	types.ApplyRadioSettingFlags(c.flags, &c.radioFlags, &channel)

	// Apply channel-specific flags
	c.flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "txfreq":
			channel.TxFreq = c.txFreq
		case "txstep":
			channel.TxStep = c.txStep
		case "lockout":
			channel.Lockout = 1
		case "no-lockout":
			channel.Lockout = 0
		case "name":
			channel.Name = c.channelName
		}
	})

	if channel == (types.Channel{Number: channelNumber}) {
		return "", nil
	}

	if channel != oldChannel {
		if err := r.SetMemoryChannel(channel); err != nil {
			return "", fmt.Errorf("failed to set channel %d: %w", channelNumber, err)
		}

		channel, err = r.GetMemoryChannel(channelNumber)
		if err != nil {
			return "", fmt.Errorf("failed to read channel %03d: %w", channelNumber, err)
		}
	}

	if ctx.Config.Pretty {
		formatter := formatters.NewTableFormatter(formatters.HeadersFromStruct(types.Channel{}))
		formatter.Update([][]string{channel.Values()})
		formatter.Render(nil)
		return "", nil
	} else {
		return fmt.Sprintf("%s\n", channel), nil
	}
}
