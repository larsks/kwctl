package commands

import (
	"fmt"
	"runtime/debug"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
)

type (
	BuildinfoCommand struct{}
)

func init() {
	Register("buildinfo", &BuildinfoCommand{})
}

func (c *BuildinfoCommand) NeedsRadio() bool {
	return false
}

func (c *BuildinfoCommand) Init() error {
	return nil
}

func (c *BuildinfoCommand) Run(r *radio.Radio, ctx config.Context, args []string) error {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("unable to read buildinfo")
	}
	fmt.Printf("%s\n", bi)
	return nil
}
