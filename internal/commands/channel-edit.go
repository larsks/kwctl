package commands

import (
	"errors"
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
	ChannelEditCommand struct {
		flags      *flag.FlagSet
		channel    types.Channel
		clear      bool
		srcChannel int
	}
)

func init() {
	Register("channel-edit", &ChannelEditCommand{}, "edit")
}

//nolint:errcheck
func (c *ChannelEditCommand) Init() error {
	c.flags = flag.NewFlagSet("channel-edit", flag.ContinueOnError)

	c.flags.StringVarP(&c.channel.Name, "name", "n", "", "set channel name")
	c.flags.VarP(types.NewFrequencyMHz(&c.channel.RxFreq), "rxfreq", "r", "frequency in MHz (e.g., 144.39)")
	c.flags.VarP(types.NewStepSize(&c.channel.RxStep), "rxstep", "s", "step size in hz (e.g., 5)")
	c.flags.VarP(types.NewMode(&c.channel.Mode), "mode", "m", "Mode (FM, NFM, AM)")
	c.flags.VarP(types.NewShift(&c.channel.Shift), "shift", "t", "Shift (simplex, up, down)")
	c.flags.VarP(types.NewFrequencyMHz(&c.channel.Offset), "offset", "", "offset in MHz (e.g., 0.6)")
	c.flags.StringP("tone-mode", "", "none", "select tone mode (none, tone, tsql, dcs)")
	c.flags.VarP(types.NewTone(&c.channel.ToneFreq), "txtone", "", "CTCSS tone when sending")
	c.flags.VarP(types.NewTone(&c.channel.CTCSSFreq), "rxtone", "", "CTCSS tone when receiving")
	c.flags.VarP(types.NewDCS(&c.channel.DCSCode), "dcs", "", "DCS code")
	c.flags.VarP(types.NewFrequencyMHz(&c.channel.TxFreq), "txfreq", "", "frequency in MHz (e.g., 144.39)")
	c.flags.VarP(types.NewStepSize(&c.channel.TxStep), "txstep", "", "step size in hz (e.g., 5)")
	c.flags.Bool("lockout", false, "skip channel during scan")
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

	c.flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "rxfreq":
			channel.RxFreq = c.channel.RxFreq
		case "rxstep":
			channel.RxStep = c.channel.RxStep
		case "txfreq":
			channel.TxFreq = c.channel.TxFreq
		case "txstep":
			channel.TxStep = c.channel.TxStep
		case "mode":
			channel.Mode = c.channel.Mode
		case "shift":
			channel.Shift = c.channel.Shift
		case "offset":
			channel.Offset = c.channel.Offset
		case "lockout":
			val, _ := c.flags.GetBool("lockout")
			if val {
				channel.Lockout = 1
			} else {
				channel.Lockout = 0
			}
		case "tone-mode":
			switch f.Value.String() {
			case "none":
				channel.Tone = 0
				channel.CTCSS = 0
				channel.DCS = 0
			case "tone":
				channel.Tone = 1
				channel.CTCSS = 0
				channel.DCS = 0
			case "tsql":
				channel.Tone = 1
				channel.CTCSS = 1
				channel.DCS = 0
			case "dcs":
				channel.Tone = 0
				channel.CTCSS = 0
				channel.DCS = 1
			}
		case "txtone":
			channel.ToneFreq = c.channel.ToneFreq
		case "rxtone":
			channel.CTCSSFreq = c.channel.CTCSSFreq
		case "dcs":
			channel.DCSCode = c.channel.DCSCode
		case "name":
			channel.Name = c.channel.Name
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

	return fmt.Sprintf("%s\n", channel), nil
}
