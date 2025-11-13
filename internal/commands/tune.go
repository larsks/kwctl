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
	TuneCommand struct{}
)

func init() {
	Register("tune", &TuneCommand{})
}

func (c TuneCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	res, err := r.SendCommand("FO", ctx.Config.Vfo)
	if err != nil {
		return "", fmt.Errorf("unable to read vfo %s: %w", ctx.Config.Vfo, err)
	}

	vfo, err := types.ParseVFO(res)
	if err != nil {
		return "", fmt.Errorf("failed to parse vfo configuration: %w", err)
	}
	oldVfo := vfo

	var forceVfoMode bool
	flags := flag.NewFlagSet("id", flag.ContinueOnError)
	flags.BoolVarP(&forceVfoMode, "force", "f", false, "change to vfo mode before tuning")
	flags.VarP(types.NewFrequencyMHz(&vfo.RxFreq), "rxfreq", "r", "frequency in MHz (e.g., 144.39)")
	flags.VarP(types.NewStepSize(&vfo.RxStep), "rxstep", "s", "step size in hz (e.g., 5)")
	flags.VarP(types.NewMode(&vfo.Mode), "mode", "m", "Mode (FM, NFM, AM)")
	if err := flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	if vfo != oldVfo {
		if forceVfoMode {
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
