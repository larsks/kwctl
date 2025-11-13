package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/commands"
	"github.com/larsks/kwctl/internal/config"
	"github.com/larsks/kwctl/internal/radio"
)

var (
	ctx config.Context
)

func getEnvWithDefault(name, default_value string) string {
	val := os.Getenv(name)
	if val == "" {
		val = default_value
	}
	return val
}

func init() {
	ctx.Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
	flag.StringVarP(&ctx.Config.Bitrate, "bitrate", "b", getEnvWithDefault("KWCTL_BPS", "9600"), "bit rate (serial only)")
	flag.CountVarP(&ctx.Config.Verbose, "verbose", "v", "increase logging verbosity")
	flag.StringVarP(&ctx.Config.Vfo, "vfo", "", getEnvWithDefault("KWCTL_VFO", "0"), "select vfo on which to operate")
	flag.StringVarP(&ctx.Config.Device, "device", "d", getEnvWithDefault("KWCTL_DEVICE", "/dev/ttyS0"), "serial device")
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

	bitrate, err := strconv.Atoi(ctx.Config.Bitrate)
	if err != nil {
		ctx.Logger.Error("invalid bitrate", "bitrate", ctx.Config.Bitrate)
		os.Exit(1)
	}
	// Parse command
	args := flag.Args()
	if len(args) == 0 {
		ctx.Logger.Error("no command specified")
		showHelp(os.Stderr)
		os.Exit(1)
	}

	command := args[0]
	commandArgs := args[1:]

	if handler := commands.Lookup(command); handler != nil {
		r := radio.NewRadio(ctx.Config.Device, bitrate).WithLogger(ctx.Logger)

		if err := r.Open(); err != nil {
			ctx.Logger.Error("failed to open radio", "device", ctx.Config.Device, "error", err)
			os.Exit(1)
		}
		defer r.Close()

		if err := r.Check(); err != nil {
			ctx.Logger.Error("radio check failed", "device", ctx.Config.Device, "error", err)
			os.Exit(1)
		}

		res, err := handler.Run(r, ctx, commandArgs)
		if err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return
			}
			ctx.Logger.Error("command failed", "command", command, "error", err)
			os.Exit(1)
		}

		if res != "" {
			fmt.Printf("%s\n", res)
		}
	} else if command == "help" {
		showHelp(os.Stdout)
	} else {
		ctx.Logger.Error("no such command", "command", command)
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
