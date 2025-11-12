package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	"go.bug.st/serial"
)

type (
	Radio struct {
		device string
		config *serial.Mode
		port   serial.Port
	}

	Config struct {
		bitrate string
		verbose bool
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
	flag.BoolVarP(&config.verbose, "verbose", "v", false, "enable verbose logging")
	flag.StringVarP(&config.vfo, "vfo", "", getEnvWithDefault("KWCTL_VFO", "0"), "select vfo on which to operate")
	flag.StringVarP(&config.device, "device", "d", getEnvWithDefault("KWCTL_DEVICE", "/dev/ttyS0"), "serial device")
}

func NewRadio(device string, bitrate int) *Radio {
	return &Radio{
		device: device,
		config: &serial.Mode{
			BaudRate: bitrate,
			Parity:   serial.NoParity,
			DataBits: 8,
			StopBits: serial.OneStopBit,
		},
	}
}

func (r *Radio) Open() error {
	port, err := serial.Open(r.device, r.config)
	if err != nil {
		return fmt.Errorf("failed to open device %s: %w", r.device, err)
	}

	// Set read timeout to prevent blocking indefinitely
	if err := port.SetReadTimeout(100 * time.Millisecond); err != nil {
		port.Close()
		return fmt.Errorf("failed to set read timeout: %w", err)
	}

	r.port = port
	return nil
}

func (r *Radio) Close() error {
	if err := r.port.Close(); err != nil {
		return fmt.Errorf("failed to close device %s: %w", r.device, err)
	}
	return nil
}

func (r *Radio) SendCommand(cmd string, args ...string) (string, error) {
	// Step 1: Clear the serial port by sending a carriage return and discarding response
	if _, err := r.port.Write([]byte("\r")); err != nil {
		return "", fmt.Errorf("failed to clear serial port: %w", err)
	}

	// Ensure data is actually sent to the device
	if err := r.port.Drain(); err != nil {
		return "", fmt.Errorf("failed to flush %s: %w", r.device, err)
	}

	// Read and discard flush response (expect CR-delimited line like "?\r")
	// This avoids hardcoded delays and works at any radio response speed
	flushBuf := make([]byte, 1)
	for {
		n, err := r.port.Read(flushBuf)
		if err != nil || n == 0 {
			break // Timeout or error - buffer was already empty
		}
		if flushBuf[0] == '\r' {
			break // Got end of flush response
		}
		// Continue reading and discarding until we get CR or timeout
	}

	// Step 2: Build and send the command
	var command string
	if len(args) > 0 {
		command = fmt.Sprintf("%s %s\r", cmd, strings.Join(args, ","))
	} else {
		command = fmt.Sprintf("%s\r", cmd)
	}

	logger.Info("sending command", "cmd", command, "device", r.device)
	if _, err := r.port.Write([]byte(command)); err != nil {
		return "", fmt.Errorf("failed to write command: %w", err)
	}

	// Ensure command is actually sent to the device
	if err := r.port.Drain(); err != nil {
		return "", fmt.Errorf("failed to flush %s: %w", r.device, err)
	}

	// Step 3: Read response until carriage return (with overall timeout protection)
	var response []byte
	readBuf := make([]byte, 1)
	deadline := time.Now().Add(2 * time.Second)
	for {
		if time.Now().After(deadline) {
			return "", fmt.Errorf("timeout waiting for response")
		}
		n, err := r.port.Read(readBuf)
		if err != nil {
			return "", fmt.Errorf("failed to read response: %w", err)
		}
		if n > 0 {
			if readBuf[0] == '\r' {
				break
			}
			response = append(response, readBuf[0])
		}
	}

	// Step 4: Parse response (format: "CMD ARG1,ARG2,...")
	responseStr := string(response)
	parts := strings.SplitN(responseStr, " ", 2)
	if len(parts) < 2 {
		return "", nil
	}

	logger.Info("received response", "response", parts[1], "device", r.device)
	return parts[1], nil
}

func (r *Radio) ID() (string, error) {
	return r.SendCommand("ID")
}

func (r *Radio) Check() error {
	id, err := r.ID()
	if err != nil {
		return fmt.Errorf("failed to identify radio at %s: %w", r.device, err)
	}

	if id != "TM-V71" {
		return fmt.Errorf("incompatible radio: want TM-V71, have %s", id)
	}

	return nil
}

func main() {
	flag.Parse()

	// Initialize logger based on verbose flag
	logLevel := slog.LevelWarn
	if config.verbose {
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
	radio := NewRadio(config.device, bitrate)

	if err := radio.Open(); err != nil {
		logger.Error("failed to open radio", "device", config.device, "error", err)
		os.Exit(1)
	}
	defer radio.Close()

	if err := radio.Check(); err != nil {
		logger.Error("radio check failed", "device", config.device, "error", err)
		os.Exit(1)
	}
}
