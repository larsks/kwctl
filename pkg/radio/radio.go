package radio

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"go.bug.st/serial"

	"github.com/larsks/kwctl/pkg/radio/types"
)

type (
	Radio struct {
		device string
		config *serial.Mode
		port   serial.Port
		logger *slog.Logger
	}
)

var (
	ErrInvalidCommand     = errors.New("invalid command")
	ErrUnavailableCommand = errors.New("command unavailable")
	SupportedRadios       = []string{"TM-V71"}
)

func NewRadio(device string, bitrate int) *Radio {
	return &Radio{
		device: device,
		config: &serial.Mode{
			BaudRate: bitrate,
			Parity:   serial.NoParity,
			DataBits: 8,
			StopBits: serial.OneStopBit,
		},
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn})).With("device", device),
	}
}

func (r *Radio) WithLogger(logger *slog.Logger) *Radio {
	r.logger = logger.With("device", r.device)
	return r
}

func (r *Radio) Open() error {
	port, err := serial.Open(r.device, r.config)
	if err != nil {
		return fmt.Errorf("failed to open device %s: %w", r.device, err)
	}

	// Set read timeout to prevent blocking indefinitely
	if err := port.SetReadTimeout(100 * time.Millisecond); err != nil {
		port.Close() //nolint:errcheck
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

// drainWithRetry wraps port.Drain() with retry logic to handle EINTR errors.
// EINTR (interrupted system call) can occur when the ioctl syscall used by
// Drain() is interrupted by signals, particularly SIGURG from Go's runtime
// scheduler (used for goroutine preemption since Go 1.14).
func (r *Radio) drainWithRetry() error {
	const maxRetries = 10
	for range maxRetries {
		err := r.port.Drain()
		if err == nil {
			return nil
		}
		// Retry only on EINTR; return all other errors immediately
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
		// EINTR received, retry the operation
	}
	return fmt.Errorf("drain failed after %d retries", maxRetries)
}

func (r *Radio) SendCommand(cmd string, args ...string) (string, error) {
	// Step 1: Clear the serial port by sending a carriage return and discarding response
	if _, err := r.port.Write([]byte("\r")); err != nil {
		return "", fmt.Errorf("failed to clear serial port: %w", err)
	}

	// Ensure data is actually sent to the device
	if err := r.drainWithRetry(); err != nil {
		return "", fmt.Errorf("failed to flush %s: %w", r.device, err)
	}

	// Read and discard flush response.
	flushBuf := make([]byte, 1)
	for {
		n, err := r.port.Read(flushBuf)
		if err != nil || n == 0 {
			break // Timeout or error - buffer was already empty
		}
		// Continue reading and discarding until we get timeout
	}

	// Step 2: Build and send the command
	var command string
	if len(args) > 0 {
		command = fmt.Sprintf("%s %s\r", cmd, strings.Join(args, ","))
	} else {
		command = fmt.Sprintf("%s\r", cmd)
	}

	r.logger.Debug("sending command", "cmd", command)
	if _, err := r.port.Write([]byte(command)); err != nil {
		return "", fmt.Errorf("failed to write command: %w", err)
	}

	// Ensure command is actually sent to the device
	if err := r.drainWithRetry(); err != nil {
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
	r.logger.Debug("raw response", "response", responseStr)

	if responseStr == "?" {
		return "", ErrInvalidCommand
	}
	if responseStr == "N" {
		return "", ErrUnavailableCommand
	}
	parts := strings.SplitN(responseStr, " ", 2)
	if len(parts) < 2 {
		return "", nil
	}

	r.logger.Debug("parsed response", "response", parts[1])
	return parts[1], nil
}

// Ensure that we are communicating with a supported radio.
func (r *Radio) Check() error {
	id, err := r.SendCommand("ID")
	if err != nil {
		return fmt.Errorf("failed to identify radio at %s: %w", r.device, err)
	}

	if !slices.Contains(SupportedRadios, id) {
		return fmt.Errorf("unsupported radio: %s", id)
	}

	return nil
}

func (r *Radio) ClearMemoryChannel(channelNumber int) error {
	channelString := fmt.Sprintf("%03d", channelNumber)
	_, err := r.SendCommand("ME", channelString, "C")
	if err != nil {
		return fmt.Errorf("failed to clear channel %d: %w", channelNumber, err)
	}

	return nil
}

func (r *Radio) GetMemoryChannel(channelNumber int) (types.Channel, error) {
	channelString := fmt.Sprintf("%03d", channelNumber)

	res, err := r.SendCommand("MN", channelString)
	if err != nil {
		return types.EmptyChannel, fmt.Errorf("failed to get name for channel %d: %w", channelNumber, err)
	}
	parts := strings.SplitN(res, ",", 2)
	if len(parts) != 2 {
		return types.EmptyChannel, fmt.Errorf("invalid response for channel %d", channelNumber)
	}
	channelName := parts[1]

	res, err = r.SendCommand("ME", channelString)
	if err != nil {
		return types.EmptyChannel, fmt.Errorf("failed to read data for channel %d: %w", channelNumber, err)
	}

	channel, err := types.ParseChannel(res)
	if err != nil {
		return types.EmptyChannel, fmt.Errorf("failed to parse data for channel %d: %w", channelNumber, err)
	}

	channel.Name = channelName

	return channel, nil
}

func (r *Radio) SetMemoryChannel(channel types.Channel) error {
	channelString := fmt.Sprintf("%03d", channel.Number)
	_, err := r.SendCommand("ME", channel.Serialize())
	if err != nil {
		return fmt.Errorf("failed to set channel %d: %w", channel.Number, err)
	}

	if channel.Name != "" {
		_, err := r.SendCommand("MN", channelString, strings.ToUpper(channel.Name))
		if err != nil {
			return fmt.Errorf("failed to set name for channel %d: %w", channel.Number, err)
		}
	}

	return nil
}

func (r *Radio) GetCurrentChannelNumber(vfo string) (int, error) {
	res, err := r.SendCommand("MR", vfo)
	if err != nil {
		return 0, fmt.Errorf("failed to get current channel: %w", err)
	}

	parts := strings.Split(res, ",")
	if len(parts) != 2 {
		return 0, fmt.Errorf("unable to determine current channel")
	}

	channelNum, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("unable to determine current channel: %w", err)
	}

	return channelNum, nil
}

func (r *Radio) GetCurrentChannel(vfo string) (types.Channel, error) {
	channelNum, err := r.GetCurrentChannelNumber(vfo)
	if err != nil {
		return types.EmptyChannel, fmt.Errorf("unable to determine current channel: %w", err)
	}

	channel, err := r.GetMemoryChannel(channelNum)
	if err != nil {
		return types.EmptyChannel, fmt.Errorf("unable to get data for channel %d: %w", channelNum, err)
	}
	return channel, nil
}

func (r *Radio) SetCurrentChannel(vfo string, channelNumber int) error {
	_, err := r.SendCommand("MR", vfo, fmt.Sprintf("%03d", channelNumber))
	if err != nil {
		return fmt.Errorf("failed to set channel: %w", err)
	}

	return nil
}

func (r *Radio) GetVFO(vfo string) (types.VFO, error) {
	res, err := r.SendCommand("FO", vfo)
	if err != nil {
		return types.EmptyVFO, fmt.Errorf("unable to read vfo %s: %w", vfo, err)
	}

	v, err := types.ParseVFO(res)
	if err != nil {
		return types.EmptyVFO, fmt.Errorf("failed to parse vfo configuration: %w", err)
	}
	return v, nil
}

func (r *Radio) SetVFO(vfo string, config types.VFO) error {
	_, err := r.SendCommand("FO", config.Serialize())
	if err != nil {
		return fmt.Errorf("failed to tune vfo %s: %w", vfo, err)
	}

	return nil
}

func (r *Radio) GetVFOMode(vfo string) (types.VfoMode, error) {
	res, err := r.SendCommand("VM", vfo)
	if err != nil {
		return 0, fmt.Errorf("failed to read mode for vfo %s: %w", vfo, err)
	}

	parts := strings.Split(res, ",")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid response: %s", res)
	}

	mode, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("unable to parse vfo response: %w", err)
	}

	return types.VfoMode(mode), nil
}

func (r *Radio) SetVFOMode(vfo string, mode types.VfoMode) error {
	_, err := r.SendCommand("VM", vfo, fmt.Sprintf("%d", mode))
	if err != nil {
		return fmt.Errorf("failed to set mode for vfo %s: %w", vfo, err)
	}

	return nil
}

func (r *Radio) GetTxPower(vfo string) (types.TxPower, error) {
	res, err := r.SendCommand("PC", vfo)
	if err != nil {
		return 0, fmt.Errorf("failed to read tx power for vfo %s: %w", vfo, err)
	}

	parts := strings.Split(res, ",")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid response: %s", res)
	}

	tx, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("unable to parse txpower response: %w", err)
	}

	return types.TxPower(tx), nil
}

func (r *Radio) SetTxPower(vfo string, tx types.TxPower) error {
	_, err := r.SendCommand("PC", vfo, fmt.Sprintf("%d", tx))
	if err != nil {
		return fmt.Errorf("failed to set txpower for vfo %s: %w", vfo, err)
	}

	return nil
}

func (r *Radio) getPttAndControl() (int, int, error) {
	res, err := r.SendCommand("BC")
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get ptt/control: %w", err)
	}

	parts := strings.SplitN(res, ",", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid response: %s", res)
	}

	ctlBand, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid control band: %s", parts[0])
	}

	pttBand, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid ptt band: %s", parts[1])
	}

	return ctlBand, pttBand, nil
}

func (r *Radio) GetControlBand() (int, error) {
	ctlBand, _, err := r.getPttAndControl()
	return ctlBand, err
}

func (r *Radio) GetPTTBand() (int, error) {
	_, pttBand, err := r.getPttAndControl()
	return pttBand, err
}
