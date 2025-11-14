package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/commands"
	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/helpers"
	"github.com/larsks/kwctl/internal/radio"
)

var (
	ctx config.Context
)

func init() {
	ctx.Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
	flag.IntVarP(&ctx.Config.Bitrate, "bitrate", "b", helpers.GetEnvWithDefault("KWCTL_BPS", 9600), "bit rate (serial only)")
	flag.CountVarP(&ctx.Config.Verbose, "verbose", "v", "increase logging verbosity")
	flag.StringVarP(&ctx.Config.Vfo, "vfo", "", helpers.GetEnvWithDefault("KWCTL_VFO", "0"), "select vfo on which to operate")
	flag.StringVarP(&ctx.Config.Device, "device", "d", helpers.GetEnvWithDefault("KWCTL_DEVICE", "/dev/ttyS0"), "serial device")
}

func main() {
	flag.SetInterspersed(false)
	flag.Parse()

	// Initialize logger based on verbose flag
	logLevel := slog.LevelWarn
	if ctx.Config.Verbose >= 2 {
		logLevel = slog.LevelDebug
	} else if ctx.Config.Verbose >= 1 {
		logLevel = slog.LevelInfo
	}
	ctx.Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	// Parse command
	args := flag.Args()
	if len(args) == 0 {
		ctx.Logger.Error("no command specified")
		showHelp(os.Stderr)
		os.Exit(1)
	}

	commandName := args[0]
	commandArgs := args[1:]

	if handler := commands.Lookup(commandName); handler != nil {
		var r *radio.Radio

		if handler.NeedsRadio() {
			r = radio.NewRadio(ctx.Config.Device, ctx.Config.Bitrate).WithLogger(ctx.Logger)

			if err := r.Open(); err != nil {
				ctx.Logger.Error("failed to open radio", "device", ctx.Config.Device, "error", err)
				os.Exit(1)
			}
			defer r.Close()

			if err := r.Check(); err != nil {
				ctx.Logger.Error("radio check failed", "device", ctx.Config.Device, "error", err)
				os.Exit(1)
			}
		}

		res, err := handler.Run(r, ctx, commandArgs)
		if err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return
			}
			ctx.Logger.Error("command failed", "command", commandName, "error", err)
			os.Exit(1)
		}

		if res != "" {
			fmt.Printf("%s\n", res)
		}
	} else if commandName == "help" {
		showHelp(os.Stdout)
	} else {
		ctx.Logger.Error("no such command", "command", commandName)
		showHelp(os.Stderr)
		os.Exit(1)
	}
}

func showHelp(out *os.File) {
	fmt.Fprintf(out, "Available commands:\n\n")
	for _, command := range commands.List() {
		fmt.Fprintf(out, "  %s\n", command)
	}
}
