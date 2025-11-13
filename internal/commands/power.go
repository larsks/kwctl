package commands

import (
	"fmt"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
)

type (
	PowerCommand struct{}
)

func init() {
	Register("power", &PowerCommand{})
}

func (c PowerCommand) Run(r *radio.Radio, ctx config.Context, args []string) (string, error) {
	flags := flag.NewFlagSet("power", flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	var response string
	var err error

	if flags.NArg() == 0 {
		// Get current power
		response, err = r.SendCommand("PC", ctx.Config.Vfo)
	} else {
		// Set power - map string to integer
		var powerVal string
		switch flags.Arg(0) {
		case "high":
			powerVal = "0"
		case "medium":
			powerVal = "1"
		case "low":
			powerVal = "2"
		default:
			return "", fmt.Errorf("invalid power setting: %s (must be high, medium, or low)", flags.Arg(0))
		}
		response, err = r.SendCommand("PC", ctx.Config.Vfo, powerVal)
	}

	if err != nil {
		return "", err
	}

	// Parse response: "vfo,power"
	parts := strings.Split(response, ",")
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected response format: %s", response)
	}

	// Map power value to human-readable string
	switch parts[1] {
	case "0":
		return "high", nil
	case "1":
		return "medium", nil
	case "2":
		return "low", nil
	default:
		return "", fmt.Errorf("unknown power value: %s", parts[1])
	}
}
