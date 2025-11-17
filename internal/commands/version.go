package commands

import (
	"fmt"
	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
)

type (
	VersionCommand struct{}
)

var (
	Version string = "dev"
	Commit  string = ""
	Date    string = ""
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
	fmt.Printf("kwctl version %s", Version)
	if Commit != "" {
		fmt.Printf(" (%s)", Commit)
	}
	if Date != "" {
		fmt.Printf(" at %s", Date)
	}
	fmt.Printf("\n")
	return nil
}
