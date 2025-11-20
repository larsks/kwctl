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
		ChannelNumber  int
		ChannelName    string
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
	status, err := r.GetStatus()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}
	jsonData, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}
