package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"

	"github.com/larsks/kwctl/internal/radio"
)

type (
	Config struct {
		bitrate string
		verbose int
		vfo     string
		device  string
	}
)

var (
	config Config
	logger *slog.Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
)

func getEnvWithDefault(name, default_value string) string {
	val := os.Getenv(name)
	if val == "" {
		val = default_value
	}
	return val
}

func init() {
	flag.StringVarP(&config.bitrate, "bitrate", "b", getEnvWithDefault("KWCTL_BPS", "9600"), "bit rate (serial only)")
	flag.CountVarP(&config.verbose, "verbose", "v", "increase logging verbosity")
	flag.StringVarP(&config.vfo, "vfo", "", getEnvWithDefault("KWCTL_VFO", "0"), "select vfo on which to operate")
	flag.StringVarP(&config.device, "device", "d", getEnvWithDefault("KWCTL_DEVICE", "/dev/ttyS0"), "serial device")
}

func main() {
	flag.Parse()

	// Initialize logger based on verbose flag
	logLevel := slog.LevelWarn
	if config.verbose >= 2 {
		logLevel = slog.LevelDebug
	} else if config.verbose >= 1 {
		logLevel = slog.LevelInfo
	}
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	bitrate, err := strconv.Atoi(config.bitrate)
	if err != nil {
		logger.Error("invalid bitrate", "bitrate", config.bitrate)
		os.Exit(1)
	}
	r := radio.NewRadio(config.device, bitrate).WithLogger(logger)

	if err := r.Open(); err != nil {
		logger.Error("failed to open radio", "device", config.device, "error", err)
		os.Exit(1)
	}
	defer r.Close()

	if err := r.Check(); err != nil {
		logger.Error("radio check failed", "device", config.device, "error", err)
		os.Exit(1)
	}

	// Parse command
	args := flag.Args()
	if len(args) == 0 {
		logger.Error("no command specified")
		os.Exit(1)
	}

	command := args[0]
	commandArgs := args[1:]

	// Route to appropriate command handler
	switch command {
	case "power":
		result, err := r.Power(config.vfo, commandArgs...)
		if err != nil {
			logger.Error("power command failed", "error", err)
			os.Exit(1)
		}
		fmt.Println(result)

	case "channel":
		result, err := r.Channel(config.vfo, commandArgs...)
		if err != nil {
			logger.Error("channel command failed", "error", err)
			os.Exit(1)
		}
		fmt.Println(result)

	case "id":
		result, err := r.ID()
		if err != nil {
			logger.Error("id command failed", "error", err)
			os.Exit(1)
		}
		fmt.Println(result)

	case "vfo":
		result, err := r.VFO(commandArgs...)
		if err != nil {
			logger.Error("id command failed", "error", err)
			os.Exit(1)
		}
		fmt.Println(result)

	default:
		logger.Error("unknown command", "command", command)
		os.Exit(1)
	}
}
