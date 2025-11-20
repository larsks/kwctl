package commands

import (
	"fmt"
	"runtime/debug"

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

func buildSettingsMap(bi *debug.BuildInfo) map[string]string {
	res := make(map[string]string)
	for _, setting := range bi.Settings {
		res[setting.Key] = setting.Value
	}
	return res
}

func (c *VersionCommand) Run(r *radio.Radio, ctx config.Context, args []string) error {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("unable to read buildinfo")
	}

	if ctx.Config.Verbose > 0 {
		fmt.Printf("%s\n", bi)
	} else {
		bsmap := buildSettingsMap(bi)
		fmt.Printf("kwctl %s/%s version %s", bsmap["GOOS"], bsmap["GOARCH"], version.Version)
		if val, ok := bsmap["vcs"]; ok && val == "git" {
			fmt.Printf(" revision %s on %s\n", bsmap["vcs.revision"][:10], bsmap["vcs.time"])
		}
		fmt.Printf("\n")
	}

	return nil
}
