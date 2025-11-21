package commands

import (
	"fmt"
	"runtime/debug"

	"github.com/larsks/gobot/tools"
	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/version"
	"github.com/larsks/kwctl/pkg/radio"
)

type (
	VersionCommand struct{}
)

func init() {
	Register("version", &VersionCommand{})
}

func (c *VersionCommand) NeedsRadio() bool {
	return false
}

func (c *VersionCommand) Init() error {
	return nil
}

func (c *VersionCommand) Run(r *radio.Radio, ctx config.Context, args []string) error {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("unable to read buildinfo")
	}

	if ctx.Config.Verbose > 0 {
		fmt.Printf("%s\n", bi)
	} else {
		bsmap := tools.BuildInfoMap(bi)
		fmt.Printf("kwctl %s/%s version %s", bsmap["GOOS"], bsmap["GOARCH"], version.Version)
		if val, ok := bsmap["vcs"]; ok && val == "git" {
			fmt.Printf(" revision %s on %s", bsmap["vcs.revision"][:10], bsmap["vcs.time"])
		}
		fmt.Printf("\n")
	}

	return nil
}
