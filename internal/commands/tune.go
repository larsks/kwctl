package commands

import (
	"fmt"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
	"github.com/larsks/kwctl/internal/types"
)

type (
	TuneCommand struct {
		flags        *flag.FlagSet
		forceVfoMode bool
	}
)

func init() {
	Register("tune", &TuneCommand{})
}

func (c *TuneCommand) Init() error {
	c.flags = flag.NewFlagSet("id", flag.ContinueOnError)
	c.flags.BoolVarP(&c.forceVfoMode, "force", "f", false, "change to vfo mode before tuning")
	return nil
}

func (c *TuneCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	res, err := r.SendCommand("FO", ctx.Config.Vfo)
	if err != nil {
		return "", fmt.Errorf("unable to read vfo %s: %w", ctx.Config.Vfo, err)
	}

	vfo, err := types.ParseVFO(res)
	if err != nil {
		return "", fmt.Errorf("failed to parse vfo configuration: %w", err)
	}
	oldVfo := vfo

	c.flags.VarP(types.NewFrequencyMHz(&vfo.RxFreq), "rxfreq", "r", "frequency in MHz (e.g., 144.39)")
	c.flags.VarP(types.NewStepSize(&vfo.RxStep), "rxstep", "s", "step size in hz (e.g., 5)")
	c.flags.VarP(types.NewMode(&vfo.Mode), "mode", "m", "Mode (FM, NFM, AM)")
	c.flags.VarP(types.NewShift(&vfo.Shift), "shift", "t", "Shift (simplex, up, down)")
	if err := c.flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	if vfo != oldVfo {
		if c.forceVfoMode {
			_, err := r.SendCommand("VM", ctx.Config.Vfo, "0")
			if err != nil {
				return "", fmt.Errorf("failed to change to vfo mode: %w", err)
			}
		}
		res, err = r.SendCommand("FO", strings.Split(vfo.Serialize(), ",")...)
		if err != nil {
			return "", fmt.Errorf("failed to tune vfo: %w", err)
		}
	}

	vfo, err = types.ParseVFO(res)
	if err != nil {
		return "", fmt.Errorf("failed to parse vfo configuration: %w", err)
	}
	return fmt.Sprintf("%s\n", vfo), nil
}
