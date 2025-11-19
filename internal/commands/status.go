package commands

import (
	"encoding/json"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/helpers"
	"github.com/larsks/kwctl/pkg/radio"
	"github.com/larsks/kwctl/pkg/radio/types"
)

type (
	StatusCommand struct {
		flags *flag.FlagSet
	}

	vfoStatus struct {
		Vfo            types.DisplayVFO
		Channel        int
		TxPower        string
		Mode           string
		SquelchSetting int
		SquelchStatus  int
	}

	radioStatus struct {
		Vfos   [2]vfoStatus
		PttVfo int
		CtlVfo int
	}
)

func init() {
	Register("status", &StatusCommand{})
}

func (c *StatusCommand) NeedsRadio() bool {
	return true
}

//nolint:errcheck
func (c *StatusCommand) Init() error {
	c.flags = flag.NewFlagSet("status", flag.ContinueOnError)
	c.flags.SetOutput(os.Stdout)
	c.flags.Usage = func() {
		fmt.Fprint(c.flags.Output(), helpers.Unindent(`
			Usage: kwctl status

			Return radio status as a JSON document
		`))
		c.flags.PrintDefaults()
	}
	return nil
}

func (c *StatusCommand) Run(r *radio.Radio, _ config.Context, args []string) error {
	var status radioStatus

	for vfoNum := range 2 {
		vfoString := fmt.Sprintf("%d", vfoNum)
		vfo, err := r.GetVFO(vfoString)
		if err != nil {
			return fmt.Errorf("failed to get vfo info: %w", err)
		}
		status.Vfos[vfoNum].Vfo = vfo.Display()

		txpower, err := r.GetTxPower(vfoString)
		if err != nil {
			return fmt.Errorf("failed to get tx power: %w", err)
		}
		status.Vfos[vfoNum].TxPower = txpower.String()

		mode, err := r.GetVFOMode(vfoString)
		if err != nil {
			return fmt.Errorf("failed to get vfo mode: %w", err)
		}
		status.Vfos[vfoNum].Mode = mode.String()

		channel, err := r.GetCurrentChannelNumber(vfoString)
		if err != nil {
			return fmt.Errorf("failed to get channel number: %w", err)
		}
		status.Vfos[vfoNum].Channel = channel
	}

	pttVfo, err := r.GetPTTBand()
	if err != nil {
		return fmt.Errorf("failed to get ptt vfo: %w", err)
	}

	ctlVfo, err := r.GetControlBand()
	if err != nil {
		return fmt.Errorf("failed to get control vfo: %w", err)
	}

	status.PttVfo = pttVfo
	status.CtlVfo = ctlVfo

	jsonData, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}
