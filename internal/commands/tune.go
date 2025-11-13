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

	flags := flag.NewFlagSet("id", flag.ContinueOnError)
	flags.VarP(types.NewFrequencyMHz(&vfo.RxFreq), "rxfreq", "r", "frequency in MHz (e.g., 144.39)")
	if err := flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	res, err = r.SendCommand("FO", strings.Split(vfo.Serialize(), ",")...)
	if err != nil {
		return "", fmt.Errorf("failed to tune vfo")
	}

	return fmt.Sprintf("%s\n", vfo), nil
}
