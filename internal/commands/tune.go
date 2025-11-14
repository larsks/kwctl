package commands

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/helpers"
	"github.com/larsks/kwctl/internal/radio"
	"github.com/larsks/kwctl/internal/types"
)

type (
	TuneCommand struct {
		flags        *flag.FlagSet
		forceVfoMode bool
		vfo          types.VFO
	}
)

func init() {
	Register("tune", &TuneCommand{})
}

//nolint:errcheck
func (c *TuneCommand) Init() error {
	c.flags = flag.NewFlagSet("id", flag.ContinueOnError)
	c.flags.BoolVarP(&c.forceVfoMode, "force", "f", false, "change to vfo mode before tuning")
	c.flags.VarP(types.NewFrequencyMHz(&c.vfo.RxFreq), "rxfreq", "r", "frequency in MHz (e.g., 144.39)")
	c.flags.VarP(types.NewStepSize(&c.vfo.RxStep), "rxstep", "s", "step size in hz (e.g., 5)")
	c.flags.VarP(types.NewMode(&c.vfo.Mode), "mode", "m", "Mode (FM, NFM, AM)")
	c.flags.VarP(types.NewShift(&c.vfo.Shift), "shift", "t", "Shift (simplex, up, down)")
	c.flags.VarP(types.NewFrequencyMHz(&c.vfo.Offset), "offset", "", "offset in MHz (e.g., 0.6)")
	c.flags.StringP("tone-mode", "", "none", "select tone mode (none, tone, tsql, dcs)")
	c.flags.VarP(types.NewTone(&c.vfo.ToneFreq), "txtone", "", "CTCSS tone when sending")
	c.flags.VarP(types.NewTone(&c.vfo.CTCSSFreq), "rxtone", "", "CTCSS tone when receiving")
	c.flags.VarP(types.NewDCS(&c.vfo.DCSCode), "dcs", "", "DCS code")
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), helpers.Unindent(`
			Usage: kwctl tune [options]

			Tune the selected VFO.

			Options:
		`))
		c.flags.PrintDefaults()
	}
	return nil
}

func (c *TuneCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	if err := c.flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	vfo, err := r.GetVFO(ctx.Config.Vfo)
	if err != nil {
		return "", fmt.Errorf("failed to read vfo %s: %w", ctx.Config.Vfo, err)
	}
	oldVfo := vfo

	c.flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "rxfreq":
			vfo.RxFreq = c.vfo.RxFreq
		case "rxstep":
			vfo.RxStep = c.vfo.RxStep
		case "mode":
			vfo.Mode = c.vfo.Mode
		case "shift":
			vfo.Shift = c.vfo.Shift
		case "offset":
			vfo.Offset = c.vfo.Offset
		case "tone-mode":
			switch f.Value.String() {
			case "none":
				vfo.Tone = 0
				vfo.CTCSS = 0
				vfo.DCS = 0
			case "tone":
				vfo.Tone = 1
				vfo.CTCSS = 0
				vfo.DCS = 0
			case "tsql":
				vfo.Tone = 1
				vfo.CTCSS = 1
				vfo.DCS = 0
			case "dcs":
				vfo.Tone = 0
				vfo.CTCSS = 0
				vfo.DCS = 1
			}
		case "txtone":
			vfo.ToneFreq = c.vfo.ToneFreq
		case "rxtone":
			vfo.CTCSSFreq = c.vfo.CTCSSFreq
		case "dcs":
			vfo.DCSCode = c.vfo.DCSCode
		}
	})

	if vfo != oldVfo {
		if c.forceVfoMode {
			err := r.SetVFOMode(ctx.Config.Vfo, radio.VFO_MODE_VFO)
			if err != nil {
				return "", fmt.Errorf("failed to change to vfo mode: %w", err)
			}
		}
		err := r.SetVFO(ctx.Config.Vfo, vfo)
		if err != nil {
			return "", fmt.Errorf("failed to tune vfo %s: %w", ctx.Config.Vfo, err)
		}
		vfo, err = r.GetVFO(ctx.Config.Vfo)
		if err != nil {
			return "", fmt.Errorf("failed to read vfo %s: %w", ctx.Config.Vfo, err)
		}
	}

	return fmt.Sprintf("%s\n", vfo), nil
}
