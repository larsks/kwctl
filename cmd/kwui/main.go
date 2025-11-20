package main

import (
	"log/slog"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 480
	windowTitle  = "KENWOOD Radio Control"
)

func main() {
	// Parse command-line flags
	kwctlCmd := flag.StringP("kwctl-command", "k", getDefaultKwctlCmd(), "command line to execute kwctl (e.g., 'ssh radio kwctl')")
	verbose := flag.BoolP("verbose", "v", false, "enable verbose logging")
	flag.Parse()

	// Configure log level based on verbose flag
	logLevel := slog.LevelInfo
	if *verbose {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	if err := run(*kwctlCmd, logger); err != nil {
		logger.Error("application error", "error", err)
		os.Exit(1)
	}
}

// getDefaultKwctlCmd returns the default kwctl command from environment or fallback
func getDefaultKwctlCmd() string {
	if cmd := os.Getenv("KWCTL_PATH"); cmd != "" {
		return cmd
	}
	return "kwctl"
}

func run(kwctlCmd string, logger *slog.Logger) error {
	windowFlags := uint32(sdl.WINDOW_SHOWN)

	// Try to initialize SDL with appropriate video driver for console
	if os.Getenv("SDL_VIDEODRIVER") == "" && os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		logger.Info("on console, attempting to use kmsdrm driver")
		os.Setenv("SDL_VIDEODRIVER", "kmsdrm")
		windowFlags = sdl.WINDOW_FULLSCREEN
	}

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		windowTitle,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		windowWidth,
		windowHeight,
		windowFlags,
	)
	if err != nil {
		return err
	}
	defer window.Destroy()

	// Use software renderer for KMSDRM compatibility
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		return err
	}
	logger.Info("using software renderer for KMSDRM compatibility")
	defer renderer.Destroy()

	app := NewApp(renderer, kwctlCmd, logger)
	if err := app.Init(); err != nil {
		return err
	}
	defer app.Cleanup()

	app.Run()
	return nil
}
